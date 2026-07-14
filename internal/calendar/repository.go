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
	GetCalendarCache(ctx context.Context, calType CalendarType) (*AcademicCalendar, error)
	UpsertCalendarCache(ctx context.Context, calType CalendarType, calendar *AcademicCalendar) error
}

type repository struct {
	queries  *sqlc.Queries
	cacheTTL time.Duration
}

func NewRepository(dbConn *sql.DB, cacheTTL time.Duration) Repository {
	return &repository{
		queries:  sqlc.New(dbConn),
		cacheTTL: cacheTTL,
	}
}

func (r *repository) GetCalendarCache(ctx context.Context, calType CalendarType) (*AcademicCalendar, error) {
	row, err := r.queries.GetCalendarCache(ctx, string(calType))
	if err != nil {
		return nil, err
	}

	scrapedAt, err := time.Parse(time.DateTime, row.ScrapedAt)
	if err != nil {
		return nil, fmt.Errorf("parse scraped_at: %w", err)
	} else if time.Since(scrapedAt) > r.cacheTTL {
		return nil, fmt.Errorf("cache expired")
	}

	var calendar AcademicCalendar
	if err := json.Unmarshal([]byte(row.Data), &calendar); err != nil {
		return nil, fmt.Errorf("unmarshal calendar: %w", err)
	}

	return &calendar, nil
}

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
