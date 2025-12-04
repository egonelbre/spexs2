package packed

import (
	"slices"
	"testing"
)

func expect(t *testing.T, s *Set, want []int) {
	if all := s.All(); !slices.Equal(want, all) {
		t.Errorf("exp %v, got %v", want, all)
	}
}

func test(t *testing.T, values []int) {
	s := New()
	for _, v := range values {
		s.Add(v)
	}
	expect(t, s, values)
}

func TestSet(t *testing.T) {
	t.Run("Simple", func(t *testing.T) { test(t, []int{1}) })
	t.Run("Increasing", func(t *testing.T) { test(t, []int{1, 5, 10, 20, 50}) })
	t.Run("SmallDeltas", func(t *testing.T) { test(t, []int{1, 2, 3, 4, 5}) })
	t.Run("LargeDeltas", func(t *testing.T) { test(t, []int{1, 100000, 200000, 300000}) })
	t.Run("VeryLargeDeltas", func(t *testing.T) { test(t, []int{1, 10000000, 20000000, 100000000}) })
	t.Run("Sparse", func(t *testing.T) { test(t, []int{1, 1000, 2000, 10000, 20000, 100000}) })
}

func TestSet_AddPanicOnNonIncreasing(t *testing.T) {
	s := New()
	s.Add(5)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Add() should panic when adding non-increasing value")
		}
	}()

	s.Add(3)
}

func TestSet_AddPanicOnEqual(t *testing.T) {
	s := New()
	s.Add(5)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Add() should panic when adding equal value")
		}
	}()

	s.Add(5)
}
