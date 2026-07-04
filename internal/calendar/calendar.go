// Package calendar provides a simple calendar implementation with basic functionalities.
package calendar

import (
	"sort"
	"time"
)

type CalendarType string

const (
	CalendarStandard  CalendarType = "standard"
	CalendarLLBBPharm CalendarType = "llb-bpharm"
)

type EventCategory string

const (
	EventExam         EventCategory = "exam"
	EventPayment      EventCategory = "payment"
	EventRegistration EventCategory = "registration"
	EventAcademic     EventCategory = "academic"
	EventDeadline     EventCategory = "deadline"
	EventLab          EventCategory = "lab"
	EventBreak        EventCategory = "break"
)

type AcademicEvent struct {
	Date          time.Time     `json:"date"`
	EndDate       *time.Time    `json:"endDate,omitempty"`
	Week          int           `json:"week,omitempty"`
	WeekDateRange string        `json:"weekDateRange,omitempty"`
	Title         string        `json:"title"`
	Category      EventCategory `json:"category"`
}

func (e *AcademicEvent) IsPast() bool {
	return isBeforeDay(e.Date, time.Now())
}

func (e *AcademicEvent) IsInRange(t time.Time) bool {
	if e.EndDate == nil {
		return isSameDay(e.Date, t)
	}
	return !isBeforeDay(t, e.Date) && !isAfterDay(t, *e.EndDate)
}

func (e *AcademicEvent) DaysUntil() int {
	now := time.Now()
	if isBeforeDay(e.Date, now) {
		return 0
	}
	return daysBetween(e.Date, now)
}

type AcademicCalendar struct {
	Semester    string          `json:"semester"`
	Year        int             `json:"year"`
	Type        CalendarType    `json:"type"`
	TotalWeeks  int             `json:"totalWeeks"`
	Events      []AcademicEvent `json:"events"`
	LastUpdated time.Time       `json:"lastUpdated"`
}

func (c *AcademicCalendar) GetCurrentWeek() int {
	now := time.Now()
	for _, event := range c.Events {
		if event.IsInRange(now) {
			return event.Week
		}
	}
	return 0
}

func (c *AcademicCalendar) GetProgressPercentage() float64 {
	currentWeek := c.GetCurrentWeek()
	if currentWeek == 0 || c.TotalWeeks == 0 {
		return 0
	}
	return float64(currentWeek) / float64(c.TotalWeeks) * 100
}

func (c *AcademicCalendar) GetUpcomingEvents(limit int) []AcademicEvent {
	var upcoming []AcademicEvent
	now := time.Now()

	sorted := make([]AcademicEvent, len(c.Events))
	copy(sorted, c.Events)
	sort.Slice(sorted, func(i, j int) bool {
		return isBeforeDay(sorted[i].Date, sorted[j].Date)
	})

	for _, event := range sorted {
		if isAfterDay(event.Date, now) {
			upcoming = append(upcoming, event)
			if limit > 0 && len(upcoming) >= limit {
				break
			}
		}
	}
	return upcoming
}

func (c *AcademicCalendar) GetEventsByCategory(category EventCategory) []AcademicEvent {
	var filtered []AcademicEvent
	for _, event := range c.Events {
		if event.Category == category {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func (c *AcademicCalendar) GetNextExam() *AcademicEvent {
	now := time.Now()

	sorted := make([]AcademicEvent, len(c.Events))
	copy(sorted, c.Events)
	sort.Slice(sorted, func(i, j int) bool {
		return isBeforeDay(sorted[i].Date, sorted[j].Date)
	})

	for _, event := range sorted {
		if event.Category == EventExam && isAfterDay(event.Date, now) {
			return &event
		}
	}
	return nil
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isBeforeDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	if y1 != y2 {
		return y1 < y2
	}
	if m1 != m2 {
		return m1 < m2
	}
	return d1 < d2
}

func isAfterDay(t1, t2 time.Time) bool {
	return isBeforeDay(t2, t1)
}

func daysBetween(t1, t2 time.Time) int {
	start := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	diff := end.Sub(start)

	days := int(diff.Hours() / 24)
	if days < 0 {
		days = -days
	}

	return days
}
