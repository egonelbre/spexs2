package spexs

type Sequence struct {
	Tokens  []Tid
	Len     int
	Section int
	Count   int
}

type Group struct {
	Id    Tid
	Elems []Tid
	Str   string
}

type Token struct {
	Id  Tid
	Str string
}

type Section struct {
	Sequences int
}

type Database struct {
	Alphabet map[Tid]Token
	Groups   map[Tid]Group

	Sequences []Sequence
	Sections  []Section

	lastId int
}

func NewDatabase(estimatedSize int) *Database {
	return &Database{
		Alphabet: make(map[Tid]Token),
		Groups:   make(map[Tid]Group),

		Sequences: make([]Sequence, 0, estimatedSize),
		Sections:  make([]Section, 0, 2),

		lastId: 0,
	}
}

func (db *Database) GetToken(seqIdx int, tokenPos int) (token Tid, ok bool, nextPos int) {
	seq := &db.Sequences[seqIdx]
	if int(tokenPos) >= len(seq.Tokens) {
		return 0, false, 0
	}

	rune, width := seq.Tokens[tokenPos], 1
	return Tid(rune), true, tokenPos + width
}

func (db *Database) AddGroup(group Group) {
	db.Groups[group.Id] = group
}

func (db *Database) AddSequence(seq Sequence) {
	db.Sequences = append(db.Sequences, seq)
	db.Sections[seq.Section].Sequences += 1
}
