package benchmark

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/search/extenders"
	"github.com/egonelbre/spexs2/search/features"
	"github.com/egonelbre/spexs2/search/filters"
	"github.com/egonelbre/spexs2/search/pool"
)

var (
	limit = flag.Int("limit", -1, "limit the number of regexps to expand")
	large = flag.Bool("large", false, "run large tests (4GB)")
	huge  = flag.Bool("huge", false, "make tests larger than (4GB)")
)

func ScanLines(filename string, line func(string)) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	buf := bufio.NewReaderSize(file, 128<<10)
	sc := bufio.NewScanner(buf)
	for sc.Scan() {
		line(sc.Text())
	}
}

func LoadDatabase(foreground, background string) *search.Database {
	db := search.NewDatabase()

	const proteins = "ACDEFGHIKLMNPRQSTVWY"
	var tokens []search.Token
	for _, token := range proteins {
		tokens = append(tokens, db.AddToken(string(token)))
	}

	db.AddGroup(&search.TokenGroup{
		Elems:    tokens,
		Name:     ".",
		FullName: ".",
	})

	fore := db.MakeSection()
	ScanLines(foreground, func(line string) {
		line = strings.TrimSpace(line)
		db.AddSequence(fore, strings.Split(line, ""), 1)
	})

	back := db.MakeSection()
	ScanLines(background, func(line string) {
		line = strings.TrimSpace(line)
		db.AddSequence(back, strings.Split(line, ""), 1)
	})

	return db
}

func NewTestSetup(db *search.Database) *search.Setup {
	s := &search.Setup{}

	s.Db = db

	var (
		ForegroundSet = []int{0}
		BackgroundSet = []int{1}
	)

	// features
	groups := features.PatGroups()
	length := features.PatLength()
	matches := features.Matches(ForegroundSet)
	binom := features.Binom(ForegroundSet, BackgroundSet)
	ratio := features.MatchesRatio(ForegroundSet, BackgroundSet)

	order := []search.Feature{
		func(q *search.Query) (float64, string) {
			v, info := binom(q)
			return -v, info
		},
		ratio,
	}

	s.In = pool.NewStack()
	s.Out = pool.NewPriority(order, 100)

	if !*huge {
		s.Extender = extenders.Regex

		s.Extendable = filters.Compose(
			filters.FromFeature(groups, []byte(`{"max": 3}`)),
			filters.FromFeature(length, []byte(`{"max": 6}`)),
			filters.FromFeature(matches, []byte(`{"min": 20}`)),
			filters.NoStartingGroup(s, nil),
		)
		s.Outputtable = filters.Compose(
			filters.FromFeature(ratio, []byte(`{"min": 2}`)),
			filters.FromFeature(binom, []byte(`{"max": 1e-3}`)),
			filters.NoEndingGroup(s, nil),
		)
	} else {
		s.Extender = extenders.Regex

		s.Extendable = filters.Compose(
			filters.FromFeature(groups, []byte(`{"max": 6}`)),
			filters.FromFeature(length, []byte(`{"max": 6}`)),
			filters.FromFeature(matches, []byte(`{"min": 1}`)),
			filters.NoStartingGroup(s, nil),
		)
		s.Outputtable = filters.Compose(
			filters.FromFeature(ratio, []byte(`{"min": 2}`)),
			filters.FromFeature(binom, []byte(`{"max": 1e-3}`)),
			filters.NoEndingGroup(s, nil),
		)
	}

	s.PreProcess = func(q *search.Query) error { return nil }

	var limiter int64
	s.PostProcess = func(q *search.Query) error {
		if *limit > 0 {
			if atomic.AddInt64(&limiter, 1) > int64(*limit) {
				return errors.New("early exit")
			}
		}
		return nil
	}

	return s
}

type Case struct {
	Name       string
	Foreground string
	Background string
	Large      bool
}

var cases = []*Case{
	&Case{"10k", "proteins_10k.inp", "proteins_10k.ref", false},
	&Case{"30k", "proteins_30k.inp", "proteins_30k.ref", false},
	&Case{"large", "proteins_large.inp", "proteins_large.ref", true},
}

func BenchmarkRun(b *testing.B) {
	if *huge {
		*large = true
	}
	for _, c := range cases {
		if c.Large && !*large {
			continue
		}

		b.Run(c.Name, func(b *testing.B) {
			db := LoadDatabase(c.Foreground, c.Background)
			for _, nprocs := range []int{0, 1, 2, 4, 8, 16, 32} {
				b.Run(strconv.Itoa(nprocs), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						s := NewTestSetup(db)
						if nprocs == 0 {
							search.Run(s)
						} else {
							search.RunParallel(s, nprocs)
						}
					}
				})
			}
		})
	}
}
