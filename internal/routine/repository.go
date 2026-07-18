package routine

import (
	"context"
	"database/sql"

	"aiub-companion/internal/database"
	"aiub-companion/internal/database/sqlc"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(Repository) error) error

	ListUserRoutine(ctx context.Context) ([]Course, error)
	ListOfferedCourses(ctx context.Context) ([]Course, error)
	InsertOfferedCourse(ctx context.Context, c Course) error
	AddToUserRoutine(ctx context.Context, classID string) error
	RemoveFromUserRoutine(ctx context.Context, classID string) error
	ClearOfferedCourses(ctx context.Context) error
}

type repository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *repository) WithTx(ctx context.Context, fn func(Repository) error) error {
	return database.RunInTx(ctx, r.db, func(tx *sql.Tx) error {
		return fn(&repository{queries: r.queries.WithTx(tx), db: r.db})
	})
}

func (r *repository) ListUserRoutine(ctx context.Context) ([]Course, error) {
	rows, err := r.queries.ListUserRoutine(ctx)
	if err != nil {
		return nil, err
	}
	courses := toCourses(rows)
	return courses, nil
}

func (r *repository) ListOfferedCourses(ctx context.Context) ([]Course, error) {
	rows, err := r.queries.ListOfferedCourses(ctx)
	if err != nil {
		return nil, err
	}
	return toCourses(rows), nil
}

func (r *repository) InsertOfferedCourse(ctx context.Context, c Course) error {
	param := sqlc.InsertOfferedCourseParams{
		ClassID:     c.ClassID,
		CourseCode:  database.StringOrNull(c.CourseCode),
		CourseTitle: c.CourseTitle,
		Section:     c.Section,
		Faculty:     database.StringOrNull(c.Faculty),
		Department:  database.StringOrNull(c.Department),
	}

	for _, schedule := range c.Schedules {
		param.ClassType = database.StringOrNull(schedule.Type)
		param.Day = database.StringOrNull(schedule.Day)
		param.StartTime = database.StringOrNull(schedule.StartTime)
		param.EndTime = database.StringOrNull(schedule.EndTime)
		param.Room = database.StringOrNull(schedule.Room)

		if err := r.queries.InsertOfferedCourse(ctx, param); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) AddToUserRoutine(ctx context.Context, classID string) error {
	return r.queries.AddToUserRoutine(ctx, classID)
}

func (r *repository) RemoveFromUserRoutine(ctx context.Context, classID string) error {
	return r.queries.RemoveFromUserRoutine(ctx, classID)
}

func (r *repository) ClearOfferedCourses(ctx context.Context) error {
	return r.queries.ClearOfferedCourses(ctx)
}

func toCourses(rows []sqlc.OfferedCourse) []Course {
	coursesByID := make(map[string]*Course, len(rows))

	for i := range rows {
		course, exists := coursesByID[rows[i].ClassID]

		if !exists {
			course = &Course{
				ClassID:     rows[i].ClassID,
				CourseCode:  rows[i].CourseCode.String,
				CourseTitle: rows[i].CourseTitle,
				Section:     rows[i].Section,
				Faculty:     rows[i].Faculty.String,
				Department:  rows[i].Department.String,
			}
			coursesByID[rows[i].ClassID] = course
		}

		course.Schedules = append(course.Schedules, Schedule{
			Type:      rows[i].ClassType.String,
			Day:       rows[i].Day.String,
			StartTime: rows[i].StartTime.String,
			EndTime:   rows[i].EndTime.String,
			Room:      rows[i].Room.String,
		})
	}

	courses := make([]Course, 0, len(coursesByID))
	for _, course := range coursesByID {
		courses = append(courses, *course)
	}

	return courses
}
