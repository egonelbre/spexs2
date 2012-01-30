package spexs

type Pos uint64

const (
	PATTERN_LENGTH_BITS = 4
	PATTERN_LENGTH_MASK = (2 << PATTERN_LENGTH_BITS) - 1
	EmptyPos = 0
)
// meaning pattern length can be at most = 2^4
// and can there can be at most 2^(64 - 4) patterns

// pos must be < 16
func PosEncode(idx int, pos byte) Pos {
	return idx << PATTERN_LENGTH_BITS | (pos & PATTERN_LENGTH_MASK)
}

func PosDecode( p Pos ) (idx int, pos byte ) {
	pos = p & PATTERN_LENGTH_MASK
	idx = p >> PATTERN_LENGTH_BITS
}

type ItemFunc func(item Pos)

type Set interface {
  Add(item Pos)
  Contains(item Pos) bool
  Length() uint
  Iterate( ItemFunc )
}

type HashSet struct {
  data map[Pos] bool
}

func NewHashSet() *HashSet {
  return &HashSet{make(map[Pos] bool), 0}
}

func ( hs *HashSet ) Add(val Pos) {
	hs.data[val] = bool
}

func ( hs *HashSet ) Contains(val Pos) bool {
  _, exists := hs.data[val]
  return exists
}

func ( hs *HashSet) Length() uint {
  return len(hs.data)
}

func ( hs *HashSet) Iterate( f ItemFunc ){
	for v, _ := range hs.data {
		f(v)
	}
}

// result can't be hs nor gs
func ( h * Set) And( g *Set, result *Set ){
	var first, second Set

	if h.Length() < g.Length() {
		first = h
		second = g
	} else {
		first = g
		second = h
	}
	
	first.Iterate( func(item Pos){
		if second.Contains(item) {
			result.Add(item)
		}
	})
}

func ( h * Set) Or( g *Set, result *Set ){
	if h != result {
		h.Iterate( func(item Pos){ result.Add(item) } )
	}
	if g != result {
		g.Iterate( func(item Pos){ result.Add(item) } )
	}
}

func ( h * Set) AddSet( g *Set ){
	g.Iterate( func(item Pos){ h.Add(item) } )
}

type FullSet struct {
	Ref *UnicodeReference
	Count int
}

func NewFullSet(ref *UnicodeReference) *FullSet{
	f := &FullSet{ ref }
	for p, _ := range ref.Pats {
		f.Count += p.Count
	}
	return f
}

func ( f * FullSet ) Add( val Pos ){ }

func ( f * FullSet ) Contains( val Pos ){
	idx, pos := PosDecode(val)
	return idx < len(f.Pats) && pos < len( f.Pats[idx].Pat )
}

func ( f *FullSet) Length() uint {
  return f.Count
}

func ( f *FullSet) Iterate( f ItemFunc ){
	for idx := range len(f.Ref.Pats) {
		pat := f.Ref.Pats[idx]
		/* TODO: iterate properly utf8 style */
		for pos := range pat.Length {
			f( PosEncode(idx, pos) )
		}
	}
}