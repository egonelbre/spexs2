package rle

const (
	bits  = 7
	mark  = 1 << bits
	mask  = mark - 1
)

type Set struct {
	cur   int
	count int
	data  []byte
}

func New() *Set {
	s := &Set{0, 0, make([]byte, 0)}
	return s
}

func (s *Set) IsSorted() {}

func (s *Set) Add(v int) {
	df := v - s.cur
	s.cur = v
	s.count += 1
	if df < 0 {
		panic("not in increasing order")
	}
	
	for ; df >= mark; df = df >> bits {
		s.data = append(s.data, byte(mark|(df&mask)))
	}
	s.data = append(s.data, byte(df))
}

func (s *Set) Iter() []int {
	vals := make([]int, s.count)
	j := 0
	base := 0
	df := 0
	k := uint(0)
	for i := 0; i < len(s.data); i += 1 {
		b := s.data[i]
		df = df | (int(b&mask) << k)
		k += 7
		if b < mark {
			base += df
			vals[j] = base
			df = 0
			k = 0
			j += 1
		}
	}
	return vals
}

func (s *Set) Len() int {
	return s.count
}

func (s *Set) AddSet(other *Set) {
	as, bs := s.Iter(), other.Iter()
	
	r := New()
	ai, bi := 0, 0
	for ai < len(as) && bi < len(bs) {
		av, bv := as[ai], bs[bi]
		if av < bv {
			r.Add(av)
			ai += 1
		} else if av > bv {
			r.Add(bv)
			bi += 1
		} else {
			r.Add(av)
			ai += 1
			bi += 1
			
		}
	}

	for ai < len(as) {
		r.Add(as[ai])
		ai += 1
	}

	for bi < len(bs) {
		r.Add(bs[bi])
		bi += 1
	}

	s.count = r.count
	s.data = r.data
	s.cur = r.cur
}

func (s *Set) AddSets(other... *Set) {
	sets := make([][]int, 0)
	if s.Len() > 0 {
		sets = append(sets, s.Iter())
	}
	for _, o := range other {
		if o.Len() > 0 {
			sets = append(sets, o.Iter())
		}
	}

	r := New()
	idx := make([]int, len(sets))
	for len(sets) > 1 {
		min := 10000000000
		for si, i := range idx {
			v := sets[si][i]
			if v < min {
				min = v
			}
		}

		r.Add(min)

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
			r.Add(s[i])
		}
	}

	s.count = r.count
	s.data = r.data
	s.cur = r.cur
}
