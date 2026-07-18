package routine

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

var sectionRe = regexp.MustCompile(`\s*\[[A-Z0-9]+\]$`)

func readRows(path string) ([][]string, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	sheet := file.GetSheetName(0)

	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}

	return rows, nil
}

func parseCourses(rows [][]string) ([]Course, error) {
	var (
		idx         columnIndex
		headerFound bool
	)

	coursesByID := make(map[string]*Course, len(rows))

	for _, row := range rows {
		if !headerFound {
			idx, headerFound = parseColumnIndex(row)
			continue
		}

		classID, meta, schedule, ok := parseCourseRow(row, idx)
		if !ok {
			continue
		}

		if course, exists := coursesByID[classID]; exists {
			course.Schedules = append(course.Schedules, schedule)
			continue
		}

		meta.Schedules = []Schedule{schedule}
		coursesByID[classID] = &meta
	}

	if !headerFound {
		return nil, errors.New("header row not found")
	}

	courses := make([]Course, 0, len(coursesByID))
	for _, course := range coursesByID {
		courses = append(courses, *course)
	}

	return courses, nil
}

func parseCourseRow(row []string, idx columnIndex) (classID string, meta Course, schedule Schedule, ok bool) {
	get := func(col int) string {
		if col < 0 || col >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[col])
	}

	classID = get(idx.classID)
	title := sectionRe.ReplaceAllString(get(idx.courseTitle), "")
	section := get(idx.section)

	if classID == "" && title == "" && section == "" {
		return "", Course{}, Schedule{}, false
	}

	meta = Course{
		ClassID:     get(idx.classID),
		CourseCode:  get(idx.courseCode),
		CourseTitle: title,
		Section:     get(idx.section),
		Faculty:     get(idx.faculty),
		Department:  get(idx.department),
	}
	schedule = Schedule{
		Type:      get(idx.classType),
		Day:       get(idx.day),
		StartTime: get(idx.startTime),
		EndTime:   get(idx.endTime),
		Room:      get(idx.room),
	}

	return classID, meta, schedule, true
}

type columnIndex struct {
	classID     int
	courseCode  int
	courseTitle int
	section     int
	faculty     int
	classType   int
	day         int
	startTime   int
	endTime     int
	room        int
	department  int
}

func (ci columnIndex) isValid() bool {
	return ci.classID != -1 && ci.courseTitle != -1
}

func parseColumnIndex(headerRow []string) (columnIndex, bool) {
	idx := columnIndex{
		classID:     -1,
		courseCode:  -1,
		courseTitle: -1,
		section:     -1,
		faculty:     -1,
		classType:   -1,
		day:         -1,
		startTime:   -1,
		endTime:     -1,
		room:        -1,
		department:  -1,
	}

	for col, cell := range headerRow {
		switch strings.TrimSpace(strings.ToLower(cell)) {
		case "class id":
			idx.classID = col
		case "course code":
			idx.courseCode = col
		case "course title":
			idx.courseTitle = col
		case "section":
			idx.section = col
		case "faculty":
			idx.faculty = col
		case "type":
			idx.classType = col
		case "day":
			idx.day = col
		case "start time":
			idx.startTime = col
		case "end time":
			idx.endTime = col
		case "room":
			idx.room = col
		case "department":
			idx.department = col
		}
	}

	return idx, idx.isValid()
}
