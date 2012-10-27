package spexs

type Sequence struct {
	Tokens  []Token
	Len     int
	Section int
	Count   int
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

type Section struct {
	Count int
}

type Database struct {
	Alphabet map[Token]TokenInfo
	Groups   map[Token]Group

	Sequences []Sequence
	Sections  []Section

	Separator string // separator for joining pattern

	nameToToken map[string]Token
	lastToken   Token
}

func NewDatabase(estimatedSize int) *Database {
	return &Database{
		Alphabet: make(map[Token]TokenInfo),
		Groups:   make(map[Token]Group),

		Sequences: make([]Sequence, 0, estimatedSize),
		Sections:  make([]Section, 0, 2),

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

func (db *Database) AddSequence(seq Sequence) {
	db.Sequences = append(db.Sequences, seq)
	if seq.Section >= len(db.Sections) {
		df := seq.Section - len(db.Sections) + 1
		extension := make([]Section, df)
		db.Sections = append(db.Sections, extension...)
	}
	db.Sections[seq.Section].Count += 1
}

func (db *Database) ToTokens(tokenNames []string) []Token {
	tokens := make([]Token, len(tokenNames))
	for i, name := range tokenNames {
		token, ok := db.nameToToken[name]
		if !ok {
			token = db.AddToken(name)
		}
		tokens[i] = token
	}
	return tokens
}
