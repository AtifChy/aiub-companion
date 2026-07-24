package calendar

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"aiub-companion/internal/database"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const cacheTTL = 24 * time.Hour

type refreshResult struct {
	cal *AcademicCalendar
	err error
}

type Service struct {
	db *database.Service

	repo   Repository
	client Client

	cache     map[CalendarType]*AcademicCalendar
	cacheTime map[CalendarType]time.Time
	inflight  map[CalendarType]chan refreshResult

	mu sync.RWMutex
}

func NewService(db *database.Service) *Service {
	return &Service{
		db:        db,
		cache:     make(map[CalendarType]*AcademicCalendar),
		cacheTime: make(map[CalendarType]time.Time),
		inflight:  make(map[CalendarType]chan refreshResult),
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
	cached, ok := s.cache[calType]
	cacheTime := s.cacheTime[calType]
	s.mu.RUnlock()

	if ok && time.Since(cacheTime) < cacheTTL {
		return cached, nil
	}

	calendar, expired, err := s.repo.GetCalendarCache(ctx, calType)
	if err != nil {
		slog.Error("get calendar cache", "error", err)
	} else {
		s.mu.Lock()
		s.cache[calType] = calendar
		s.cacheTime[calType] = calendar.LastUpdated
		s.mu.Unlock()

		if expired {
			go func() {
				if _, err := s.singleFlightRefresh(context.Background(), calType); err != nil {
					slog.Warn("Failed to refresh calendar in background", "type", calType, "error", err)
				}
			}()
		}

		return calendar, nil
	}

	return s.singleFlightRefresh(ctx, calType)
}

func (s *Service) RefreshCalendar(ctx context.Context, calType CalendarType) error {
	_, err := s.singleFlightRefresh(ctx, calType)
	return err
}

func (s *Service) singleFlightRefresh(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	s.mu.Lock()
	if ch, ok := s.inflight[calType]; ok {
		s.mu.Unlock()

		select {
		case result, open := <-ch:
			if !open {
				s.mu.RLock()
				defer s.mu.RUnlock()
				if cached, ok := s.cache[calType]; ok {
					return cached, nil
				}
				return nil, errors.New("calendar refresh channel closed without result")
			}
			return result.cal, result.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	ch := make(chan refreshResult, 1)
	s.inflight[calType] = ch
	s.mu.Unlock()

	calendar, err := s.refresh(ctx, calType)

	ch <- refreshResult{cal: calendar, err: err}
	close(ch)

	s.mu.Lock()
	delete(s.inflight, calType)
	s.mu.Unlock()

	return calendar, err
}

func (s *Service) refresh(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	calendar, err := s.fetchAndParse(ctx, calType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	calendar.LastUpdated = now

	if err := s.repo.UpsertCalendarCache(ctx, calType, calendar); err != nil {
		slog.Warn("Failed to cache calendar", "type", calType, "error", err)
	}

	s.mu.Lock()
	s.cacheTime[calType] = now
	s.cache[calType] = calendar
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
