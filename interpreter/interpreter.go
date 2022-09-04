package interpreter

import (
	"lsbasi/lexer"
	"lsbasi/token"
	"strconv"
)

type Interpreter struct {
	currentToken *token.Token
	lexer        *lexer.Lexer
}

func (i *Interpreter) error() {
	panic("Invalid syntax")
}

func (i *Interpreter) eat(tokenType token.TokenType) {
	if i.currentToken.Type == tokenType {
		i.currentToken = i.lexer.GetNextToken()
	} else {
		i.error()
	}
}

func (i *Interpreter) factor() int64 {
	var result int64
	if i.currentToken.Type == token.LPAREN {
		i.eat(token.LPAREN)
		result = i.expr()
		i.eat(token.RPAREN)
	} else if i.currentToken.Type == token.INTEGER {
		result, _ = strconv.ParseInt(i.currentToken.Value, 10, 64)
		i.eat(token.INTEGER)
	}
	return result
}

func (i *Interpreter) term() int64 {
	result := i.factor()

	tokens := []token.TokenType{
		token.MUL,
		token.DIV,
		token.MOD,
	}

	for i.in(i.currentToken.Type, tokens) {
		if i.currentToken.Type == token.MUL {
			i.eat(token.MUL)
			result *= i.factor()
		} else if i.currentToken.Type == token.DIV {
			i.eat(token.DIV)
			result /= i.factor()
		} else if i.currentToken.Type == token.MOD {
			i.eat(token.MOD)
			result %= i.factor()
		}
	}

	return result
}

func (i *Interpreter) expr() int64 {

	result := i.term()
	tokens := []token.TokenType{
		token.PLUS,
		token.SUB,
	}

	for i.in(i.currentToken.Type, tokens) {
		if i.currentToken.Type == token.PLUS {
			i.eat(token.PLUS)
			result += i.term()
		} else if i.currentToken.Type == token.SUB {
			i.eat(token.SUB)
			result -= i.term()
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
	i.currentToken = i.lexer.GetNextToken()
	return i.expr()
}

func New(text []rune) *Interpreter {
	return &Interpreter{
		currentToken: nil,
		lexer:        lexer.New(text),
	}
}
