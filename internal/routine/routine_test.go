package routine

import "testing"

func TestParseColumnIndex(t *testing.T) {
	tests := []struct {
		name       string
		headerRow  []string
		wantOK     bool
		expectedID int
	}{
		{
			name:       "All columns matching exactly",
			headerRow:  []string{"Class ID", "Course Code", "Course Title", "Section", "Faculty"},
			wantOK:     true,
			expectedID: 0,
		},
		{
			name:       "Mixed case and spaces",
			headerRow:  []string{"  class id  ", "course code", "  COURSE TITLE  "},
			wantOK:     true,
			expectedID: 0,
		},
		{
			name:       "Missing Class ID",
			headerRow:  []string{"Course Code", "Course Title", "Section"},
			wantOK:     false,
			expectedID: -1,
		},
		{
			name:       "Missing Course Title",
			headerRow:  []string{"Class ID", "Course Code", "Section"},
			wantOK:     false,
			expectedID: 0, // ID is found, but overall match is false because Course Title is missing
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, ok := parseColumnIndex(tt.headerRow)
			if ok != tt.wantOK {
				t.Errorf("parseColumnIndex() ok = %v, want %v", ok, tt.wantOK)
			}
			if idx.classID != tt.expectedID {
				t.Errorf("expected classID index %d, got %d", tt.expectedID, idx.classID)
			}
		})
	}
}
