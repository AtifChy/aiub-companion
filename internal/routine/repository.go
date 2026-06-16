package routine

import (
	"context"
	"database/sql"

	"aiub-companion/internal/database"
	"aiub-companion/internal/database/db"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(Repository) error) error
	GetUserRoutine(ctx context.Context) ([]Course, error)
	SearchOfferedCourses(ctx context.Context, query string) ([]Course, error)
	InsertOfferedCourse(ctx context.Context, c Course) error
	AddToUserRoutine(ctx context.Context, classID string) error
	RemoveFromUserRoutine(ctx context.Context, classID string) error
	ClearOfferedCourses(ctx context.Context) error
}

type sqliteRepository struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewRepository(dbConn *sql.DB) Repository {
	return &sqliteRepository{
		queries: db.New(dbConn),
		dbConn:  dbConn,
	}
}

func (r *sqliteRepository) WithTx(ctx context.Context, fn func(Repository) error) error {
	return database.RunInTx(ctx, r.dbConn, r.queries, func(qtx *db.Queries) error {
		return fn(&sqliteRepository{queries: qtx, dbConn: r.dbConn})
	})
}

func (r *sqliteRepository) GetUserRoutine(ctx context.Context) ([]Course, error) {
	rows, err := r.queries.GetUserRoutine(ctx)
	if err != nil {
		return nil, err
	}
	courses := toCourses(rows)
	return courses, nil
}

func (r *sqliteRepository) SearchOfferedCourses(ctx context.Context, query string) ([]Course, error) {
	search := database.StringOrNull(query)
	rows, err := r.queries.SearchOfferedCourses(ctx, search)
	if err != nil {
		return nil, err
	}
	return toCourses(rows), nil
}

func (r *sqliteRepository) InsertOfferedCourse(ctx context.Context, c Course) error {
	return r.queries.InsertOfferedCourse(ctx, db.InsertOfferedCourseParams{
		ClassID:     c.ClassID,
		CourseCode:  database.StringOrNull(c.CourseCode),
		CourseTitle: c.CourseTitle,
		Section:     c.Section,
		Faculty:     database.StringOrNull(c.Faculty),
		ClassType:   database.StringOrNull(c.Type),
		Day:         database.StringOrNull(c.Day),
		StartTime:   database.StringOrNull(c.StartTime),
		EndTime:     database.StringOrNull(c.EndTime),
		Room:        database.StringOrNull(c.Room),
		Department:  database.StringOrNull(c.Department),
	})
}

func (r *sqliteRepository) AddToUserRoutine(ctx context.Context, classID string) error {
	return r.queries.AddToUserRoutine(ctx, classID)
}

func (r *sqliteRepository) RemoveFromUserRoutine(ctx context.Context, classID string) error {
	return r.queries.RemoveFromUserRoutine(ctx, classID)
}

func (r *sqliteRepository) ClearOfferedCourses(ctx context.Context) error {
	return r.queries.ClearOfferedCourses(ctx)
}

func toCourses(rows []db.OfferedCourse) []Course {
	courses := make([]Course, len(rows))
	for i := range rows {
		courses[i] = Course{
			ClassID:     rows[i].ClassID,
			CourseCode:  rows[i].CourseCode.String,
			CourseTitle: rows[i].CourseTitle,
			Section:     rows[i].Section,
			Faculty:     rows[i].Faculty.String,
			Type:        rows[i].ClassType.String,
			Day:         rows[i].Day.String,
			StartTime:   rows[i].StartTime.String,
			EndTime:     rows[i].EndTime.String,
			Room:        rows[i].Room.String,
			Department:  rows[i].Department.String,
		}
	}
	return courses
}
