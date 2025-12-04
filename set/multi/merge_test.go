package multi

import (
	"runtime"
	"slices"
	"testing"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name  string
		left  []int
		right []int
		want  []int
	}{
		{
			name:  "both empty",
			left:  []int{},
			right: []int{},
			want:  []int{},
		},
		{
			name:  "left empty",
			left:  []int{},
			right: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "right empty",
			left:  []int{1, 2, 3},
			right: []int{},
			want:  []int{1, 2, 3},
		},
		{
			name:  "no overlap",
			left:  []int{1, 2, 3},
			right: []int{4, 5, 6},
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "interleaved",
			left:  []int{1, 3, 5},
			right: []int{2, 4, 6},
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "with duplicates",
			left:  []int{1, 2, 3},
			right: []int{2, 3, 4},
			want:  []int{1, 2, 2, 3, 3, 4},
		},
		{
			name:  "single element each",
			left:  []int{1},
			right: []int{2},
			want:  []int{1, 2},
		},
		{
			name:  "right before left",
			left:  []int{5, 6, 7},
			right: []int{1, 2, 3},
			want:  []int{1, 2, 3, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := make([]int, len(tt.left)+len(tt.right))
			merge(tt.left, tt.right, dst)
			if !slices.Equal(dst, tt.want) {
				t.Errorf("got %v, want %v", dst, tt.want)
			}
		})
	}
}

func TestMergeSortedInts(t *testing.T) {
	tests := []struct {
		name string
		sets [][]int
		want []int
	}{
		{
			name: "no sets",
			sets: [][]int{},
			want: nil,
		},
		{
			name: "single set",
			sets: [][]int{{1, 2, 3}},
			want: []int{1, 2, 3},
		},
		{
			name: "two sets no overlap",
			sets: [][]int{{1, 2, 3}, {4, 5, 6}},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "two sets interleaved",
			sets: [][]int{{1, 3, 5}, {2, 4, 6}},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "three sets",
			sets: [][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "with duplicates",
			sets: [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}},
			want: []int{1, 2, 2, 3, 3, 3, 4, 4, 5},
		},
		{
			name: "different sizes",
			sets: [][]int{{1}, {2, 3, 4, 5}, {6, 7}},
			want: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "with empty sets",
			sets: [][]int{{1, 2}, {}, {3, 4}},
			want: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeSortedInts(tt.sets...)
			if !slices.Equal(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByLen(t *testing.T) {
	sets := [][]int{
		{1, 2, 3, 4, 5},
		{1},
		{1, 2, 3},
		{1, 2},
	}

	bl := bylen(sets)

	if bl.Len() != 4 {
		t.Errorf("Len() = %d, want 4", bl.Len())
	}

	if !bl.Less(1, 0) {
		t.Error("Less(1, 0) should be true (len 1 < len 5)")
	}

	if bl.Less(0, 1) {
		t.Error("Less(0, 1) should be false (len 5 > len 1)")
	}

	bl.Swap(0, 1)
	if len(sets[0]) != 1 || len(sets[1]) != 5 {
		t.Error("Swap didn't work correctly")
	}
}

func BenchmarkMerge(b *testing.B) {
	left := make([]int, 1000)
	right := make([]int, 1000)
	for i := range left {
		left[i] = i * 2
		right[i] = i*2 + 1
	}
	dst := make([]int, 2000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merge(left, right, dst)
	}
}

func BenchmarkMergeSortedInts(b *testing.B) {
	sets := make([][]int, 10)
	for i := range sets {
		sets[i] = make([]int, 100)
		for j := range sets[i] {
			sets[i][j] = i + j*10
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(mergeSortedInts(sets...))
	}
}
