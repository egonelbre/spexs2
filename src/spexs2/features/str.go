package features

import . "spexs"

func Pat() FeatureFunc {
	return func(q *Query) (float64, string) {
		return 0.0, q.String()
	}
}

func PatRegex() FeatureFunc {
	return func(q *Query) (float64, string) {
		return 0.0, q.StringLong()
	}
}
