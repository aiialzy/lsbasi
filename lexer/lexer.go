package lexer

import (
	"lsbasi/token"
	"strings"
	"unicode"
)

type Lexer struct {
	text        []rune
	pos         int
	currentChar rune
}

func (l *Lexer) integer() string {
	var number strings.Builder
	for unicode.IsDigit(l.currentChar) {
		number.WriteRune(l.currentChar)
		l.advance()
	}

	return number.String()
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

func (l *Lexer) GetNextToken() *token.Token {

	if unicode.IsSpace(l.currentChar) {
		l.skipWhiteSpace()
	}

	if l.currentChar == 0 {
		return token.New(token.EOF, "")
	}

	if unicode.IsDigit(l.currentChar) {
		number := l.integer()
		return token.New(token.INTEGER, number)
	}

	switch l.currentChar {
	case '+':
		l.advance()
		return token.New(token.PLUS, string(l.currentChar))

	case '-':
		l.advance()
		return token.New(token.SUB, string(l.currentChar))

	case '*':
		l.advance()
		return token.New(token.MUL, string(l.currentChar))

	case '/':
		l.advance()
		return token.New(token.DIV, string(l.currentChar))

	case '%':
		l.advance()
		return token.New(token.MOD, string(l.currentChar))
	}

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
