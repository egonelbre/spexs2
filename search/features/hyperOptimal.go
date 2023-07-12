package features

// find optimal hypergeometric split
import (
	"fmt"
	"math"

	"github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/stats/hyper"
)

// the count of matching unique sequences
func countseqs(q *search.Query, group []int, iter []int) int {
	db := q.Db

	prevseq := -1
	seqs := make([]int, len(db.Total))

	for _, p := range iter {
		seq := db.PosToSequence[p]

		if seq.Index == prevseq {
			continue
		}
		prevseq = seq.Index
		seqs[seq.Section]++
	}
	return count(seqs, group)
}

func HyperOptimal(fore []int) search.Feature {
	return func(q *search.Query) (float64, string) {
		db := q.Db

		iter := q.Loc.Iter()

		allmatches := countseqs(q, fore, iter)
		totalSeqs := count(db.Total, fore)

		var best struct {
			p       float64
			matches int
			seqs    int
		}
		best.p = math.Inf(1.0)

		include := make([]bool, len(db.Total))
		for _, sec := range fore {
			include[sec] = true
		}

		prevseq := -1

		count := 0
		for _, i := range iter {
			seq := db.PosToSequence[i]
			if seq.Index == prevseq {
				continue
			}
			prevseq = seq.Index

			if !include[seq.Section] {
				continue
			}
			count += int(seq.Count)

			p := hyper.ComplementCdfSlow(count, allmatches, seq.Index+1, totalSeqs)
			if p < best.p {
				best.p = p
				best.matches = count
				best.seqs = seq.Index + 1
			}
		}

		//fmt.Printf("<%v/%v> <%v/%v>\t%v\t%v\n", best.matches, best.seqs, allMatches, totalSeqs, best.p, q)

		return best.p, fmt.Sprintf("<%v/%v>", best.matches, best.seqs)
	}
}
