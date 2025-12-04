package packed

type bucket uint16

const (
	bits = 15        // bits per bucket
	cbit = 1 << bits // continuation bit
	mask = cbit - 1  // bits mask
)

type Set struct {
	cur   int
	count int
	data  []bucket
}

func New() *Set {
	return &Set{}
}

func (s *Set) Add(v int) {
	df := v - s.cur

	if df <= 0 {
		panic("Not in increasing order!")
	}

	s.cur = v
	s.count++

	// fast detection of numbits
	switch {
	case df < cbit:
		s.data = append(s.data, bucket(df))

	case df>>(bits*1) < cbit:
		s.data = append(s.data,
			bucket(cbit|(df&mask)),
			bucket(df>>bits),
		)
	case df>>(bits*2) < cbit:
		s.data = append(s.data,
			bucket(cbit|(df&mask)),
			bucket(cbit|((df>>(bits*1))&mask)),
			bucket(df>>(bits*2)),
		)
	case df>>(bits*3) < cbit:
		s.data = append(s.data,
			bucket(cbit|(df&mask)),
			bucket(cbit|((df>>(bits*1))&mask)),
			bucket(cbit|((df>>(bits*2))&mask)),
			bucket(df>>(bits*3)),
		)
	case df>>(bits*4) < cbit:
		s.data = append(s.data,
			bucket(cbit|(df&mask)),
			bucket(cbit|((df>>(bits*1))&mask)),
			bucket(cbit|((df>>(bits*2))&mask)),
			bucket(cbit|((df>>(bits*3))&mask)),
			bucket(df>>(bits*4)),
		)
	default:
		for ; df >= cbit; df >>= bits {
			s.data = append(s.data, bucket(cbit|(df&mask)))
		}
		s.data = append(s.data, bucket(df))
	}

}

func (s *Set) All() []int {
	vals := make([]int, s.count)
	j := 0
	base := 0
	df := 0
	k := uint(0)
	for _, b := range s.data {
		df |= int(b&mask) << k
		k += bits
		if b < cbit {
			base += df
			vals[j] = base
			df = 0
			k = 0
			j++
		}
	}
	return vals
}

func (s *Set) Len() int {
	return s.count
}
