package interpreter

import (
	"lsbasi/token"
	"strconv"
	"strings"
	"unicode"
)

type Interpreter struct {
	text         []rune
	pos          int
	currentToken *token.Token
}

func (i *Interpreter) error() {
	panic("Error parsing input")
}

func (i *Interpreter) readNumber() string {
	text := i.text
	var number strings.Builder
	for i.pos < len(text) && unicode.IsDigit(text[i.pos]) {
		number.WriteRune(text[i.pos])
		i.pos += 1
	}

	return number.String()
}

func (i *Interpreter) skipWhiteSpace() {
	text := i.text
	for unicode.IsSpace(text[i.pos]) {
		i.pos += 1
	}
}

func (i *Interpreter) getNextToken() *token.Token {
	text := i.text
	if i.pos > len(text)-1 {
		return token.New(token.EOF, "")
	}

	currentChar := text[i.pos]

	if unicode.IsSpace(currentChar) {
		i.skipWhiteSpace()
	}

	currentChar = text[i.pos]

	if unicode.IsDigit(currentChar) {
		number := i.readNumber()
		t := token.New(token.INTEGER, number)
		return t
	}

	if currentChar == '+' {
		t := token.New(token.PLUS, string(currentChar))
		i.pos += 1
		return t
	} else if currentChar == '-' {
		t := token.New(token.SUB, string(currentChar))
		i.pos += 1
		return t
	}

	i.error()

	return nil
}

func (i *Interpreter) eat(tokenType token.TokenType) {
	if i.currentToken.Type == tokenType {
		i.currentToken = i.getNextToken()
	} else {
		i.error()
	}
}

func (i *Interpreter) Expr() int64 {
	i.currentToken = i.getNextToken()
	left, _ := strconv.ParseInt(i.currentToken.Value, 10, 64)
	i.eat(token.INTEGER)

	op := i.currentToken
	i.eat(op.Type)

	right, _ := strconv.ParseInt(i.currentToken.Value, 10, 64)
	i.eat(token.INTEGER)

	var result int64
	if op.Type == token.PLUS {
		result = left + right
	} else if op.Type == token.SUB {
		result = left - right
	}
	return result
}

func New(text []rune) *Interpreter {
	return &Interpreter{
		text:         text,
		pos:          0,
		currentToken: nil,
	}
}
