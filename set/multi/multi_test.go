package multi

import (
	"slices"
	"testing"

	"github.com/egonelbre/spexs2/set/array"
	"github.com/egonelbre/spexs2/set/packed"
)

func expect(t *testing.T, s *Set, want []int) {
	if all := s.All(); !slices.Equal(want, all) {
		t.Errorf("exp %v, got %v", want, all)
	}
}

func TestSet_AddSet(t *testing.T) {
	multi := New()
	multi.AddSet(array.From(1, 3))
	expect(t, multi, []int{1, 3})

	multi.AddSet(array.From(2, 5))
	expect(t, multi, []int{1, 2, 3, 5})

	multi.AddSet(array.From(4))
	expect(t, multi, []int{1, 2, 3, 4, 5})
}

func TestSet_AddSet_Mixed(t *testing.T) {
	multi := New()
	multi.AddSet(array.From(1, 3))
	expect(t, multi, []int{1, 3})

	multi.AddSet(packed.From(2, 5))
	expect(t, multi, []int{1, 2, 3, 5})

	multi.AddSet(packed.From(4))
	expect(t, multi, []int{1, 2, 3, 4, 5})

	multi.AddSet(array.From())
	multi.AddSet(packed.From())
	expect(t, multi, []int{1, 2, 3, 4, 5})
}
