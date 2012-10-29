package features

import . "spexs"

func Pat() Feature {
	return func(q *Query) (float64, string) {
		return 0.0, q.String()
	}
}

func PatRegex() Feature {
	return func(q *Query) (float64, string) {
		return 0.0, q.StringLong()
	}
}
