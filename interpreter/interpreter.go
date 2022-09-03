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
	currentChar  rune
}

func (i *Interpreter) error() {
	panic("Invalid syntax")
}

func (i *Interpreter) integer() string {
	var number strings.Builder
	for unicode.IsDigit(i.currentChar) {
		number.WriteRune(i.currentChar)
		i.advance()
	}

	return number.String()
}

func (i *Interpreter) skipWhiteSpace() {
	for unicode.IsSpace(i.currentChar) {
		i.advance()
	}
}

func (i *Interpreter) advance() {
	i.pos += 1
	if i.pos < len(i.text) {
		i.currentChar = i.text[i.pos]
	} else {
		i.currentChar = 0
	}
}

func (i *Interpreter) getNextToken() *token.Token {

	if unicode.IsSpace(i.currentChar) {
		i.skipWhiteSpace()
	}

	if i.currentChar == 0 {
		return token.New(token.EOF, "")
	}

	if unicode.IsDigit(i.currentChar) {
		number := i.integer()
		return token.New(token.INTEGER, number)
	}

	switch i.currentChar {
	case '+':
		i.advance()
		return token.New(token.PLUS, string(i.currentChar))

	case '-':
		i.advance()
		return token.New(token.SUB, string(i.currentChar))

	case '*':
		i.advance()
		return token.New(token.MUL, string(i.currentChar))

	case '/':
		i.advance()
		return token.New(token.DIV, string(i.currentChar))

	case '%':
		i.advance()
		return token.New(token.MOD, string(i.currentChar))
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

func (i *Interpreter) term() int64 {
	result, _ := strconv.ParseInt(i.currentToken.Value, 10, 64)
	i.eat(token.INTEGER)
	return result
}

func (i *Interpreter) expr() int64 {

	result := i.term()
	tokens := []token.TokenType{
		token.PLUS,
		token.SUB,
	}

	for i.in(i.currentToken.Type, tokens) {
		switch i.currentToken.Type {
		case token.PLUS:
			i.eat(token.PLUS)
			result += i.term()

		case token.SUB:
			i.eat(token.SUB)
			result -= i.term()

		case token.MUL:
			i.eat(token.MUL)
			result *= i.term()

		case token.DIV:
			i.eat(token.DIV)
			result /= i.term()

		case token.MOD:
			i.eat(token.MOD)
			result %= i.term()
		}
	}
	return result
}

func (i *Interpreter) in(tokenType token.TokenType, tokenTypes []token.TokenType) bool {
	for _, tt := range tokenTypes {
		if tokenType == tt {
			return true
		}
	}

	return false
}

func (i *Interpreter) Interprete() int64 {
	i.currentToken = i.getNextToken()
	return i.expr()
}

func New(text []rune) *Interpreter {
	if len(text) == 0 {
		panic("text length must bigger than 0")
	}

	return &Interpreter{
		text:         text,
		pos:          0,
		currentToken: nil,
		currentChar:  text[0],
	}
}
