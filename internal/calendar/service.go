package calendar

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"aiub-companion/internal/database"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const cacheTTL = 24 * time.Hour

type Service struct {
	db *database.Service

	repo   Repository
	client Client

	cache     map[CalendarType]*AcademicCalendar
	cacheTime map[CalendarType]time.Time

	mu sync.RWMutex
}

func NewService(db *database.Service) *Service {
	return &Service{
		db:        db,
		cache:     make(map[CalendarType]*AcademicCalendar),
		cacheTime: make(map[CalendarType]time.Time),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.repo = NewRepository(s.db.DB(), cacheTTL)
	s.client = NewClient()
	return nil
}

// GetAcademicCalendar retrieves the academic calendar for the specified type.
// It first checks the in-memory cache, then the database cache,
// and finally fetches and parses the calendar if needed.
func (s *Service) GetAcademicCalendar(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	s.mu.RLock()
	if cached, ok := s.cache[calType]; ok && time.Since(s.cacheTime[calType]) < cacheTTL {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()

	calendar, err := s.repo.GetCalendarCache(ctx, calType)
	if err == nil {
		s.mu.Lock()
		s.cache[calType] = calendar
		s.cacheTime[calType] = calendar.LastUpdated
		s.mu.Unlock()
		return calendar, nil
	}

	return s.refresh(ctx, calType)
}

func (s *Service) RefreshCalendar(ctx context.Context, calType CalendarType) error {
	_, err := s.refresh(ctx, calType)
	return err
}

func (s *Service) refresh(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	calendar, err := s.fetchAndParse(ctx, calType)
	if err != nil {
		s.mu.RLock()
		stale := s.cache[calType]
		s.mu.RUnlock()
		if stale != nil {
			return stale, nil
		}
		return nil, err
	}

	if err := s.repo.UpsertCalendarCache(ctx, calType, calendar); err != nil {
		slog.Warn("Failed to cache calendar", "type", calType, "error", err)
	}

	s.mu.Lock()
	s.cache[calType] = calendar
	s.cacheTime[calType] = time.Now()
	s.mu.Unlock()

	return calendar, nil
}

func (s *Service) fetchAndParse(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	table, semester, err := s.client.ScrapeCalendar(ctx, calType)
	if err != nil {
		return nil, fmt.Errorf("scrape calendar: %w", err)
	}

	parser := NewParser(calType)
	calendar, err := parser.Parse(table, semester)
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %w", err)
	}

	return calendar, nil
}

// Public methods to expose calendar data

func (s *Service) GetCurrentWeek(ctx context.Context, calType CalendarType) (int, error) {
	calendar, err := s.GetAcademicCalendar(ctx, calType)
	if err != nil {
		return 0, err
	}
	return calendar.GetCurrentWeek(), nil
}

func (s *Service) GetUpcomingEvents(ctx context.Context, calType CalendarType, limit int) ([]AcademicEvent, error) {
	calendar, err := s.GetAcademicCalendar(ctx, calType)
	if err != nil {
		return nil, err
	}
	return calendar.GetUpcomingEvents(limit), nil
}

func (s *Service) GetCurrentOrNextExam(ctx context.Context, calType CalendarType) (*AcademicEvent, error) {
	calendar, err := s.GetAcademicCalendar(ctx, calType)
	if err != nil {
		return nil, err
	}
	return calendar.GetCurrentOrNextExam(), nil
}
