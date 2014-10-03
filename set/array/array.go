package array

type Set []int

func New() *Set {
	arr := Set(make([]int, 0, 8))
	return &arr
}

func (s *Set) Add(v int) {
	*s = append(*s, v)
}

func (s *Set) Len() int {
	return len(*s)
}

func (s *Set) Unpack() []int {
	return []int(*s)
}

func (s *Set) Iter(fn func(val int)) {
	for _, v := range *s {
		fn(v)
	}
}
