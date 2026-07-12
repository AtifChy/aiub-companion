package notice

import (
	"cmp"
	"slices"
	"strings"

	"github.com/hbollon/go-edlib"
)

const similarityThreshold float32 = 0.7

type scoreNotice struct {
	notice Notice
	score  float32
}

func fuzzyFilter(notices []Notice, query string) []Notice {
	query = strings.TrimSpace(query)

	scored := make([]scoreNotice, 0, len(notices))

	for _, n := range notices {
		score := max(fuzzyScore(query, n.Title), fuzzyScore(query, n.Summary))

		if score >= similarityThreshold {
			scored = append(scored, scoreNotice{
				notice: n,
				score:  score,
			})
		}
	}

	slices.SortFunc(scored, func(a, b scoreNotice) int {
		if c := cmp.Compare(b.score, a.score); c != 0 {
			return c
		}
		return cmp.Compare(b.notice.PostedDate, a.notice.PostedDate)
	})

	result := make([]Notice, len(scored))
	for i, sn := range scored {
		result[i] = sn.notice
	}

	return result
}

func fuzzyScore(query, text string) float32 {
	if text == "" {
		return 0
	}

	text = strings.ToLower(text)

	var best float32
	for word := range strings.FieldsSeq(text) {
		sim, err := edlib.StringsSimilarity(query, word, edlib.JaroWinkler)
		if err == nil && sim > best {
			best = sim
		}
	}

	return best
}
