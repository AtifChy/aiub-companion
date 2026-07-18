package routine

import "testing"

// fullHeader returns a complete, ordered header row that covers every column.
func fullHeader() []string {
	return []string{
		"Class ID", "Course Code", "Course Title", "Section",
		"Faculty", "Type", "Day", "Start Time", "End Time", "Room", "Department",
	}
}

// fullHeaderIndex returns the columnIndex that fullHeader() should produce.
func fullHeaderIndex() columnIndex {
	return columnIndex{
		classID:     0,
		courseCode:  1,
		courseTitle: 2,
		section:     3,
		faculty:     4,
		classType:   5,
		day:         6,
		startTime:   7,
		endTime:     8,
		room:        9,
		department:  10,
	}
}

func TestParseColumnIndex(t *testing.T) {
	tests := []struct {
		name      string
		headerRow []string
		wantOK    bool
		wantIdx   columnIndex
	}{
		{
			name:      "full header in order",
			headerRow: fullHeader(),
			wantOK:    true,
			wantIdx:   fullHeaderIndex(),
		},
		{
			name: "full header shuffled",
			headerRow: []string{
				"Day", "Class ID", "Room", "Course Title", "Start Time",
				"End Time", "Section", "Faculty", "Type", "Course Code", "Department",
			},
			wantOK: true,
			wantIdx: columnIndex{
				day:         0,
				classID:     1,
				room:        2,
				courseTitle: 3,
				startTime:   4,
				endTime:     5,
				section:     6,
				faculty:     7,
				classType:   8,
				courseCode:  9,
				department:  10,
			},
		},
		{
			name:      "case insensitive and extra whitespace",
			headerRow: []string{"  CLASS ID  ", "  course title  "},
			wantOK:    true,
			wantIdx:   columnIndex{classID: 0, courseCode: -1, courseTitle: 1, section: -1, faculty: -1, classType: -1, day: -1, startTime: -1, endTime: -1, room: -1, department: -1},
		},
		{
			name:      "missing Class ID — invalid",
			headerRow: []string{"Course Code", "Course Title", "Section"},
			wantOK:    false,
			wantIdx:   columnIndex{classID: -1, courseCode: 0, courseTitle: 1, section: 2, faculty: -1, classType: -1, day: -1, startTime: -1, endTime: -1, room: -1, department: -1},
		},
		{
			name:      "missing Course Title — invalid",
			headerRow: []string{"Class ID", "Course Code", "Section"},
			wantOK:    false,
			wantIdx:   columnIndex{classID: 0, courseCode: 1, courseTitle: -1, section: 2, faculty: -1, classType: -1, day: -1, startTime: -1, endTime: -1, room: -1, department: -1},
		},
		{
			name:      "empty header — invalid",
			headerRow: []string{},
			wantOK:    false,
			wantIdx:   columnIndex{classID: -1, courseCode: -1, courseTitle: -1, section: -1, faculty: -1, classType: -1, day: -1, startTime: -1, endTime: -1, room: -1, department: -1},
		},
		{
			name:      "unknown columns ignored",
			headerRow: []string{"Foo", "Bar", "Class ID", "Baz", "Course Title"},
			wantOK:    true,
			wantIdx:   columnIndex{classID: 2, courseCode: -1, courseTitle: 4, section: -1, faculty: -1, classType: -1, day: -1, startTime: -1, endTime: -1, room: -1, department: -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseColumnIndex(tt.headerRow)
			if ok != tt.wantOK {
				t.Errorf("ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.wantIdx {
				t.Errorf("idx =\n  got  %+v\n  want %+v", got, tt.wantIdx)
			}
		})
	}
}

func TestParseCourseRow(t *testing.T) {
	idx := fullHeaderIndex()

	makeRow := func(classID, code, title, section, faculty, typ, day, start, end, room, dept string) []string {
		return []string{classID, code, title, section, faculty, typ, day, start, end, room, dept}
	}

	t.Run("valid full row", func(t *testing.T) {
		row := makeRow("C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE")
		classID, meta, sched, ok := parseCourseRow(row, idx)

		if !ok {
			t.Fatal("expected ok=true")
		}
		if classID != "C001" {
			t.Errorf("classID = %q, want %q", classID, "C001")
		}
		if meta.CourseTitle != "Intro to CS" {
			t.Errorf("CourseTitle = %q, want %q", meta.CourseTitle, "Intro to CS")
		}
		if meta.Faculty != "Dr. Smith" {
			t.Errorf("Faculty = %q, want %q", meta.Faculty, "Dr. Smith")
		}
		if sched.Day != "Sun" {
			t.Errorf("Day = %q, want %q", sched.Day, "Sun")
		}
		if sched.StartTime != "08:00" {
			t.Errorf("StartTime = %q, want %q", sched.StartTime, "08:00")
		}
		if sched.Room != "101" {
			t.Errorf("Room = %q, want %q", sched.Room, "101")
		}
	})

	t.Run("section suffix stripped from course title", func(t *testing.T) {
		row := makeRow("C002", "CSE201", "Data Structures [A]", "A", "Dr. Jones", "Lab", "Mon", "10:00", "11:30", "Lab1", "CSE")
		_, meta, _, ok := parseCourseRow(row, idx)
		if !ok {
			t.Fatal("expected ok=true")
		}
		if meta.CourseTitle != "Data Structures" {
			t.Errorf("CourseTitle = %q, want %q", meta.CourseTitle, "Data Structures")
		}
	})

	t.Run("section suffix with alphanumeric tag stripped", func(t *testing.T) {
		row := makeRow("C003", "EEE301", "Circuit Analysis [B2]", "B2", "", "Theory", "Tue", "09:00", "10:30", "201", "EEE")
		_, meta, _, ok := parseCourseRow(row, idx)
		if !ok {
			t.Fatal("expected ok=true")
		}
		if meta.CourseTitle != "Circuit Analysis" {
			t.Errorf("CourseTitle = %q, want %q", meta.CourseTitle, "Circuit Analysis")
		}
	})

	t.Run("all-empty row returns not-ok", func(t *testing.T) {
		row := makeRow("", "", "", "", "", "", "", "", "", "", "")
		_, _, _, ok := parseCourseRow(row, idx)
		if ok {
			t.Error("expected ok=false for empty row")
		}
	})

	t.Run("short row does not panic", func(t *testing.T) {
		row := []string{"C004"} // fewer columns than the index expects
		classID, _, _, ok := parseCourseRow(row, idx)
		if !ok {
			t.Fatal("expected ok=true (classID is present)")
		}
		if classID != "C004" {
			t.Errorf("classID = %q, want %q", classID, "C004")
		}
	})

	t.Run("whitespace is trimmed from all fields", func(t *testing.T) {
		row := makeRow("  C005  ", "  CSE999  ", "  Advanced Topics  ", "  Z  ", "  Prof. X  ", "  Lab  ", "  Wed  ", "  13:00  ", "  14:30  ", "  305  ", "  CS  ")
		classID, meta, sched, ok := parseCourseRow(row, idx)
		if !ok {
			t.Fatal("expected ok=true")
		}
		if classID != "C005" {
			t.Errorf("classID = %q, want trimmed", classID)
		}
		if meta.CourseCode != "CSE999" {
			t.Errorf("CourseCode = %q, want trimmed", meta.CourseCode)
		}
		if sched.Type != "Lab" {
			t.Errorf("Type = %q, want trimmed", sched.Type)
		}
	})
}

func TestParseCourses(t *testing.T) {
	t.Run("no header row returns error", func(t *testing.T) {
		rows := [][]string{
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE"},
		}
		// None of the above looks like a header, so parseCourses should error.
		_, err := parseCourses(rows)
		if err == nil {
			t.Error("expected error when no header is found")
		}
	})

	t.Run("empty input returns error", func(t *testing.T) {
		_, err := parseCourses([][]string{})
		if err == nil {
			t.Error("expected error for empty input")
		}
	})

	t.Run("header only — no courses — returns empty slice", func(t *testing.T) {
		rows := [][]string{fullHeader()}
		courses, err := parseCourses(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(courses) != 0 {
			t.Errorf("expected 0 courses, got %d", len(courses))
		}
	})

	t.Run("single course single schedule", func(t *testing.T) {
		rows := [][]string{
			fullHeader(),
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE"},
		}
		courses, err := parseCourses(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(courses) != 1 {
			t.Fatalf("expected 1 course, got %d", len(courses))
		}
		if len(courses[0].Schedules) != 1 {
			t.Errorf("expected 1 schedule, got %d", len(courses[0].Schedules))
		}
		if courses[0].Schedules[0].Day != "Sun" {
			t.Errorf("schedule Day = %q, want %q", courses[0].Schedules[0].Day, "Sun")
		}
	})

	t.Run("same classID rows are merged into one course with multiple schedules", func(t *testing.T) {
		rows := [][]string{
			fullHeader(),
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE"},
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Lab", "Mon", "10:00", "11:30", "Lab1", "CSE"},
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Tue", "08:00", "09:30", "101", "CSE"},
		}
		courses, err := parseCourses(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(courses) != 1 {
			t.Fatalf("expected 1 course, got %d", len(courses))
		}
		if len(courses[0].Schedules) != 3 {
			t.Errorf("expected 3 schedules, got %d", len(courses[0].Schedules))
		}
	})

	t.Run("different classIDs produce separate courses", func(t *testing.T) {
		rows := [][]string{
			fullHeader(),
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE"},
			{"C002", "CSE201", "Data Structures", "B", "Dr. Jones", "Theory", "Mon", "09:00", "10:30", "202", "CSE"},
			{"C003", "EEE301", "Circuit Analysis", "A", "Prof. Ali", "Lab", "Tue", "11:00", "12:30", "Lab2", "EEE"},
		}
		courses, err := parseCourses(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(courses) != 3 {
			t.Errorf("expected 3 courses, got %d", len(courses))
		}
	})

	t.Run("blank rows between data rows are ignored", func(t *testing.T) {
		rows := [][]string{
			fullHeader(),
			{"C001", "CSE101", "Intro to CS", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE"},
			{"", "", "", "", "", "", "", "", "", "", ""},
			{"C002", "CSE201", "Data Structures", "B", "Dr. Jones", "Theory", "Mon", "09:00", "10:30", "202", "CSE"},
		}
		courses, err := parseCourses(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(courses) != 2 {
			t.Errorf("expected 2 courses (blank row skipped), got %d", len(courses))
		}
	})

	t.Run("section tag is stripped from course title", func(t *testing.T) {
		rows := [][]string{
			fullHeader(),
			{"C001", "CSE101", "Intro to CS [A]", "A", "Dr. Smith", "Theory", "Sun", "08:00", "09:30", "101", "CSE"},
		}
		courses, err := parseCourses(rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if courses[0].CourseTitle != "Intro to CS" {
			t.Errorf("CourseTitle = %q, want section tag stripped", courses[0].CourseTitle)
		}
	})
}
