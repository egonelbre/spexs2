package set

import "math"

func MergeSortedInts(sets ...[]int) []int {
	t := 0
	for _, s := range sets {
		t += len(s)
	}
	r := make([]int, t)
	ri := 0

	idx := make([]int, len(sets))
	for len(sets) > 1 {
		min := math.MaxInt64
		for si, i := range idx {
			v := sets[si][i]
			if v < min {
				min = v
			}
		}

		r[ri] = min
		ri += 1

		for si := len(idx) - 1; si >= 0; si -= 1 {
			i := idx[si]
			if sets[si][i] == min {
				i += 1
				if i >= len(sets[si]) {
					idx = append(idx[:si], idx[si+1:]...)
					sets = append(sets[:si], sets[si+1:]...)
				} else {
					idx[si] = i
				}
			}
		}
	}

	if len(sets) == 1 {
		s := sets[0]
		for i := idx[0]; i < len(s); i += 1 {
			r[ri] = s[i]
			ri += 1
		}
	}

	return r
}

func MergeSortedUniqueInts(sets ...[]int) []int {
	t := 0
	for _, s := range sets {
		t += len(s)
	}
	r := make([]int, t)
	ri := 0

	idx := make([]int, len(sets))
	for len(sets) > 1 {
		min := math.MaxInt64
		min_si := -1
		for si, i := range idx {
			v := sets[si][i]
			if v < min {
				min = v
				min_si = si
			}
		}

		r[ri] = min
		ri += 1

		si := min_si
		idx[si] += 1
		if idx[si] >= len(sets[si]) {
			idx = append(idx[:si], idx[si+1:]...)
			sets = append(sets[:si], sets[si+1:]...)
		}
	}

	if len(sets) == 1 {
		s := sets[0]
		for i := idx[0]; i < len(s); i += 1 {
			r[ri] = s[i]
			ri += 1
		}
	}

	return r
}
