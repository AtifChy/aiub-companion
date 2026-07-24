package calendar

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"aiub-companion/internal/database/sqlc"
)

type Repository interface {
	GetCalendarCache(ctx context.Context, calType CalendarType) (*AcademicCalendar, bool, error)
	UpsertCalendarCache(ctx context.Context, calType CalendarType, calendar *AcademicCalendar) error
}

type repository struct {
	queries  *sqlc.Queries
	cacheTTL time.Duration
}

func NewRepository(db *sql.DB, cacheTTL time.Duration) Repository {
	return &repository{
		queries:  sqlc.New(db),
		cacheTTL: cacheTTL,
	}
}

// GetCalendarCache retrieves the cached academic calendar for the specified type from the database.
// It returns the calendar, a boolean indicating whether the cache is expired, and an error if any occurred.
func (r *repository) GetCalendarCache(ctx context.Context, calType CalendarType) (*AcademicCalendar, bool, error) {
	row, err := r.queries.GetCalendarCache(ctx, string(calType))
	if err != nil {
		return nil, false, err
	}

	var calendar AcademicCalendar
	if err := json.Unmarshal([]byte(row.Data), &calendar); err != nil {
		return nil, false, fmt.Errorf("unmarshal calendar: %w", err)
	}

	calendar.LastUpdated = row.ScrapedAt
	expired := time.Since(row.ScrapedAt) > r.cacheTTL
	return &calendar, expired, nil
}

// UpsertCalendarCache updates the cached academic calendar for the specified type in the database.
func (r *repository) UpsertCalendarCache(ctx context.Context, calType CalendarType, calendar *AcademicCalendar) error {
	data, err := json.Marshal(calendar)
	if err != nil {
		return fmt.Errorf("marshal calendar: %w", err)
	}

	return r.queries.UpsertCalendarCache(ctx, sqlc.UpsertCalendarCacheParams{
		CalendarType: string(calType),
		Data:         string(data),
	})
}
