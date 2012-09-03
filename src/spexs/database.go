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
	Alias string
	Str   string
}

type Token struct {
	Id  Tid
	Str string
}

type Section struct {
	Count int
}

type Database struct {
	Alphabet map[Tid]Token
	Groups   map[Tid]Group

	Sequences []Sequence
	Sections  []Section

	tokenToId map[string]Tid
	lastId    Tid
}

func NewDatabase(estimatedSize int) *Database {
	return &Database{
		Alphabet: make(map[Tid]Token),
		Groups:   make(map[Tid]Group),

		Sequences: make([]Sequence, 0, estimatedSize),
		Sections:  make([]Section, 0, 2),

		tokenToId: make(map[string]Tid),
		lastId:    Tid(0),
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

func (db *Database) nextId() Tid {
	id := db.lastId
	db.lastId += 1
	return id
}

func (db *Database) AddGroup(group Group) Tid {
	tid := db.nextId()
	group.Id = tid
	db.Groups[tid] = group
	return tid
}

func (db *Database) AddToken(token string) Tid {
	tid := db.nextId()
	db.tokenToId[token] = tid
	db.Alphabet[tid] = Token{tid, token}
	return tid
}

func (db *Database) AddSequence(seq Sequence) {
	db.Sequences = append(db.Sequences, seq)
	db.Sections[seq.Section].Count += 1
}

func (db *Database) ToTids(tokens []string) []Tid {
	tids := make([]Tid, len(tokens))
	for i, token := range tokens {
		tid, ok := db.tokenToId[token]
		if !ok {
			tid = db.AddToken(token)
		}
		tids[i] = tid
	}
	return tids
}
