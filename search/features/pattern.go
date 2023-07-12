package features

import "github.com/egonelbre/spexs2/search"

// pattern as a string
func Pat() search.Feature {
	return func(q *search.Query) (float64, string) {
		return 0.0, q.String()
	}
}

// pattern as regular expression
func PatRegex() search.Feature {
	return func(q *search.Query) (float64, string) {
		return 0.0, q.StringLong()
	}
}

// length of the pattern
func PatLength() search.Feature {
	return func(q *search.Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			t++
			if e.Flags&search.IsStar != 0 {
				t++
			}
		}
		return float64(t), ""
	}
}

// count of characters
func PatChars() search.Feature {
	return func(q *search.Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			if e.Flags&search.IsGroup == 0 {
				t++
			}
		}
		return float64(t), ""
	}
}

// count of groups
func PatGroups() search.Feature {
	return func(q *search.Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			if e.Flags&search.IsGroup != 0 {
				t++
			}
		}
		return float64(t), ""
	}
}

// count of stars
func PatStars() search.Feature {
	return func(q *search.Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			if e.Flags&search.IsStar != 0 {
				t++
			}
		}
		return float64(t), ""
	}
}
