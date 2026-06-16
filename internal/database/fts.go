package database

import (
	"regexp"
	"strings"
)

var ftsSpecial = regexp.MustCompile(`[^\p{L}\p{N}\s]`)

func SanitizeQuery(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	cleaned := ftsSpecial.ReplaceAllString(input, " ")
	words := strings.Fields(cleaned)
	if len(words) == 0 {
		return ""
	}

	terms := make([]string, len(words))
	for i, word := range words {
		terms[i] = word + "*"
	}

	return strings.Join(terms, " AND ")
}
