package array

import (
	"runtime"
	"slices"
	"testing"
)

func expect(t *testing.T, s *Set, want []int) {
	if all := s.All(); !slices.Equal(want, all) {
		t.Errorf("exp %v, got %v", want, all)
	}
}

func TestSet_Add(t *testing.T) {
	s := New()
	expect(t, s, []int{})

	s.Add(1)
	expect(t, s, []int{1})

	s.Add(3)
	s.Add(2)
	expect(t, s, []int{1, 3, 2})
}

func TestSet_AddDuplicates(t *testing.T) {
	s := New()

	s.Add(1)
	s.Add(1)
	s.Add(2)
	s.Add(1)

	expect(t, s, []int{1, 1, 2, 1})
}

func TestSet_AllReturnsSlice(t *testing.T) {
	s := New()
	s.Add(1)
	s.Add(2)

	all1 := s.All()
	all2 := s.All()

	if &all1[0] != &all2[0] {
		t.Error("All() should return the same underlying data")
	}
}

func BenchmarkSet_Add(b *testing.B) {
	for b.Loop() {
		s := New()
		for i := range 1000 {
			s.Add(i)
		}
		runtime.KeepAlive(s)
	}
}
