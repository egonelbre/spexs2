package spexs

type Pos uint64

const (
	PATTERN_LENGTH_BITS = 5
	PATTERN_LENGTH_MASK = (1 << PATTERN_LENGTH_BITS) - 1
	EmptyPos            = 0
)

// meaning pattern length can be at most = 2^4-1
// and can there can be at most 2^(64 - 4) patterns

// pos must be < 16
func PosEncode(idx int, pos byte) Pos {
	idxe := uint64(idx << PATTERN_LENGTH_BITS)
	pose := uint64(pos & PATTERN_LENGTH_MASK)
	return Pos(idxe | pose)
}

func PosDecode(p Pos) (idx int, pos byte) {
	pos = byte(p & PATTERN_LENGTH_MASK)
	idx = int(p >> PATTERN_LENGTH_BITS)
	return
}

type Set interface {
	Add(val Pos)
	Contains(val Pos) bool
	Length() int
	Iter() chan Pos
}

type HashSet struct {
	data map[Pos]bool
}

func NewHashSet() *HashSet {
	return &HashSet{make(map[Pos]bool)}
}

func (hs *HashSet) Add(val Pos) {
	hs.data[val] = true
}

func (hs *HashSet) Contains(val Pos) bool {
	_, exists := hs.data[val]
	return exists
}

func (hs *HashSet) Length() int {
	return len(hs.data)
}

func (hs *HashSet) Iter() chan Pos {
	ch := make(chan Pos, 100)
	go func(){
		for v, _ := range hs.data {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// result can't be hs nor gs
func SetAnd(h Set, g Set, result Set) {
	var first, second Set

	if h.Length() < g.Length() {
		first = h
		second = g
	} else {
		first = g
		second = h
	}

	for item := range first.Iter() {
		if second.Contains(item) {
			result.Add(item)
		}
	}
}

func SetOr(h Set, g Set, result Set) {
	if h != result {
		for item := range h.Iter() {
			result.Add(item)
		}
	}
	if g != result {
		for item := range g.Iter() {
			result.Add(item)
		}
	}
}

func SetAddSet(h Set, g Set) {
	for item := range g.Iter() {
		h.Add(item)
	}
}

type FullSet struct {
	Ref   *UnicodeReference
	Count int
}

func NewFullSet(ref *UnicodeReference) *FullSet {
	f := &FullSet{ref, 0}
	for _, p := range ref.Pats {
		f.Count += p.Count
	}
	return f
}

func (f *FullSet) Add(val Pos) {}

func (f *FullSet) Contains(val Pos) bool {
	idx, pos := PosDecode(val)
	return idx < len(f.Ref.Pats) && int(pos) < len(f.Ref.Pats[idx].Pat)
}

func (f *FullSet) Length() int {
	return f.Count
}

func (f *FullSet) Iter() chan Pos {
	ch := make(chan Pos, 100)
	go func(){
		for idx, pat := range f.Ref.Pats {
			/* TODO: iterate properly utf8 style */
			for pos, _ := range pat.Pat {
				ch <- PosEncode(idx, uint8(pos))
			}
		}
		close(ch)
	}()
	return ch
}
