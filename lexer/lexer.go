package lexer

// import "regexp/syntax
// import "fmt"

// TokenType : the type of a token is string
type TokenType string

// Token : token struct
type Token struct {
	Type TokenType // read program as string
	Val  string    // holds the value of the token
}

// keywords table
var key = map[string]TokenType{
	"print": PRINT,
	"if":    IF,
	"else":  ELSE,
	"while": WHILE,
}

// keyLookup checks the keywords table and return either the keyword or identifier
func keyLookup(ident string) TokenType {
	if tok, ok := key[ident]; ok {
		return tok
	}
	return IDENT
}

// Tokens
const (
	// Punctuation and operators
	LPAR    = "("
	RPAR    = ")"
	LBRAC   = "{"
	RBRAC   = "}"
	COMMA   = ","
	PLUS    = "+"
	MINUS   = "-"
	MULTIP  = "*"
	DIVIDE  = "/"
	MODULO  = "%"
	ASSIGN  = "="
	EQUAL   = "=="
	N_EQUAL = "!="
	AND     = "&&"
	OR      = "||"
	LESS    = "<"
	MORE    = ">"
	LESS_EQ = "<="
	MORE_EQ = ">="
	NEWLINE = "\n"

	// Identifier
	IDENT = "IDENT"
	NUM   = "NUM"

	// Keywords
	PRINT = "PRINT"
	IF    = "IF"
	ELSE  = "ELSE"
	WHILE = "WHILE"

	// Other
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

// Lexer : lexer struct
type Lexer struct {
	input         string //program input
	char          byte   // current char under examination
	position      int    // current position in input (points to current char)
	positionIndex int    // current reading position in input (after current char)
}

// LexConstructor : constructor function of a lexer
func LexConstructor(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.scanChar()
	return lex
}

// scanChar() gives us the next char and moves one step in the input string
func (lex *Lexer) scanChar() {
	if lex.positionIndex >= len(lex.input) {
		lex.char = 0 // 0 = NULL character in ASCII
	} else {
		lex.char = lex.input[lex.positionIndex]
	}
	lex.position = lex.positionIndex
	lex.positionIndex++
}

// NextToken : read next token
func (lex *Lexer) NextToken() Token {
	var tok Token

	// remove all spaces except newline characters
	lex.spaceTrim()

	switch lex.char {

	case '(':
		tok = newToken(LPAR, lex.char)
	case ')':
		tok = newToken(RPAR, lex.char)
	case '{':
		tok = newToken(LBRAC, lex.char)
	case '}':
		tok = newToken(RBRAC, lex.char)
	case ',':
		tok = newToken(COMMA, lex.char)
	case '+':
		tok = newToken(PLUS, lex.char)
	case '-':
		tok = newToken(MINUS, lex.char)
	case '*':
		tok = newToken(MULTIP, lex.char)
	case '/':
		tok = newToken(DIVIDE, lex.char)
	case '%':
		tok = newToken(MODULO, lex.char)
	case '\n':
		tok = newToken(NEWLINE, lex.char)

	// 2-character operators
	case '=':
		if lex.lookAhead() == '=' {
			char := lex.char
			lex.scanChar()
			value := string(char) + string(lex.char)
			tok = Token{Type: EQUAL, Val: value}
		} else {
			tok = newToken(ASSIGN, lex.char)
		}
	case '>':
		if lex.lookAhead() == '=' {
			char := lex.char
			lex.scanChar()
			value := string(char) + string(lex.char)
			tok = Token{Type: MORE_EQ, Val: value}
		} else {
			tok = newToken(MORE, lex.char)
		}
	case '<':
		if lex.lookAhead() == '=' {
			char := lex.char
			lex.scanChar()
			value := string(char) + string(lex.char)
			tok = Token{Type: LESS_EQ, Val: value}
		} else {
			tok = newToken(LESS, lex.char)
		}

	case '!':
		if lex.lookAhead() == '=' {
			char := lex.char
			lex.scanChar()
			value := string(char) + string(lex.char)
			tok = Token{Type: N_EQUAL, Val: value}
		} else {
			tok = newToken(ILLEGAL, lex.char)
		}

	case '&':
		if lex.lookAhead() == '&' {
			char := lex.char
			lex.scanChar()
			value := string(char) + string(lex.char)
			tok = Token{Type: AND, Val: value}
		} else {
			tok = newToken(ILLEGAL, lex.char)
		}

	case '|':
		if lex.lookAhead() == '|' {
			char := lex.char
			lex.scanChar()
			value := string(char) + string(lex.char)
			tok = Token{Type: OR, Val: value}
		} else {
			tok = newToken(ILLEGAL, lex.char)
		}

	case 0:
		tok.Val = ""
		tok.Type = EOF

	default:
		if 'a' <= lex.char && lex.char <= 'z' || 'A' <= lex.char && lex.char <= 'Z' {
			tok.Val = lex.readIdentifier()
			tok.Type = keyLookup(tok.Val)
			return tok
		} else if '0' <= lex.char && lex.char <= '9' {
			tok.Val = lex.readNumber()
			tok.Type = NUM
			return tok
		} else {
			tok = newToken(ILLEGAL, lex.char)
		}
	}

	lex.scanChar()
	return tok
}

// ignore whitespace
func (lex *Lexer) spaceTrim() {
	// Don't ignore newline character - used to indentify the end of assignment
	for lex.char == ' ' || lex.char == '\t' {
		lex.scanChar()
	}
}

func newToken(tokenType TokenType, char byte) Token {
	return Token{Type: tokenType, Val: string(char)}
}

// read ideantifiers
func (lex *Lexer) readIdentifier() string {
	position := lex.position
	for 'a' <= lex.char && lex.char <= 'z' || 'A' <= lex.char && lex.char <= 'Z' {
		lex.scanChar()
	}
	return lex.input[position:lex.position]
}

// read numbers
func (lex *Lexer) readNumber() string {
	position := lex.position
	for '0' <= lex.char && lex.char <= '9' {
		lex.scanChar()
	}
	return lex.input[position:lex.position]
}

// check next character to identify 2-char operators
func (lex *Lexer) lookAhead() byte {
	if lex.positionIndex >= len(lex.input) {
		return 0
	}
	return lex.input[lex.positionIndex]
}
