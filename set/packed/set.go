package packed

type bucket uint16

const (
	bits = 15        // bits per bucket
	cbit = 1 << bits // continuation bit
	mask = cbit - 1  // bits mask
)

// Set represents a set of integers.
//
// It's guaranteed that the elements are added in increasing order
// and there are no duplicates.
type Set struct {
	cur   int
	count int
	data  []bucket
}

// New creates a new empty set.
func New() *Set {
	return &Set{}
}

// From creates a set with the provided values.
func From(values ...int) *Set {
	s := New()
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Add adds a new element to the set.
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

// All returns all elements in the set.
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

// Len returns the number of elements in the set.
func (s *Set) Len() int {
	return s.count
}
