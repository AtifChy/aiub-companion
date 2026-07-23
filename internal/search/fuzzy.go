// Package search provides a fuzzy search implementation for searching through a list of items based on a query string.
package search

import (
	"cmp"
	"math"
	"slices"
	"strings"

	"github.com/hbollon/go-edlib"
)

const (
	similarityThreshold float32 = 0.7
	tieBreakerThreshold float64 = 0.05
)

// FieldFunc is a function type that extracts a string field from an item of type T.
type FieldFunc[T any] func(T) string

// CompareFunc is a function type that compares two items of type T and returns an integer indicating their order.
type CompareFunc[T any] func(a, b T) int

type scoreItem[T any] struct {
	item  T
	score float32
}

// FuzzySearch performs a fuzzy search on a slice of items of type T based on the provided query string.
func FuzzySearch[T any](
	items []T,
	query string,
	compare CompareFunc[T],
	fields ...FieldFunc[T],
) []T {
	query = strings.ToLower(query)

	if len(fields) == 0 {
		panic("FuzzySearch: at least one field function must be provided")
	}

	scored := make([]scoreItem[T], 0, len(items))

	for _, item := range items {
		var best float32

		for _, field := range fields {
			score := fuzzyScore(query, field(item))
			best = max(best, score)
		}

		if best >= similarityThreshold {
			scored = append(scored, scoreItem[T]{
				item:  item,
				score: best,
			})
		}
	}

	slices.SortFunc(scored, func(a, b scoreItem[T]) int {
		if math.Abs(float64(a.score-b.score)) > tieBreakerThreshold {
			return cmp.Compare(b.score, a.score)
		}

		if compare != nil {
			return compare(a.item, b.item)
		}

		return cmp.Compare(b.score, a.score)
	})

	result := make([]T, len(scored))
	for i, s := range scored {
		result[i] = s.item
	}
	return result
}

func fuzzyScore(query, text string) float32 {
	if text == "" {
		return 0
	}

	text = strings.ToLower(text)

	wholeScore, err := edlib.StringsSimilarity(query, text, edlib.JaroWinkler)
	if err != nil {
		return 0
	}

	queryFields := strings.Fields(query)

	var sum float32
	for _, qw := range queryFields {
		var best float32
		for tw := range strings.FieldsSeq(text) {
			score, err := edlib.StringsSimilarity(qw, tw, edlib.JaroWinkler)
			if err != nil {
				continue
			}
			best = max(best, score)
		}
		sum += best
	}

	tokenScore := sum / float32(len(queryFields))

	return wholeScore*0.6 + tokenScore*0.4
}
