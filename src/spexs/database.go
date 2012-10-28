package spexs

type Sequence struct {
	Tokens []Token
	Count  map[int]int
}

type Group struct {
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
	Groups   map[Token]Group

	Sequences []Sequence
	Total     []int // total for each section

	Separator string // separator for joining pattern

	nameToToken map[string]Token
	lastToken   Token
}

func NewDatabase(estimatedSize int) *Database {
	return &Database{
		Alphabet: make(map[Token]TokenInfo),
		Groups:   make(map[Token]Group),

		Sequences: make([]Sequence, 0, estimatedSize),
		Total:     make([]int, 0),

		Separator: "",

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

func (db *Database) AddGroup(group Group) Token {
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

func sum(count []int) int {
	total := 0
	for _, val := range count {
		total += val
	}
	return total
}

func (db *Database) MakeSection(seqs [][]string, count []int) int {
	ext := make([]Sequence, len(seqs))
	db.Total = append(db.Total, sum(count))
	sec := len(db.Total) - 1

	for i, raw := range seqs {
		seq := Sequence{db.ToTokens(raw), make(map[int]int)}
		seq.Count[sec] = count[i]
		ext[i] = seq
	}

	db.Sequences = append(db.Sequences, ext...)
	return sec
}
