package notice

import "regexp"

type Category string

const (
	Admission    Category = "admission"
	Exam         Category = "exam"
	Registration Category = "registration"
	Internship   Category = "internship"
	Scholarship  Category = "scholarship"
	Payment      Category = "payment"
	Holiday      Category = "holiday"
	General      Category = "general"
)

var categoryPatterns = []struct {
	pattern  *regexp.Regexp
	category Category
}{
	{regexp.MustCompile(`(?i)admission\s*(test|result|final|written|slot|viva)`), Admission},
	{regexp.MustCompile(`(?i)(seat\s*plan|exam\s*(schedule|permit)|final[-\s]term\s*exam|set[-\s]b\s*exam)`), Exam},
	{regexp.MustCompile(`(?i)registration`), Registration},
	{regexp.MustCompile(`(?i)internship`), Internship},
	{regexp.MustCompile(`(?i)scholarship`), Scholarship},
	{regexp.MustCompile(`(?i)payment`), Payment},
	{regexp.MustCompile(`(?i)(holiday|eid|vacation|break)`), Holiday},
}

func matchCategory(title string) Category {
	for _, entry := range categoryPatterns {
		if entry.pattern.MatchString(title) {
			return entry.category
		}
	}
	return General
}
