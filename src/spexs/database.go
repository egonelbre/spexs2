package spexs

import "bytes"

type Sequence struct {
	Tokens []Token
	Count  map[int]int
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
}

type Database struct {
	Alphabet map[Token]TokenInfo
	Groups   map[Token]TokenGroup

	Sequences []Sequence
	Total     []int // total for each section

	Separator string // separator for joining pattern

	nameToToken map[string]Token
	strToSeq    map[string]int
	lastToken   Token
}

func NewDatabase(estimatedSize int) *Database {
	return &Database{
		Alphabet: make(map[Token]TokenInfo),
		Groups:   make(map[Token]TokenGroup),

		Sequences: make([]Sequence, 0, estimatedSize),
		Total:     make([]int, 0),

		Separator: "",

		strToSeq:    make(map[string]int),
		nameToToken: make(map[string]Token),
		lastToken:   Token(0),
	}
}

func (db *Database) GetToken(seqIdx uint, tokenPos uint) (token Token, ok bool, nextPos uint) {
	seq := &db.Sequences[seqIdx]
	if int(tokenPos) >= len(seq.Tokens) {
		return 0, false, 0
	}
	token = seq.Tokens[tokenPos]
	return token, true, tokenPos + 1
}

func (db *Database) nextToken() Token {
	newToken := db.lastToken
	db.lastToken += 1
	return newToken
}

func (db *Database) AddGroup(group TokenGroup) Token {
	token := db.nextToken()
	group.Token = token
	db.Groups[token] = group
	return token
}

func (db *Database) AddToken(tokenName string) Token {
	token := db.nextToken()
	db.nameToToken[tokenName] = token
	db.Alphabet[token] = TokenInfo{token, tokenName}
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
	return len(db.Total) - 1
}

func sum(count []int) int {
	total := 0
	for _, val := range count {
		total += val
	}
	return total
}

func (db *Database) AddSequences(sec int, seqs [][]string, count []int) {
	db.Total[sec] += sum(count)

	ext := make([]Sequence, len(seqs))
	for i, raw := range seqs {
		seq := Sequence{db.ToTokens(raw), make(map[int]int)}
		seq.Count[sec] = count[i]
		ext[i] = seq
	}

	db.Sequences = append(db.Sequences, ext...)
}

func (db *Database) AddSequence(sec int, raw []string, count int) {
	db.Total[sec] += count
	tokens := db.ToTokens(raw)
	hash := hashTokens(tokens)
	if si, ok := db.strToSeq[hash]; ok {
		db.Sequences[si].Count[sec] += count
		return
	}
	seq := Sequence{tokens, make(map[int]int)}
	seq.Count[sec] = count
	db.Sequences = append(db.Sequences, seq)
	db.strToSeq[hash] = len(db.Sequences) - 1
}

func hashTokens(toks []Token) string {
	var buf bytes.Buffer
	for _, t := range toks {
		buf.WriteRune(rune(t))
	}
	return string(buf.Bytes())
}
