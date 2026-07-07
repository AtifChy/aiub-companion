package routine

import (
	"context"
	"fmt"

	"aiub-companion/internal/database"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
	db   *database.Service
	repo Repository
}

func NewService(db *database.Service) *Service {
	return &Service{db: db}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.repo = NewRepository(s.db.DB())
	return nil
}

func (s *Service) ImportOfferedCourses(ctx context.Context, filePath string) error {
	rows, err := readRows(filePath)
	if err != nil {
		return err
	}

	courses, err := parseCourses(rows)
	if err != nil {
		return err
	}

	return s.saveCourses(ctx, courses)
}

func (s *Service) saveCourses(ctx context.Context, courses []Course) error {
	return s.repo.WithTx(ctx, func(txRepo Repository) error {
		if err := txRepo.ClearOfferedCourses(ctx); err != nil {
			return fmt.Errorf("clear offered courses: %w", err)
		}

		for i := range courses {
			if err := txRepo.InsertOfferedCourse(ctx, courses[i]); err != nil {
				return fmt.Errorf("insert offered course: %w", err)
			}
		}

		return nil
	})
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
