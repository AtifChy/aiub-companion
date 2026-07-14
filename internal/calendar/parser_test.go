package calendar

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

// parseHTML parses a raw HTML string using ParseFromNode.
func parseHTML(p *Parser, content string) (*AcademicCalendar, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %v", err)
	}

	table := findNode(doc, "table")
	if table == nil {
		return nil, errors.New("no table found in HTML")
	}

	return p.Parse(table, "")
}

func TestAcademicEvent_IsPast(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"yesterday", yesterday, true},
		{"tomorrow", tomorrow, false},
		{"now", now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := AcademicEvent{Date: tt.date}
			if got := event.IsPast(); got != tt.expected {
				t.Errorf("IsPast() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAcademicEvent_DaysUntil(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		date     time.Time
		expected int
	}{
		{"tomorrow", now.Add(24 * time.Hour), 1},
		{"next week", now.Add(7 * 24 * time.Hour), 7},
		{"yesterday", now.Add(-24 * time.Hour), 0},
		{"30 days", now.Add(30 * 24 * time.Hour), 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := AcademicEvent{Date: tt.date}
			if got := event.DaysUntil(); got != tt.expected {
				t.Errorf("DaysUntil() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestAcademicEvent_IsInRange(t *testing.T) {
	start := time.Date(2026, time.June, 28, 0, 0, 0, 0, time.Local)
	end := time.Date(2026, time.July, 4, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"start date", start, true},
		{"end date", end, true},
		{"middle date", time.Date(2026, time.July, 1, 0, 0, 0, 0, time.Local), true},
		{"before range", time.Date(2026, time.June, 27, 0, 0, 0, 0, time.Local), false},
		{"after range", time.Date(2026, time.July, 5, 0, 0, 0, 0, time.Local), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := AcademicEvent{Date: start, EndDate: &end}
			if got := event.IsInRange(tt.date); got != tt.expected {
				t.Errorf("IsInRange() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAcademicCalendar_GetCurrentWeek(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		weeks    []Week
		expected int
	}{
		{
			name: "current week found",
			weeks: []Week{
				{Number: 3, Start: now.Add(-14 * 24 * time.Hour), End: now.Add(-7 * 24 * time.Hour)},
				{Number: 4, Start: now.Add(-1 * time.Hour), End: now.Add(7 * 24 * time.Hour)},
			},
			expected: 4,
		},
		{
			name:     "no weeks",
			weeks:    []Week{},
			expected: 0,
		},
		{
			name: "no current week",
			weeks: []Week{
				{Number: 1, Start: now.Add(-21 * 24 * time.Hour), End: now.Add(-14 * 24 * time.Hour)},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := &AcademicCalendar{Weeks: tt.weeks}
			if got := cal.GetCurrentWeek(); got != tt.expected {
				t.Errorf("GetCurrentWeek() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestAcademicCalendar_GetProgressPercentage(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		totalWeeks  int
		currentWeek int
		minExpected float64
		maxExpected float64
	}{
		{"week 4 of 14", 14, 4, 28.0, 29.0},
		{"week 1 of 14", 14, 1, 7.0, 8.0},
		{"zero weeks", 0, 0, 0, 0},
		{"week 0", 14, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := &AcademicCalendar{
				TotalWeeks: tt.totalWeeks,
				Weeks: []Week{
					{Number: tt.currentWeek, Start: now.Add(-1 * time.Hour), End: now.Add(7 * 24 * time.Hour)},
				},
			}
			got := cal.GetProgressPercentage()
			if got < tt.minExpected || got > tt.maxExpected {
				t.Errorf("GetProgressPercentage() = %v, want between %v and %v", got, tt.minExpected, tt.maxExpected)
			}
		})
	}
}

func TestAcademicCalendar_GetUpcomingEvents(t *testing.T) {
	now := time.Now()

	events := []AcademicEvent{
		{Date: now.Add(-48 * time.Hour), Title: "past"},
		{Date: now.Add(24 * time.Hour), Title: "tomorrow"},
		{Date: now.Add(48 * time.Hour), Title: "day after"},
		{Date: now.Add(72 * time.Hour), Title: "3 days"},
	}

	cal := &AcademicCalendar{Events: events}

	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{"no limit", 0, 3},
		{"limit 1", 1, 1},
		{"limit 2", 2, 2},
		{"limit greater than events", 10, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.GetUpcomingEvents(tt.limit)
			if len(got) != tt.expected {
				t.Errorf("GetUpcomingEvents(%d) returned %d events, want %d", tt.limit, len(got), tt.expected)
			}
		})
	}
}

func TestAcademicCalendar_GetEventsByCategory(t *testing.T) {
	events := []AcademicEvent{
		{Title: "Midterm", Category: EventExam},
		{Title: "Payment", Category: EventPayment},
		{Title: "Final", Category: EventExam},
		{Title: "Registration", Category: EventRegistration},
	}

	cal := &AcademicCalendar{Events: events}

	tests := []struct {
		category EventCategory
		expected int
	}{
		{EventExam, 2},
		{EventPayment, 1},
		{EventRegistration, 1},
		{EventLab, 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			got := cal.GetEventsByCategory(tt.category)
			if len(got) != tt.expected {
				t.Errorf("GetEventsByCategory(%s) returned %d events, want %d", tt.category, len(got), tt.expected)
			}
		})
	}
}

func TestAcademicCalendar_GetNextExam(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		events   []AcademicEvent
		expected string
	}{
		{
			name: "finds next exam",
			events: []AcademicEvent{
				{Date: now.Add(-24 * time.Hour), Title: "Past Exam", Category: EventExam},
				{Date: now.Add(24 * time.Hour), Title: "Midterm", Category: EventExam},
				{Date: now.Add(72 * time.Hour), Title: "Final", Category: EventExam},
			},
			expected: "Midterm",
		},
		{
			name: "no upcoming exams",
			events: []AcademicEvent{
				{Date: now.Add(-24 * time.Hour), Title: "Past Exam", Category: EventExam},
			},
			expected: "",
		},
		{
			name: "only non-exam events",
			events: []AcademicEvent{
				{Date: now.Add(24 * time.Hour), Title: "Payment", Category: EventPayment},
			},
			expected: "",
		},
		{
			name: "midterm before final - unsorted input",
			events: []AcademicEvent{
				{Date: now.Add(30 * 24 * time.Hour), Title: "Final Exam", Category: EventExam},
				{Date: now.Add(7 * 24 * time.Hour), Title: "Midterm Exam", Category: EventExam},
				{Date: now.Add(14 * 24 * time.Hour), Title: "Set B Midterm", Category: EventExam},
			},
			expected: "Midterm Exam",
		},
		{
			name: "set a and set b exams",
			events: []AcademicEvent{
				{Date: now.Add(14 * 24 * time.Hour), Title: "Set B Midterm Exam", Category: EventExam},
				{Date: now.Add(7 * 24 * time.Hour), Title: "Set A Midterm Exam", Category: EventExam},
			},
			expected: "Set A Midterm Exam",
		},
		{
			name: "lab exams are categorized as lab not exam",
			events: []AcademicEvent{
				{Date: now.Add(7 * 24 * time.Hour), Title: "Laboratory Midterm exams", Category: EventLab},
				{Date: now.Add(14 * 24 * time.Hour), Title: "Midterm Exam", Category: EventExam},
			},
			expected: "Midterm Exam",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := &AcademicCalendar{Events: tt.events}
			got := cal.GetNextExam()
			if got == nil && tt.expected != "" {
				t.Errorf("GetNextExam() = nil, want %s", tt.expected)
			} else if got != nil && got.Title != tt.expected {
				t.Errorf("GetNextExam() = %s, want %s", got.Title, tt.expected)
			}
		})
	}
}

func TestParser_RealWorldHTML_DebugExams(t *testing.T) {
	oldTimeNow := timeNow
	defer func() { timeNow = oldTimeNow }()
	timeNow = func() time.Time {
		return time.Date(2026, 7, 10, 0, 0, 0, 0, time.Local)
	}

	html := `
<table>
<tbody>
<tr>
<td><p><strong>2026</strong></p></td>
<td><p>Date Range</p></td>
<td><p>Week</p></td>
<td><p>Day</p></td>
<td><p>Events</p></td>
</tr>
<tr>
<td><p><strong>Jul</strong></p></td>
<td><p>12 – 18</p></td>
<td><p>Week 6</p></td>
<td><p>14</p></td>
<td><p>Midterm Exam**</p></td>
</tr>
<tr>
<td><p><strong>Jul</strong></p></td>
<td><p>19 – 25</p></td>
<td><p>Week 7</p></td>
<td><p>21</p></td>
<td><p>Laboratory Midterm exams</p></td>
</tr>
<tr>
<td><p><strong>Aug</strong></p></td>
<td><p>9 – 15</p></td>
<td><p>Week 10</p></td>
<td><p>11</p></td>
<td><p>Final Exam</p></td>
</tr>
<tr>
<td><p><strong>Aug</strong></p></td>
<td><p>16 – 22</p></td>
<td><p>Week 11</p></td>
<td><p>18</p></td>
<td><p>Laboratory Final exams</p></td>
</tr>
</tbody>
</table>`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Debug: print all events with their categories
	t.Logf("Parsed %d events:", len(calendar.Events))
	for i, e := range calendar.Events {
		t.Logf("  [%d] title=%q category=%s date=%v", i, e.Title, e.Category, e.Date)
	}

	if len(calendar.Events) == 0 {
		t.Fatal("No events parsed")
	}

	// Check categories
	midtermExam := findEventByTitle(calendar.Events, "Midterm Exam")
	if midtermExam == nil {
		// List all titles for debugging
		titles := make([]string, 0, len(calendar.Events))
		for _, e := range calendar.Events {
			titles = append(titles, e.Title)
		}
		t.Fatalf("Midterm Exam not found. Available titles: %v", titles)
	}
	if midtermExam.Category != EventExam {
		t.Errorf("Midterm Exam category = %s, want %s", midtermExam.Category, EventExam)
	}

	labMidterm := findEventByTitle(calendar.Events, "Laboratory Midterm exams")
	if labMidterm == nil {
		t.Fatal("Laboratory Midterm exams not found")
	}
	if labMidterm.Category != EventLab {
		t.Errorf("Laboratory Midterm exams category = %s, want %s", labMidterm.Category, EventLab)
	}

	// GetNextExam should return Midterm Exam, not Final Exam
	nextExam := calendar.GetNextExam()
	if nextExam == nil {
		t.Fatal("GetNextExam() returned nil")
	}
	if nextExam.Title != "Midterm Exam" {
		t.Errorf("GetNextExam() = %s, want Midterm Exam", nextExam.Title)
	}
}

func TestParser_AIUBExamFormat(t *testing.T) {
	oldTimeNow := timeNow
	defer func() { timeNow = oldTimeNow }()
	timeNow = func() time.Time {
		return time.Date(2026, 7, 10, 0, 0, 0, 0, time.Local)
	}

	// Test with actual AIUB calendar exam format
	html := `
<table>
<tbody>
<tr>
<td><p><strong>2026</strong></p></td>
<td><p>Date Range</p></td>
<td><p>Week</p></td>
<td><p>Day</p></td>
<td><p>Events</p></td>
</tr>
<tr>
<td><p><strong>Jul</strong></p></td>
<td><p>12 – 18</p></td>
<td><p>Week 6</p></td>
<td><p> </p></td>
<td><p>Midterm Exam for uneven sections**</p></td>
</tr>
<tr>
<td><p><strong>Jul</strong></p></td>
<td><p>19 – 25</p></td>
<td><p>Week 7</p></td>
<td><p> </p></td>
<td><p>Midterm Exam for even sections**</p></td>
</tr>
<tr>
<td><p><strong>Aug</strong></p></td>
<td><p>9 – 15</p></td>
<td><p>Week 10</p></td>
<td><p> </p></td>
<td><p>Set A Final Exam</p></td>
</tr>
<tr>
<td><p><strong>Aug</strong></p></td>
<td><p>16 – 22</p></td>
<td><p>Week 11</p></td>
<td><p> </p></td>
<td><p>Set B Final Exam</p></td>
</tr>
</tbody>
</table>`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	t.Logf("Parsed %d events:", len(calendar.Events))
	for i, e := range calendar.Events {
		t.Logf("  [%d] title=%q category=%s date=%v weekDateRange=%q", i, e.Title, e.Category, e.Date, e.WeekDateRange)
	}

	// All should be categorized as exams
	for _, e := range calendar.Events {
		if e.Category != EventExam {
			t.Errorf("Event %q category = %s, want %s", e.Title, e.Category, EventExam)
		}
	}

	// Check weekDateRange is populated
	for _, e := range calendar.Events {
		if e.WeekDateRange == "" {
			t.Errorf("Event %q has empty weekDateRange", e.Title)
		}
	}

	// GetNextExam should return the first midterm (earliest date)
	nextExam := calendar.GetNextExam()
	if nextExam == nil {
		t.Fatal("GetNextExam() returned nil")
	}
	if nextExam.Title != "Midterm Exam for uneven sections" {
		t.Errorf("GetNextExam() = %s, want Midterm Exam for uneven sections", nextExam.Title)
	}
}

func TestGetNextExam_WithPastExams(t *testing.T) {
	// Simulate current date being after midterm but before final
	// We'll use fixed dates relative to "now"
	now := time.Now()

	// Create calendar with past midterm and future final
	events := []AcademicEvent{
		{Date: now.Add(-30 * 24 * time.Hour), Title: "Midterm Exam", Category: EventExam}, // 30 days ago
		{Date: now.Add(10 * 24 * time.Hour), Title: "Final Exam", Category: EventExam},    // 10 days from now
	}

	cal := &AcademicCalendar{Events: events}
	nextExam := cal.GetNextExam()

	if nextExam == nil {
		t.Fatal("GetNextExam() returned nil")
	}
	if nextExam.Title != "Final Exam" {
		t.Errorf("GetNextExam() = %s, want Final Exam (midterm should be past)", nextExam.Title)
	}
}

func TestGetNextExam_AllPast(t *testing.T) {
	now := time.Now()

	events := []AcademicEvent{
		{Date: now.Add(-30 * 24 * time.Hour), Title: "Midterm Exam", Category: EventExam},
		{Date: now.Add(-10 * 24 * time.Hour), Title: "Final Exam", Category: EventExam},
	}

	cal := &AcademicCalendar{Events: events}
	nextExam := cal.GetNextExam()

	if nextExam != nil {
		t.Errorf("GetNextExam() = %s, want nil (all exams past)", nextExam.Title)
	}
}

func TestParser_DateRangeRowspanWithEmptyDay(t *testing.T) {
	// Actual AIUB calendar structure - day column has rowspan and is empty
	html := `
<table>
<tbody>
<tr>
<td><p><strong>2026</strong></p></td>
<td><p>Date Range</p></td>
<td><p>Week</p></td>
<td><p>Day</p></td>
<td><p>Events</p></td>
</tr>
<tr>
<td rowspan="2"><p><strong>Jul</strong></p></td>
<td rowspan="2"><p>19 – 25</p></td>
<td rowspan="2"><p>Week 7</p></td>
<td rowspan="2"><p>&nbsp;</p></td>
<td><p><strong>Midterm Exam**</strong></p></td>
</tr>
<tr>
<td><p>Some other event</p></td>
</tr>
<tr>
<td><p><strong>Aug</strong></p></td>
<td><p>9 – 15</p></td>
<td><p>Week 10</p></td>
<td><p>&nbsp;</p></td>
<td><p>Final Exam</p></td>
</tr>
</tbody>
</table>`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	t.Logf("Parsed %d events:", len(calendar.Events))
	for i, e := range calendar.Events {
		t.Logf("  [%d] title=%q week=%d weekDateRange=%q date=%v category=%s",
			i, e.Title, e.Week, e.WeekDateRange, e.Date, e.Category)
	}

	if len(calendar.Events) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(calendar.Events))
	}

	// Check midterm exam
	midterm := findEventByTitle(calendar.Events, "Midterm Exam")
	if midterm == nil {
		t.Fatal("Midterm Exam not found")
	}
	if midterm.Category != EventExam {
		t.Errorf("Midterm Exam category = %s, want %s", midterm.Category, EventExam)
	}
	if midterm.Week != 7 {
		t.Errorf("Midterm Exam week = %d, want 7", midterm.Week)
	}
	if midterm.WeekDateRange != "19 – 25" {
		t.Errorf("Midterm Exam weekDateRange = %q, want '19 – 25'", midterm.WeekDateRange)
	}
	// Date should be from the date range column (July 19)
	if midterm.Date.Month() != time.July {
		t.Errorf("Midterm Exam month = %v, want July", midterm.Date.Month())
	}
	if midterm.Date.Day() != 19 {
		t.Errorf("Midterm Exam day = %d, want 19", midterm.Date.Day())
	}

	// Check final exam
	final := findEventByTitle(calendar.Events, "Final Exam")
	if final == nil {
		t.Fatal("Final Exam not found")
	}
	if final.Category != EventExam {
		t.Errorf("Final Exam category = %s, want %s", final.Category, EventExam)
	}

	// GetNextExam should return Midterm (earlier date)
	nextExam := calendar.GetNextExam()
	if nextExam == nil {
		t.Fatal("GetNextExam() returned nil")
	}
	if nextExam.Title != "Midterm Exam" {
		t.Errorf("GetNextExam() = %s, want Midterm Exam", nextExam.Title)
	}
}

func findEventByTitle(events []AcademicEvent, title string) *AcademicEvent {
	for i := range events {
		if events[i].Title == title {
			return &events[i]
		}
	}
	return nil
}

func TestParser_ParseIndividualDates(t *testing.T) {
	p := NewParser(CalendarStandard)

	tests := []struct {
		name     string
		input    string
		month    time.Month
		expected []int // expected days
	}{
		{"single day", "7", time.June, []int{7}},
		{"multiple days with &", "11 & 14", time.June, []int{11, 14}},
		{"day with weekday", "6 (Sat)", time.June, []int{6}},
		{"multiple with spaces", "11 & 14", time.July, []int{11, 14}},
		{"empty string", "", time.June, []int{}},
		{"just weekday", "(Sat)", time.June, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dates := p.parseIndividualDates(tt.input, tt.month)

			if len(dates) != len(tt.expected) {
				t.Errorf("parseIndividualDates(%q) returned %d dates, want %d", tt.input, len(dates), len(tt.expected))
				return
			}

			for i, date := range dates {
				if date.Month() != tt.month || date.Day() != tt.expected[i] {
					t.Errorf("date[%d] = %v, want day %d in %v", i, date, tt.expected[i], tt.month)
				}
			}
		})
	}
}

func TestParser_ParseDateRange(t *testing.T) {
	p := NewParser(CalendarStandard)
	p.year = 2026

	tests := []struct {
		name       string
		input      string
		month      time.Month
		startDay   int
		startMonth time.Month
		startYear  int
		endDay     int
		endMonth   time.Month
		endYear    int
	}{
		{"same month range", "7 - 13", time.June, 7, time.June, 2026, 13, time.June, 2026},
		{"same month with en-dash", "14 – 20", time.June, 14, time.June, 2026, 20, time.June, 2026},
		{"cross-month range", "28 - 4 Jul", time.June, 28, time.June, 2026, 4, time.July, 2026},
		{"cross-month no explicit month", "30 - 5", time.June, 30, time.June, 2026, 5, time.July, 2026},
		{"cross-year range with explicit month", "28 - 3 Jan", time.December, 28, time.December, 2026, 3, time.January, 2027},
		{"cross-year range no explicit month", "30 - 5", time.December, 30, time.December, 2026, 5, time.January, 2027},
		{"cross-month range with en-dash and explicit month", "28 Sep – 4 Oct", time.September, 28, time.September, 2026, 4, time.October, 2026},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := p.parseDateRange(tt.input, tt.month)
			if !ok {
				t.Fatalf("parseDateRange(%q) failed to parse", tt.input)
			}

			if !result.IsRange {
				t.Fatal("expected IsRange to be true")
			}

			if result.StartDate.Day() != tt.startDay || result.StartDate.Month() != tt.startMonth || result.StartDate.Year() != tt.startYear {
				t.Errorf("start date = %v, want %02d-%v-%d", result.StartDate, tt.startDay, tt.startMonth, tt.startYear)
			}

			if result.EndDate.Day() != tt.endDay || result.EndDate.Month() != tt.endMonth || result.EndDate.Year() != tt.endYear {
				t.Errorf("end date = %v, want %02d-%v-%d", result.EndDate, tt.endDay, tt.endMonth, tt.endYear)
			}
		})
	}
}

func TestParser_ParseFullDate(t *testing.T) {
	p := NewParser(CalendarStandard)

	tests := []struct {
		name        string
		input       string
		expected    time.Time
		shouldParse bool
	}{
		{"full date", "Oct 27, 2026", time.Date(2026, time.October, 27, 0, 0, 0, 0, time.Local), true},
		{"invalid date", "Invalid 27, 2026", time.Time{}, false},
		{"partial date", "Oct 27", time.Time{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.parseFullDate(tt.input)
			if ok != tt.shouldParse {
				t.Errorf("parseFullDate(%q) ok = %v, want %v", tt.input, ok, tt.shouldParse)
				return
			}
			if ok && !got.Equal(tt.expected) {
				t.Errorf("parseFullDate(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParser_ParseMonthDay(t *testing.T) {
	p := NewParser(CalendarStandard)

	tests := []struct {
		name        string
		input       string
		expected    time.Time
		shouldParse bool
	}{
		{"Aug 1", "Aug 1", time.Date(2026, time.August, 1, 0, 0, 0, 0, time.Local), true},
		{"Sep 27", "Sep 27", time.Date(2026, time.September, 27, 0, 0, 0, 0, time.Local), true},
		{"invalid month", "Invalid 1", time.Time{}, false},
		{"no day", "Aug", time.Time{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := p.parseMonthDay(tt.input)
			if ok != tt.shouldParse {
				t.Errorf("parseMonthDay(%q) ok = %v, want %v", tt.input, ok, tt.shouldParse)
				return
			}
			if ok && !got.Equal(tt.expected) {
				t.Errorf("parseMonthDay(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCategorizeEvent(t *testing.T) {
	tests := []struct {
		title    string
		expected EventCategory
	}{
		{"Midterm Exam", EventExam},
		{"Final Exam", EventExam},
		{"Midterm Exam**", EventExam},
		{"2nd installment payment", EventPayment},
		{"Pre-registration for Fall", EventRegistration},
		{"Laboratory midterm exams", EventLab},
		{"Laboratory Final exams", EventLab},
		{"Semester break", EventBreak},
		{"Deadline for submission", EventDeadline},
		{"Last date of registration", EventDeadline},
		{"First Day of Classes", EventAcademic},
		{"Scheduled Classes", EventAcademic},
		{"Permit collection for final exam", EventDeadline},
		{"Permit collection for midterm exam", EventDeadline},
		{"Result publication", EventAcademic},
		{"Grade submission", EventAcademic},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if got := categorizeEvent(tt.title); got != tt.expected {
				t.Errorf("categorizeEvent(%q) = %s, want %s", tt.title, got, tt.expected)
			}
		})
	}
}

func TestParser_ParseHTML(t *testing.T) {
	html := `
	<table>
	<tbody>
		<tr>
			<td><p><strong>2026</strong></p></td>
			<td><p><strong>Date Range</strong></p></td>
			<td><p><strong>Week</strong></p></td>
			<td><p><strong>Day</strong></p></td>
			<td><p><strong>Events/Activities</strong></p></td>
		</tr>
		<tr>
			<td><p><strong>Jun</strong></p></td>
			<td><p>7 – 13</p></td>
			<td><p>Week 1</p></td>
			<td><p>7</p></td>
			<td><p><strong>First Day of Classes</strong></p></td>
		</tr>
		<tr>
			<td><p><strong>Jun</strong></p></td>
			<td><p>14 – 20</p></td>
			<td><p>Week 2</p></td>
			<td><p>15</p></td>
			<td><p><strong>Registration closed</strong></p></td>
		</tr>
	</tbody>
	</table>
	`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if calendar.Year != 2026 {
		t.Errorf("expected year 2026, got %d", calendar.Year)
	}

	if calendar.TotalWeeks != 2 {
		t.Errorf("expected total weeks 2, got %d", calendar.TotalWeeks)
	}

	if len(calendar.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(calendar.Events))
	}

	// Check first event
	if len(calendar.Events) > 0 {
		event := calendar.Events[0]
		if event.Title != "First Day of Classes" {
			t.Errorf("expected title 'First Day of Classes', got %q", event.Title)
		}
		if event.Date.Day() != 7 {
			t.Errorf("expected day 7, got %d", event.Date.Day())
		}
		if event.Date.Month() != time.June {
			t.Errorf("expected June, got %v", event.Date.Month())
		}
		if event.Week != 1 {
			t.Errorf("expected week 1, got %d", event.Week)
		}
		if event.Category != EventAcademic {
			t.Errorf("expected category Academic, got %v", event.Category)
		}
	}
}

func TestParser_ParseComplexHTML(t *testing.T) {
	html := `
	<table>
	<tbody>
		<tr>
			<td colspan="2"><p><strong>2026</strong></p></td>
			<td><p><strong>Week</strong></p></td>
			<td><p><strong>Day</strong></p></td>
			<td><p><strong>Events</strong></p></td>
		</tr>
		<tr>
			<td><p><strong>Jun</strong></p></td>
			<td><p>28 – 4 Jul</p></td>
			<td><p>Week 4</p></td>
			<td><p> </p></td>
			<td><p>Scheduled Classes</p></td>
		</tr>
		<tr>
			<td><p><strong>Jul</strong></p></td>
			<td><p>26 Jul – 1 Aug</p></td>
			<td><p>Week 8</p></td>
			<td><p>Aug 1</p></td>
			<td><p>Deadline for grades</p></td>
		</tr>
		<tr>
			<td><p><strong>Sep</strong></p></td>
			<td><p>Sep 14 – Oct 3</p></td>
			<td><p> </p></td>
			<td><p>14 – 17</p></td>
			<td><p>Set B exams</p></td>
		</tr>
		<tr>
			<td><p> </p></td>
			<td><p>Oct 27, 2026</p></td>
			<td><p> </p></td>
			<td><p> </p></td>
			<td><p>Grade conversion</p></td>
		</tr>
	</tbody>
	</table>
	`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check cross-month date range
	foundWeek4 := false
	for _, event := range calendar.Events {
		if event.Week == 4 {
			foundWeek4 = true
			if event.EndDate == nil {
				t.Error("expected Week 4 to have date range (EndDate)")
			} else {
				if event.Date.Month() != time.June || event.Date.Day() != 28 {
					t.Errorf("expected Week 4 start Jun 28, got %v %d", event.Date.Month(), event.Date.Day())
				}
				if event.EndDate.Month() != time.July || event.EndDate.Day() != 4 {
					t.Errorf("expected Week 4 end Jul 4, got %v %d", event.EndDate.Month(), event.EndDate.Day())
				}
			}
		}

		// Check "Aug 1" parsing
		if event.Title == "Deadline for grades" {
			if event.Date.Month() != time.August || event.Date.Day() != 1 {
				t.Errorf("expected Aug 1, got %v %d", event.Date.Month(), event.Date.Day())
			}
		}

		// Check "14 – 17" range
		if event.Title == "Set B exams" {
			if event.EndDate == nil {
				t.Error("expected Set B exams to have date range")
			}
		}

		// Check "Oct 27, 2026" full date
		if event.Title == "Grade conversion" {
			if event.Date.Year() != 2026 || event.Date.Month() != time.October || event.Date.Day() != 27 {
				t.Errorf("expected Oct 27, 2026, got %v", event.Date)
			}
		}
	}

	if !foundWeek4 {
		t.Error("expected to find Week 4 event")
	}
}

func TestParser_MultipleEventsPerCell(t *testing.T) {
	html := `
	<table>
	<tbody>
		<tr>
			<td><p><strong>2026</strong></p></td>
			<td><p>Date Range</p></td>
			<td><p>Week</p></td>
			<td><p>Day</p></td>
			<td><p>Events</p></td>
		</tr>
		<tr>
			<td><p><strong>Jul</strong></p></td>
			<td><p>21 – 27</p></td>
			<td><p>Week 13</p></td>
			<td><p>25</p></td>
			<td>
				<p>Laboratory Final exams</p>
				<p>Permit collection for final exam</p>
			</td>
		</tr>
	</tbody>
	</table>
	`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should create two separate events
	titles := make(map[string]bool)
	for _, event := range calendar.Events {
		titles[event.Title] = true
	}

	if !titles["Laboratory Final exams"] {
		t.Error("expected 'Laboratory Final exams' event")
	}
	if !titles["Permit collection for final exam"] {
		t.Error("expected 'Permit collection for final exam' event")
	}
}

func TestSplitEventText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single event",
			input:    "First Day of Classes",
			expected: []string{"First Day of Classes"},
		},
		{
			name:     "multiple p tags",
			input:    "<p>Laboratory Final exams</p><p>Permit collection</p>",
			expected: []string{"Laboratory Final exams", "Permit collection"},
		},
		{
			name:     "html entities",
			input:    "Adding &amp; Dropping",
			expected: []string{"Adding & Dropping"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitEvent(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("splitEvent() returned %d events, want %d", len(got), len(tt.expected))
				return
			}
			for i, title := range got {
				if title != tt.expected[i] {
					t.Errorf("event[%d] = %q, want %q", i, title, tt.expected[i])
				}
			}
		})
	}
}

func TestParser_RealWorldHTML(t *testing.T) {
	html := `
<table>
<tbody>
<tr>
<td colspan="2">
<p><strong>2026</strong></p>
</td>
<td>
<p><strong>Week</strong></p>
</td>
<td>
<p><strong>Day</strong></p>
</td>
<td>
<p><strong>Events/Activities†</strong></p>
</td>
</tr>
<tr>
<td rowspan="8">
<p><strong>Jun</strong></p>
</td>
<td>
<p> </p>
</td>
<td>
<p> </p>
</td>
<td>
<p><strong>6 (Sat)</strong></p>
</td>
<td>
<p>Freshman Students' Orientation**</p>
</td>
</tr>
<tr>
<td rowspan="3">
<p>7 – 13</p>
</td>
<td rowspan="3">
<p>Week 1</p>
</td>
<td>
<p>7</p>
</td>
<td>
<p><strong>First Day of Classes </strong></p>
</td>
</tr>
<tr>
<td>
<p>9</p>
</td>
<td>
<p>Last date of registration payment and validating registration.</p>
</td>
</tr>
<tr>
<td>
<p>11 &amp; 14</p>
</td>
<td>
<p>Adding/ Dropping**</p>
</td>
</tr>
<tr>
<td rowspan="2">
<p>14 – 20</p>
</td>
<td rowspan="2">
<p>Week 2</p>
</td>
<td>
<p>15</p>
</td>
<td>
<p><strong>Registration closed for Summer 2025-2026</strong></p>
</td>
</tr>
<tr>
<td>
<p>18</p>
</td>
<td>
<p>Submission of TSF and course description</p>
</td>
</tr>
<tr>
<td>
<p>21 – 27</p>
</td>
<td>
<p>Week 3</p>
</td>
<td>
<p>25</p>
</td>
<td>
<p>Automatic conversion of UW, I, blank grades of previous semester to F</p>
</td>
</tr>
<tr>
<td>
<p>28 – 4 Jul</p>
</td>
<td>
<p>Week 4</p>
</td>
<td>
<p> </p>
</td>
<td>
<p>Scheduled Classes</p>
</td>
</tr>
</tbody>
</table>`

	p := NewParser(CalendarStandard)
	calendar, err := parseHTML(p, html)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if calendar.Year != 2026 {
		t.Errorf("expected year 2026, got %d", calendar.Year)
	}

	if calendar.TotalWeeks != 4 {
		t.Errorf("expected total weeks 4, got %d", calendar.TotalWeeks)
	}

	// Verify specific events
	eventMap := make(map[string]AcademicEvent)
	for _, e := range calendar.Events {
		eventMap[e.Title] = e
	}

	// Check "11 & 14" creates two events with same title
	addingDroppingCount := 0
	for _, e := range calendar.Events {
		if e.Title == "Adding/ Dropping" {
			addingDroppingCount++
		}
	}
	if addingDroppingCount != 2 {
		t.Errorf("expected 2 'Adding/ Dropping' events (for Jun 11 and 14), got %d", addingDroppingCount)
	}

	// Check cross-month range "28 – 4 Jul"
	if event, ok := eventMap["Scheduled Classes"]; ok {
		if event.Date.Month() != time.June || event.Date.Day() != 28 {
			t.Errorf("expected Scheduled Classes start Jun 28, got %v %d", event.Date.Month(), event.Date.Day())
		}
		if event.EndDate == nil {
			t.Error("expected Scheduled Classes to have EndDate")
		} else if event.EndDate.Month() != time.July || event.EndDate.Day() != 4 {
			t.Errorf("expected Scheduled Classes end Jul 4, got %v %d", event.EndDate.Month(), event.EndDate.Day())
		}
	} else {
		t.Error("expected 'Scheduled Classes' event")
	}
}
