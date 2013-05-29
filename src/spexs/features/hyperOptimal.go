package features
/*
// find optimal hypergeometric split
import (
	"fmt"
	"math"
	"sort"
	. "spexs"
	"stats/hyper"
)

func uniquePositions(q *Query) []uint {
	positions := make([]uint, q.Loc.Len())
	k := 0
	for _, val := range q.Loc.Iter() {
		p, _ := DecodePos(val)
		positions[k] = p
		k += 1
	}
	sort.Sort(uintSlice(positions))
	return uniq(positions)
}

func HyperOptimal(fore []int) Feature {
	return func(q *Query) (float64, string) {
		db := q.Db
		positions := uniquePositions(q)

		allMatches := count(q.Matches(), fore)
		totalSeqs := count(db.Total, fore)

		var best struct {
			p       float64
			matches int
			seqs    int
		}
		best.p = math.Inf(1.0)

		acc := 0
		for _, i := range positions {
			seq := db.Sequences[i]

			for _, sec := range fore {
				acc += seq.Count[sec]
			}

			p := hyper.Split(acc, allMatches, int(i+1), totalSeqs)
			if p < best.p {
				best.p = p
				best.matches = acc
				best.seqs = int(i + 1)
			}
		}

		return best.p, fmt.Sprintf("<%v/%v>", best.matches, best.seqs)
	}
}

type uintSlice []uint

func (p uintSlice) Len() int           { return len(p) }
func (p uintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p uintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func uniq(data []uint) []uint {
	if len(data) <= 0 {
		return data
	}
	k := 1
	for i := 0; i < len(data); i += 1 {
		if data[k-1] != data[i] {
			data[k] = data[i]
			k += 1
		}
	}
	return data[0:k]
}
*/