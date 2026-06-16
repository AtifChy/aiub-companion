package routine

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"aiub-companion/internal/database"

	"github.com/xuri/excelize/v2"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportOfferedCourses(ctx context.Context, filePath string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	sheetName := file.GetSheetName(0)
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	var (
		idx         columnIndex
		headerFound bool
		courses     []Course
	)

	courses = make([]Course, 0, len(rows))

	sectionRe := regexp.MustCompile(`\s*\[[A-Z0-9]+\]$`)

	for _, row := range rows {
		if !headerFound {
			if idx, headerFound = parseColumnIndex(row); !headerFound {
				continue
			}
			continue
		}

		get := func(col int) string {
			if col < 0 || col >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[col])
		}

		classID := get(idx.classID)
		courseTitle := get(idx.courseTitle)
		courseTitle = sectionRe.ReplaceAllString(courseTitle, "")
		section := get(idx.section)

		if classID == "" || courseTitle == "" || section == "" {
			continue
		}

		courses = append(courses, Course{
			ClassID:     classID,
			CourseCode:  get(idx.courseCode),
			CourseTitle: courseTitle,
			Section:     section,
			Faculty:     get(idx.faculty),
			Type:        get(idx.classType),
			Day:         get(idx.day),
			StartTime:   get(idx.startTime),
			EndTime:     get(idx.endTime),
			Room:        get(idx.room),
			Department:  get(idx.department),
		})
	}

	return s.repo.WithTx(ctx, func(txRepo Repository) error {
		err := txRepo.ClearOfferedCourses(ctx)
		if err != nil {
			return fmt.Errorf("failed to clear offered courses: %w", err)
		}
		for i := range courses {
			err := txRepo.InsertOfferedCourse(ctx, courses[i])
			if err != nil {
				return fmt.Errorf("failed to insert offered course: %w", err)
			}
		}
		return nil
	})
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
		switch strings.ToLower(strings.TrimSpace(cell)) {
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

	return idx, idx.classID != -1 && idx.courseTitle != -1
}

func (s *Service) GetUserRoutine(ctx context.Context) ([]Course, error) {
	return s.repo.GetUserRoutine(ctx)
}

func (s *Service) SearchOfferedCourses(ctx context.Context, query string) ([]Course, error) {
	input := database.SanitizeQuery(query)
	if input == "" {
		return []Course{}, nil
	}
	return s.repo.SearchOfferedCourses(ctx, input)
}

func (s *Service) AddToUserRoutine(ctx context.Context, classID string) error {
	return s.repo.AddToUserRoutine(ctx, classID)
}

func (s *Service) RemoveFromUserRoutine(ctx context.Context, classID string) error {
	return s.repo.RemoveFromUserRoutine(ctx, classID)
}

func (s *Service) ClearOfferedCourses(ctx context.Context) error {
	return s.repo.ClearOfferedCourses(ctx)
}
