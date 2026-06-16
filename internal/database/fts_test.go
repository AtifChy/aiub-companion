package database

import "testing"

func TestSanitizeQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Only whitespace",
			input:    "   \t\n  ",
			expected: "",
		},
		{
			name:     "Single word",
			input:    "compiler",
			expected: "compiler*",
		},
		{
			name:     "Multiple words",
			input:    "computer science",
			expected: "computer* AND science*",
		},
		{
			name:     "Punctuation stripped",
			input:    "compiler-design; 101!",
			expected: "compiler* AND design* AND 101*",
		},
		{
			name:     "Extra spaces collapsed",
			input:    "  physics   chemistry  ",
			expected: "physics* AND chemistry*",
		},
		{
			name:     "Unicode characters preserved",
			input:    "বাংলা ভাষা",
			expected: "বাংলা* AND ভাষা*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeQuery(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeQuery(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
