package features

import . "spexs"

// length of the pattern
func PatLength() FeatureFunc {
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
func PatChars() FeatureFunc {
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
func PatGroups() FeatureFunc {
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
func PatGroups() FeatureFunc {
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
