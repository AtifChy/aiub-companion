// Package routine provides routine domain for the application.
package routine

// Course represents a course in the routine.
type Course struct {
	ClassID     string     `json:"classID"`
	CourseCode  string     `json:"courseCode"`
	CourseTitle string     `json:"courseTitle"`
	Section     string     `json:"section"`
	Faculty     string     `json:"faculty"`
	Department  string     `json:"department"`
	Schedules   []Schedule `json:"schedules"`
}

type Schedule struct {
	Type      string `json:"type"`
	Day       string `json:"day"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Room      string `json:"room"`
}
