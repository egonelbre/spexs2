package set

func merge(left, right, into []int) {
	if len(left) == 0 {
		copy(into, right)
		return
	}
	if len(right) == 0 {
		copy(into, left)
		return
	}

	i, rlast := 0, 0
	for _, lv := range left {
		for _, rv := range right[rlast:] {
			if lv < rv {
				break
			}
			into[i] = rv
			i += 1
			rlast += 1
		}
		into[i] = lv
		i += 1
	}
	copy(into[i:], right[rlast:])
}

func MergeSortedInts(sets ...[]int) []int {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0]
	}

	t := 0
	for _, s := range sets {
		t += len(s)
	}

	r := make([]int, t)
	prev := sets[0]
	for i := 1; i < len(sets); i += 1 {
		target := r[t-len(sets[i])-len(prev):]
		merge(prev, sets[i], target)
		prev = target
	}

	return r
}
