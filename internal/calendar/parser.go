package calendar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var monthMap = map[string]time.Month{
	"jan": time.January,
	"feb": time.February,
	"mar": time.March,
	"apr": time.April,
	"may": time.May,
	"jun": time.June,
	"jul": time.July,
	"aug": time.August,
	"sep": time.September,
	"oct": time.October,
	"nov": time.November,
	"dec": time.December,
}

var (
	reFullDate  = regexp.MustCompile(`(?i)^([a-z]+)\s+(\d{1,2}),?\s+(\d{4})$`)       // Pattern: "July 13, 2023"
	reMonthDay  = regexp.MustCompile(`(?i)^([a-z]+)\s+(\d{1,2})$`)                   // Pattern: "July 13"
	reDateRange = regexp.MustCompile(`(\d{1,2})\s*-\s*(\d{1,2})(?:\s+([A-Za-z]+))?`) // Pattern: "7 - 13" or "28 - 4 Jul"
	reParen     = regexp.MustCompile(`\([^)]*\)`)                                    // Pattern: "(Sat)"
	reWeekNum   = regexp.MustCompile(`(\d+)`)                                        // Pattern: "Week 1", "Week 2", etc.

	reParagraph = regexp.MustCompile(`<p[^>]*>([\s\S]*?)</p>`) // Pattern to match <p>...</p> blocks
	reTag       = regexp.MustCompile(`<[^>]*>`)                // Pattern to match HTML tags
)

type Parser struct {
	year    int
	calType CalendarType
}

func NewParser(calendarType CalendarType) *Parser {
	return &Parser{
		year:    time.Now().Year(),
		calType: calendarType,
	}
}

func (p *Parser) Parse(content string) (*AcademicCalendar, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	semester, year := p.extractHeaderInfo(doc)
	if year != 0 {
		p.year = year
	}

	events, totalWeeks, err := p.parseTable(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table: %v", err)
	}

	return &AcademicCalendar{
		Semester:    semester,
		Year:        p.year,
		Type:        p.calType,
		Events:      events,
		TotalWeeks:  totalWeeks,
		LastUpdated: time.Now(),
	}, nil
}

type rowspanState struct {
	month     time.Month
	week      int
	dateRange string
	counters  map[int]int    // column index to remaining rowspan count
	values    map[int]string // column index to value
}

func (p *Parser) parseTable(doc *html.Node) ([]AcademicEvent, int, error) {
	var events []AcademicEvent
	var totalWeek int

	table := findNode(doc, "table")
	if table == nil {
		return nil, 0, fmt.Errorf("no table found in HTML")
	}

	tbody := findNode(table, "tbody")
	if tbody == nil {
		tbody = table // If no tbody, use table directly
	}

	state := &rowspanState{
		counters: make(map[int]int),
		values:   make(map[int]string),
	}

	for row := tbody.FirstChild; row != nil; row = row.NextSibling {
		if row.Type != html.ElementNode || row.Data != "tr" {
			continue
		}

		rowEvents, week := p.parseRow(row, state)
		if week > totalWeek {
			totalWeek = week
		}
		events = append(events, rowEvents...)
	}

	return events, totalWeek, nil
}

func (p *Parser) parseRow(row *html.Node, state *rowspanState) ([]AcademicEvent, int) {
	cells := collectCells(row)
	if len(cells) < 1 {
		return nil, 0
	}

	// Build a map of logical column index to cell index, considering rowspan
	// Skip columns that are occupied by rowspan from previous rows
	colIdx := 0
	cellIdx := 0
	colValues := make(map[int]string)
	colHTML := make(map[int]string)

	for cellIdx < len(cells) {
		// Skip columns occupied by active rowspans
		for state.counters[colIdx] > 0 {
			// Use the value from the rowspan state
			colValues[colIdx] = state.values[colIdx]
			colIdx++
		}

		if cellIdx >= len(cells) {
			break
		}

		cell := cells[cellIdx]
		colspan := getAttrInt(cell, "colspan", 1)
		rowspan := getAttrInt(cell, "rowspan", 1)
		text := strings.TrimSpace(extractText(cell))
		text = strings.ReplaceAll(text, "\u00a0", " ") // Replace non-breaking spaces

		// Store the value for this column
		colValues[colIdx] = text
		colHTML[colIdx] = extractHTML(cell)

		// If this cell has rowspan, update the rowspan state
		if rowspan > 1 {
			state.counters[colIdx] = rowspan
			state.values[colIdx] = text
		}

		colIdx += colspan
		cellIdx++
	}

	// Fill remaining columns from rowspan if needed
	for state.counters[colIdx] > 0 && colIdx <= 4 {
		colValues[colIdx] = state.values[colIdx]
		colIdx++
	}

	// Decrement rowspan counters for the next row
	for k := range state.counters {
		state.counters[k]--
		if state.counters[k] <= 0 {
			delete(state.counters, k)
			delete(state.values, k)
		}
	}

	// Extract values from columns
	if month := colValues[0]; month != "" {
		if m, ok := parseMonth(month); ok {
			state.month = m
		}
	}
	if dateRange := colValues[1]; dateRange != "" {
		state.dateRange = dateRange
	}
	if week := colValues[2]; week != "" {
		if w := parseWeek(week); w > 0 {
			state.week = w
		}
	}

	dayStr := colValues[3]
	eventHTML := colHTML[4]

	if eventHTML == "" {
		return nil, state.week
	}

	// Parse dates
	dateResult := p.parseDayColumn(dayStr, state.month, state.dateRange)

	// Split multiple events
	eventTitles := splitEvent(eventHTML)

	var events []AcademicEvent
	for _, title := range eventTitles {
		title = cleanTitle(title)
		if title == "" {
			continue
		}

		if dateResult.IsRange {
			events = append(events, AcademicEvent{
				Date:          dateResult.StartDate,
				EndDate:       &dateResult.EndDate,
				Week:          state.week,
				WeekDateRange: state.dateRange,
				Title:         title,
				Category:      categorizeEvent(title),
			})
		} else {
			for _, date := range dateResult.Dates {
				events = append(events, AcademicEvent{
					Date:          date,
					Week:          state.week,
					WeekDateRange: state.dateRange,
					Title:         title,
					Category:      categorizeEvent(title),
				})
			}
		}
	}

	return events, state.week
}

type DateParseResult struct {
	Dates     []time.Time
	IsRange   bool
	StartDate time.Time
	EndDate   time.Time
}

func (p *Parser) parseDayColumn(dayStr string, month time.Month, fallbackRange string) DateParseResult {
	dayStr = strings.TrimSpace(dayStr)

	// If day column is empty, try to use the fallback date range
	if dayStr == "" && fallbackRange != "" {
		if result, ok := p.parseDateRange(fallbackRange, month); ok {
			return result
		}
		return DateParseResult{}
	}

	if dayStr == "" {
		return DateParseResult{}
	}

	// Try parsing as full date first: "July 13, 2023"
	if date, ok := p.parseFullDate(dayStr); ok {
		return DateParseResult{Dates: []time.Time{date}}
	}

	// Try parsing as month-day: "July 13"
	if date, ok := p.parseMonthDay(dayStr); ok {
		return DateParseResult{Dates: []time.Time{date}}
	}

	// Try parsing as date range: "7 - 13" or "28 - 4 Jul"
	if result, ok := p.parseDateRange(dayStr, month); ok {
		return result
	}

	// Parse as individual dates: "7", "11 & 14", "6 (Sat)"
	dates := p.parseIndividualDates(dayStr, month)
	return DateParseResult{Dates: dates}
}

func (p *Parser) parseFullDate(text string) (time.Time, bool) {
	match := reFullDate.FindStringSubmatch(text)
	if len(match) < 4 {
		return time.Time{}, false
	}

	month, ok := parseMonth(match[1])
	if !ok {
		return time.Time{}, false
	}

	day, err1 := strconv.Atoi(match[2])
	year, err2 := strconv.Atoi(match[3])
	if err1 != nil || err2 != nil {
		return time.Time{}, false
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.Local), true
}

func (p *Parser) parseMonthDay(text string) (time.Time, bool) {
	match := reMonthDay.FindStringSubmatch(text)
	if len(match) < 3 {
		return time.Time{}, false
	}

	month, ok := parseMonth(match[1])
	if !ok {
		return time.Time{}, false
	}

	day, err := strconv.Atoi(match[2])
	if err != nil {
		return time.Time{}, false
	}

	return time.Date(p.year, month, day, 0, 0, 0, 0, time.Local), true
}

func (p *Parser) parseDateRange(text string, month time.Month) (DateParseResult, bool) {
	// Normalize dashes
	replacer := strings.NewReplacer("–", "-", "—", "-", "−", "-")
	text = replacer.Replace(text)

	match := reDateRange.FindStringSubmatch(text)
	if len(match) < 3 {
		return DateParseResult{}, false
	}

	startDay, err1 := strconv.Atoi(match[1])
	endDay, err2 := strconv.Atoi(match[2])
	if err1 != nil || err2 != nil {
		return DateParseResult{}, false
	}

	// Determine the end month
	var endMonth time.Month
	if match[3] != "" {
		if m, ok := parseMonth(match[3]); ok {
			endMonth = m
		} else {
			endMonth = month
		}
	} else {
		endMonth = month
	}

	startDate := time.Date(p.year, month, startDay, 0, 0, 0, 0, time.Local)

	var endDate time.Time
	if endMonth != month {
		// Cross-month range
		endYear := p.year
		if endMonth < month {
			endYear++ // Next year
		}
		endDate = time.Date(endYear, endMonth, endDay, 0, 0, 0, 0, time.Local)
	} else if endDay < startDay {
		// Cross-month without explicit month
		nextMonth := month + 1
		nextYear := p.year
		if nextMonth > time.December {
			nextMonth = time.January
			nextYear++
		}
		endDate = time.Date(nextYear, nextMonth, endDay, 0, 0, 0, 0, time.Local)
	} else {
		endDate = time.Date(p.year, endMonth, endDay, 0, 0, 0, 0, time.Local)
	}

	return DateParseResult{
		IsRange:   true,
		StartDate: startDate,
		EndDate:   endDate,
	}, true
}

func (p *Parser) parseIndividualDates(text string, month time.Month) []time.Time {
	var dates []time.Time

	// Clean text
	replacer := strings.NewReplacer("&amp;", "&", "\u00a0", " ")
	text = replacer.Replace(text)

	// Remove parentheses
	text = reParen.ReplaceAllString(text, " ")

	text = strings.ReplaceAll(text, "&", " ")
	parts := strings.FieldsSeq(text)

	for part := range parts {
		day, err := strconv.Atoi(part)
		if err != nil || day < 1 || day > 31 {
			continue
		}
		dates = append(dates, time.Date(p.year, month, day, 0, 0, 0, 0, time.Local))
	}

	return dates
}

func (p *Parser) extractHeaderInfo(doc *html.Node) (semester string, year int) {
	// Find h2 with semester info
	if h2 := findNode(doc, "h2"); h2 != nil {
		semester = extractText(h2)
		parts := strings.Split(semester, ":")
		if len(parts) > 1 {
			semester = parts[1]
		}
	}

	// Find year from table header
	table := findNode(doc, "table")
	if table != nil {
		// Look for <string>2023</string> in first row
		if strong := findNode(table, "strong"); strong != nil {
			text := extractText(strong)
			if y, err := strconv.Atoi(text); err == nil {
				year = y
			}
		}
	}

	if year == 0 {
		year = time.Now().Year()
	}

	return semester, year
}

func parseMonth(text string) (time.Month, bool) {
	text = strings.TrimSpace(strings.ToLower(text))
	if len(text) < 3 {
		return 0, false
	}

	if month, ok := monthMap[text]; ok {
		return month, true
	}

	return 0, false
}

func parseWeek(text string) int {
	text = strings.TrimSpace(strings.ToLower(text))
	if !strings.Contains(text, "week") {
		return 0
	}

	match := reWeekNum.FindStringSubmatch(text)
	if len(match) > 1 {
		week, err := strconv.Atoi(match[1])
		if err == nil {
			return week
		}
	}

	return 0
}

func splitEvent(text string) []string {
	// Clean HTML entites
	replacer := strings.NewReplacer("&amp;", "&", "&nbsp;", " ")
	text = replacer.Replace(text)

	// Try to split by <p> tags
	matches := reParagraph.FindAllStringSubmatch(text, -1)

	var titles []string
	if len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				title := cleanTitle(match[1])
				if title != "" {
					titles = append(titles, title)
				}
			}
		}
	} else {
		// Single event without <p> tags
		title := cleanTitle(text)
		if title != "" {
			titles = append(titles, title)
		}
	}

	return titles
}

func cleanTitle(text string) string {
	text = strings.TrimSpace(text)
	text = stripTags(text)
	// Remove multiple spaces
	text = strings.Join(strings.Fields(text), " ")
	// Remove trailing markers
	text = strings.TrimRight(text, "*† ")
	return text
}

type categoryRule struct {
	keywords []string
	category EventCategory
}

var categoryRules = []categoryRule{
	{[]string{"payment", "installment"}, EventPayment},
	{[]string{"permit", "collection for"}, EventDeadline},
	{[]string{"deadline", "last date"}, EventDeadline},
	{[]string{"result", "grade"}, EventAcademic},
	{[]string{"lab", "laboratory"}, EventLab},
	{[]string{"exam", "midterm", "final"}, EventExam},
	{[]string{"registration"}, EventRegistration},
	{[]string{"break", "holiday"}, EventBreak},
}

func categorizeEvent(title string) EventCategory {
	title = strings.ToLower(title)
	for _, rule := range categoryRules {
		for _, kw := range rule.keywords {
			if strings.Contains(title, kw) {
				return rule.category
			}
		}
	}
	return EventAcademic // Default category
}

func findNode(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNode(c, tag); result != nil {
			return result
		}
	}
	return nil
}

func collectCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	for c := row.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			cells = append(cells, c)
		}
	}
	return cells
}

func extractText(n *html.Node) string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			sb.WriteString(c.Data)
		case html.ElementNode:
			sb.WriteString(extractText(c))
		}
	}
	return strings.TrimSpace(sb.String())
}

func extractHTML(n *html.Node) string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	var render func(*html.Node)
	render = func(node *html.Node) {
		switch node.Type {
		case html.TextNode:
			sb.WriteString(node.Data)
		case html.ElementNode:
			sb.WriteString("<")
			sb.WriteString(node.Data)
			for _, attr := range node.Attr {
				sb.WriteString(" ")
				sb.WriteString(attr.Key)
				sb.WriteString(`="`)
				sb.WriteString(attr.Val)
				sb.WriteString(`"`)
			}
			sb.WriteString(">")
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				render(child)
			}
			sb.WriteString("</")
			sb.WriteString(node.Data)
			sb.WriteString(">")
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		render(c)
	}
	return sb.String()
}

func stripTags(html string) string {
	return reTag.ReplaceAllString(html, "")
}

func getAttrInt(n *html.Node, attr string, defaultVal int) int {
	for _, a := range n.Attr {
		if a.Key == attr {
			if val, err := strconv.Atoi(a.Val); err == nil {
				return val
			}
		}
	}
	return defaultVal
}
