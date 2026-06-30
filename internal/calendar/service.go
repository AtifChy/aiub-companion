package calendar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"aiub-companion/internal/notice"
)

var noticeIDMap = map[CalendarType]string{
	CalendarStandard:  "academic-calendar-summer-2025-26--except-llb--bpharm",
	CalendarLLBBPharm: "academic-calendar-summer-2025-26--llb-bpharm",
}

type Service struct {
	notice *notice.Service

	mu        sync.RWMutex
	cache     map[CalendarType]*AcademicCalendar
	cacheTime map[CalendarType]time.Time
	cacheTTL  time.Duration
}

func NewService(notice *notice.Service) *Service {
	return &Service{
		notice:    notice,
		cache:     make(map[CalendarType]*AcademicCalendar),
		cacheTime: make(map[CalendarType]time.Time),
		cacheTTL:  24 * time.Hour,
	}
}

func (s *Service) GetAcademicCalendar(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	s.mu.RLock()
	if cached, ok := s.cache[calType]; ok && time.Since(s.cacheTime[calType]) < s.cacheTTL {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()

	calendar, err := s.fetchAndParse(ctx, calType)
	if err != nil {
		s.mu.RLock()
		cached := s.cache[calType]
		s.mu.RUnlock()
		if cached != nil {
			return cached, nil
		}
		return nil, err
	}

	s.mu.Lock()
	s.cache[calType] = calendar
	s.cacheTime[calType] = time.Now()
	s.mu.Unlock()

	return calendar, nil
}

func (s *Service) RefreshCalendar(ctx context.Context, calType CalendarType) error {
	calendar, err := s.fetchAndParse(ctx, calType)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.cache[calType] = calendar
	s.cacheTime[calType] = time.Now()
	s.mu.Unlock()

	return nil
}

func (s *Service) fetchAndParse(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	noticeID, ok := noticeIDMap[calType]
	if !ok {
		return nil, fmt.Errorf("unknown calendar type: %v", calType)
	}

	n, err := s.notice.GetNoticeDetails(ctx, noticeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notice %s: %w", noticeID, err)
	}

	if n.Content == "" {
		return nil, fmt.Errorf("notice %s has empty content", noticeID)
	}

	parser := NewParser(calType)
	calendar, err := parser.Parse(n.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse calendar: %w", err)
	}

	return calendar, nil
}

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

func (s *Service) GetNextExam(ctx context.Context, calType CalendarType) (*AcademicEvent, error) {
	calendar, err := s.GetAcademicCalendar(ctx, calType)
	if err != nil {
		return nil, err
	}
	return calendar.GetNextExam(), nil
}
