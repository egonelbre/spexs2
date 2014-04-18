package features

// find optimal hypergeometric split
import (
	"fmt"
	"math"

	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/stats/hyper"
)

type HyperOptimal struct {
	feature
	Fore []int
}

// the count of matching unique sequences
func countseqs(q *Query, group []int, iter []int) int {
	db := q.Db

	prevseq := -1
	seqs := make([]int, len(db.Total))

	for _, p := range iter {
		seq := db.PosToSequence[p]

		if seq.Index == prevseq {
			continue
		}
		prevseq = seq.Index
		seqs[seq.Section] += 1
	}
	return count(seqs, group)
}

func (f *HyperOptimal) Evaluate(q *Query) (float64, string) {
	iter := q.Loc.Iter()

	allmatches := countseqs(q, f.Fore, iter)
	totalSeqs := count(f.Db.Total, f.Fore)

	var best struct {
		p       float64
		matches int
		seqs    int
	}
	best.p = math.Inf(1.0)

	include := make([]bool, len(f.Db.Total))
	for _, sec := range f.Fore {
		include[sec] = true
	}

	prevseq := -1

	count := 0
	for _, i := range iter {
		seq := f.Db.PosToSequence[i]
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

	return best.p, fmt.Sprintf("<%v/%v>", best.matches, best.seqs)
}
