package search

import "github.com/egonelbre/spexs2/set"

type Sequence struct {
	Index   int
	Section uint32
	Count   uint32
}

type TokenGroup struct {
	Token    Token
	Elems    []Token
	Name     string
	FullName string
}

type TokenInfo struct {
	Token Token
	Name  string
	Count int
}

type Database struct {
	Alphabet map[Token]*TokenInfo
	Groups   map[Token]*TokenGroup

	PosToSequence []Sequence // mapping from position to sequence
	FullSequence  []Token    // concatenated sequences
	Total         []int      // total number sequence for each section
	TotalTokens   []int      // total number of tokens for each section

	Separator string // separator for joining pattern

	nameToToken map[string]Token
	genSeqID    int
	genToken    Token
}

const initialSize = 1e4
const ZeroToken = Token(0)

func NewDatabase() *Database {
	return &Database{
		Alphabet: make(map[Token]*TokenInfo),
		Groups:   make(map[Token]*TokenGroup),

		PosToSequence: make([]Sequence, 0, initialSize),
		FullSequence:  make([]Token, 0, initialSize),
		Total:         make([]int, 0),
		TotalTokens:   make([]int, 0),

		Separator: "",

		nameToToken: make(map[string]Token),
		genSeqID:    0,
		genToken:    Token(1),
	}
}

func (db *Database) AddAllPositions(s set.Set) {
	for i, v := range db.FullSequence {
		if v != ZeroToken {
			s.Add(i)
		}
	}
}

func (db *Database) mkNewToken() Token {
	newToken := db.genToken
	db.genToken++
	return newToken
}

func (db *Database) MatchesOccs(s set.Set) (matches []int, occs []int) {
	matches = make([]int, len(db.Total))
	occs = make([]int, len(db.Total))

	prevseq := -1
	for _, p := range s.All() {
		seq := db.PosToSequence[p]
		occs[seq.Section] += int(seq.Count)
		if seq.Index == prevseq {
			continue
		}
		matches[seq.Section] += int(seq.Count)
		prevseq = seq.Index
	}

	return
}

func (db *Database) AddGroup(group *TokenGroup) Token {
	token := db.mkNewToken()
	group.Token = token
	db.Groups[token] = group
	return token
}

func (db *Database) AddToken(tokenName string) Token {
	token := db.mkNewToken()
	db.nameToToken[tokenName] = token
	db.Alphabet[token] = &TokenInfo{token, tokenName, 0}
	return token
}

func (db *Database) ToTokens(raw []string) []Token {
	tokens := make([]Token, len(raw))
	for i, name := range raw {
		token, ok := db.nameToToken[name]
		if !ok {
			token = db.AddToken(name)
		}
		tokens[i] = token
	}
	return tokens
}

func (db *Database) MakeSection() int {
	db.Total = append(db.Total, 0)
	db.TotalTokens = append(db.TotalTokens, 0)
	return len(db.Total) - 1
}

func (db *Database) addToTokenCount(sec int, tokens []Token, count int) {
	db.TotalTokens[sec] += count * len(tokens)
	for _, t := range tokens {
		db.Alphabet[t].Count += count
	}
}

func (db *Database) AddSequence(sec int, raw []string, count int) {
	db.Total[sec] += count
	tokens := db.ToTokens(raw)
	db.addToTokenCount(sec, tokens, count)

	seq := Sequence{db.genSeqID, uint32(sec), uint32(count)}
	db.genSeqID++

	// add sequence tokens to a single array
	seqstart := len(db.FullSequence)
	db.FullSequence = append(db.FullSequence, tokens...)
	seqend := len(db.FullSequence)

	// add sequence info for each positions
	db.PosToSequence = append(db.PosToSequence, make([]Sequence, len(tokens))...)
	for i := seqstart; i < seqend; i++ {
		db.PosToSequence[i] = seq
	}

	// add separators
	db.FullSequence = append(db.FullSequence, ZeroToken)
	db.PosToSequence = append(db.PosToSequence, seq)
}
