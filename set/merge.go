package set

import "sort"

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
			i++
			rlast++
		}
		into[i] = lv
		i++
	}
	copy(into[i:], right[rlast:])
}

type bylen [][]int

func (s bylen) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bylen) Len() int           { return len(s) }
func (s bylen) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

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

	sort.Sort(bylen(sets))

	r := make([]int, t)
	prev := sets[0]
	for _, cur := range sets[1:] {
		target := r[t-len(cur)-len(prev):]
		merge(prev, cur, target)
		prev = target
	}

	return r
}
