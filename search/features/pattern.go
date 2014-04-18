package features

import . "github.com/egonelbre/spexs2/search"

// pattern as a string
type Pat struct {
	feature
}

func (f *Pat) Evaluate(q *Query) (float64, string) {
	return 0.0, q.String(f.Db)
}

// pattern as regular expression
type PatRegex struct {
	feature
}

func (f *PatRegex) Evaluate(q *Query) (float64, string) {
	return 0.0, q.StringLong(f.Db)
}

// length of the pattern
type PatLength struct {
	feature
}

func (f *PatLength) Evaluate(q *Query) (float64, string) {
	t := 0
	for _, e := range q.Pat {
		t += 1
		if e.Flags&IsStar != 0 {
			t += 1
		}
	}
	return float64(t), ""
}

// count of characters
type PatChars struct {
	feature
}

func (f *PatChars) Evaluate(q *Query) (float64, string) {
	t := 0
	for _, e := range q.Pat {
		if e.Flags&IsGroup == 0 {
			t += 1
		}
	}
	return float64(t), ""
}

// count of groups
type PatGroups struct {
	feature
}

func (f *PatGroups) Evaluate(q *Query) (float64, string) {
	t := 0
	for _, e := range q.Pat {
		if e.Flags&IsGroup != 0 {
			t += 1
		}
	}
	return float64(t), ""
}

// count of stars
type PatStars struct {
	feature
}

func (f *PatStars) Evaluate(q *Query) (float64, string) {
	t := 0
	for _, e := range q.Pat {
		if e.Flags&IsStar != 0 {
			t += 1
		}
	}
	return float64(t), ""
}
