package features

import . "github.com/egonelbre/spexs2/search"

// pattern as a string
func Pat() Feature {
	return func(q *Query) (float64, string) {
		return 0.0, q.String()
	}
}

// pattern as regular expression
func PatRegex() Feature {
	return func(q *Query) (float64, string) {
		return 0.0, q.StringLong()
	}
}

// length of the pattern
func PatLength() Feature {
	return func(q *Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			t += 1
			if e.IsStar {
				t += 1
			}
		}
		return float64(t), ""
	}
}

// count of characters
func PatChars() Feature {
	return func(q *Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			if !e.IsGroup {
				t += 1
			}
		}
		return float64(t), ""
	}
}

// count of groups
func PatGroups() Feature {
	return func(q *Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			if e.IsGroup {
				t += 1
			}
		}
		return float64(t), ""
	}
}

// count of stars
func PatStars() Feature {
	return func(q *Query) (float64, string) {
		t := 0
		for _, e := range q.Pat {
			if e.IsStar {
				t += 1
			}
		}
		return float64(t), ""
	}
}
