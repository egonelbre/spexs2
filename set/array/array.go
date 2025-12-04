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

func (s *Set) All() []int {
	return []int(*s)
}
