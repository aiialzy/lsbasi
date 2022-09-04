package lexer

import (
	"fmt"
	"lsbasi/token"
	"strings"
	"unicode"
)

var reserved_keywords = map[string]*token.Token{
	"var": token.New(token.VAR, "var"),
}

type Lexer struct {
	text        []rune
	pos         int
	currentChar rune
}

func (l *Lexer) integer() string {
	var number strings.Builder
	for unicode.IsDigit(l.currentChar) ||
		l.currentChar == 'x' ||
		l.currentChar == 'X' ||
		l.currentChar == 'o' ||
		l.currentChar == 'O' ||
		l.currentChar == 'b' ||
		l.currentChar == 'B' {
		number.WriteRune(l.currentChar)
		l.advance()
	}

	return number.String()
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

func (l *Lexer) skipComment(isMultiLine bool) {
	var comment strings.Builder
	if isMultiLine {
		comment.WriteRune(l.currentChar)
		l.advance()
		comment.WriteRune(l.currentChar)
		l.advance()
		for {
			if l.currentChar == '*' && l.peek() == '/' {
				break
			}
			comment.WriteRune(l.currentChar)
			l.advance()
		}
		comment.WriteRune(l.currentChar)
		l.advance()
		comment.WriteRune(l.currentChar)
		l.advance()
	} else {
		for l.currentChar != '\n' {
			comment.WriteRune(l.currentChar)
			l.advance()
		}
		comment.WriteRune(l.currentChar)
		l.advance()
	}

	// fmt.Println("comment: ", comment.String())
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

	for unicode.IsSpace(l.currentChar) || l.currentChar == '/' {
		if unicode.IsSpace(l.currentChar) {
			l.skipWhiteSpace()
		}

		if l.currentChar == '/' {
			nextChar := l.peek()
			if nextChar == '/' {
				l.skipComment(false)
			} else if nextChar == '*' {
				l.skipComment(true)
			} else {
				break
			}
		}

	}

	if l.currentChar == 0 {
		return token.New(token.EOF, "")
	}

	if unicode.IsDigit(l.currentChar) {
		number := l.integer()
		return token.New(token.INTEGER, number)
	}

	if l.currentChar == '_' || unicode.IsLetter(l.currentChar) {
		word := l.word()
		if t, exists := reserved_keywords[word]; exists {
			return t
		}
		return token.New(token.ID, word)
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
		return token.New(token.DIV, "/")

	case '%':
		l.advance()
		return token.New(token.MOD, "%")

	case '(':
		l.advance()
		return token.New(token.LPAREN, "(")

	case ')':
		l.advance()
		return token.New(token.RPAREN, ")")

	case '{':
		l.advance()
		return token.New(token.LBRACE, "{")

	case '}':
		l.advance()
		return token.New(token.RBRACE, "}")

	case ';':
		l.advance()
		return token.New(token.SEMI, ";")

	case '=':
		l.advance()
		return token.New(token.ASSIGN, "=")
	}

	fmt.Printf("%v\n", l.text[l.pos:])
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
