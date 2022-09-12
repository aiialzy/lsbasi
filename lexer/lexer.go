package lexer

import (
	"fmt"
	"lsbasi/token"
	"strings"
	"unicode"
)

var reserved_keywords = map[string]*token.Token{
	"BEGIN":   token.New(token.BEGIN, "BEGIN"),
	"END":     token.New(token.END, "END"),
	"PROGRAM": token.New(token.PROGRAM, "PROGRAM"),
	"VAR":     token.New(token.VAR, "VAR"),
	"DIV":     token.New(token.INTEGER_DIV, "DIV"),
	"INTEGER": token.New(token.INTEGER, "INTEGER"),
	"REAL":    token.New(token.REAL, "REAL"),
}

type Lexer struct {
	text        []rune
	pos         int
	currentChar rune
}

func (l *Lexer) number() *token.Token {
	var number strings.Builder
	for unicode.IsDigit(l.currentChar) ||
		l.currentChar == 'x' ||
		l.currentChar == 'X' ||
		l.currentChar == 'o' ||
		l.currentChar == 'O' ||
		l.currentChar == 'b' ||
		l.currentChar == 'B' ||
		l.currentChar == '.' {
		number.WriteRune(l.currentChar)
		l.advance()
	}

	numStr := number.String()
	if strings.Contains(numStr, ".") {
		return token.New(token.FLOAT_CONST, numStr)
	} else {
		return token.New(token.INTEGER_CONST, numStr)
	}
}

func (l *Lexer) word() string {
	var word strings.Builder
	for unicode.IsLetter(l.currentChar) ||
		unicode.IsDigit(l.currentChar) ||
		l.currentChar == '_' {

		word.WriteRune(l.currentChar)
		l.advance()
	}

	return word.String()
}

func (l *Lexer) peek() rune {
	peekPos := l.pos + 1
	if peekPos < len(l.text) {
		return l.text[peekPos]
	}
	return 0
}

func (l *Lexer) error() {
	panic("unexpected input")
}

func (l *Lexer) skipWhiteSpace() {
	for unicode.IsSpace(l.currentChar) {
		l.advance()
	}
}

func (l *Lexer) advance() {
	l.pos += 1
	if l.pos < len(l.text) {
		l.currentChar = l.text[l.pos]
	} else {
		l.currentChar = 0
	}
}

func (l *Lexer) skipComment() {
	for l.currentChar != '}' {
		l.advance()
	}
	l.advance()
}

func (l *Lexer) GetNextToken() *token.Token {
	t := l.nextToken()
	// fmt.Println(t)

	return t
}

func (l *Lexer) nextToken() *token.Token {

	for unicode.IsSpace(l.currentChar) || l.currentChar == '{' {
		if unicode.IsSpace(l.currentChar) {
			l.skipWhiteSpace()
		}

		if l.currentChar == '{' {
			l.skipComment()
		}
	}

	if l.currentChar == 0 {
		return token.New(token.EOF, "")
	}

	if unicode.IsDigit(l.currentChar) {
		return l.number()
	}

	if l.currentChar == '_' || unicode.IsLetter(l.currentChar) {
		word := l.word()
		word = strings.ToUpper(word)
		if t, exists := reserved_keywords[word]; exists {
			return t
		}
		return token.New(token.ID, word)
	}

	if l.currentChar == ':' && l.peek() == '=' {
		l.advance()
		l.advance()
		return token.New(token.ASSIGN, ":=")
	}

	switch l.currentChar {
	case '+':
		l.advance()
		return token.New(token.PLUS, "+")

	case '-':
		l.advance()
		return token.New(token.SUB, "-")

	case '*':
		l.advance()
		return token.New(token.MUL, "*")

	case '/':
		l.advance()
		return token.New(token.FLOAT_DIV, "/")

	case '%':
		l.advance()
		return token.New(token.MOD, "%")

	case '(':
		l.advance()
		return token.New(token.LPAREN, "(")

	case ')':
		l.advance()
		return token.New(token.RPAREN, ")")

	case ';':
		l.advance()
		return token.New(token.SEMI, ";")

	case '.':
		l.advance()
		return token.New(token.DOT, ".")

	case ':':
		l.advance()
		return token.New(token.COLON, ":")

	case ',':
		l.advance()
		return token.New(token.COMMA, ",")
	}

	fmt.Printf("%v\n", string(l.text[l.pos:]))
	l.error()

	return nil
}

func New(text []rune) *Lexer {
	if len(text) == 0 {
		panic("text length must bigger than 0")
	}
	return &Lexer{
		text:        text,
		pos:         0,
		currentChar: text[0],
	}
}
