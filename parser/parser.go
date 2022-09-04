package parser

import (
	"lsbasi/lexer"
	"lsbasi/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken *token.Token
}

func (p *Parser) error() {
	panic("Invalid syntax")
}

func (p *Parser) eat(tokenType token.TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.GetNextToken()
	} else {
		p.error()
	}
}

func (p *Parser) factor() interface{} {
	var l interface{}
	if p.currentToken.Type == token.LPAREN {
		p.eat(token.LPAREN)
		l = p.expr()
		p.eat(token.RPAREN)
	} else if p.currentToken.Type == token.INTEGER {
		l = Num{
			Value: p.currentToken,
		}
		p.eat(token.INTEGER)
	}
	return l
}

func (p *Parser) term() interface{} {
	l := p.factor()

	tokens := []token.TokenType{
		token.MUL,
		token.DIV,
		token.MOD,
	}

	for p.in(p.currentToken.Type, tokens) {
		currentToken := p.currentToken
		if p.currentToken.Type == token.MUL {
			p.eat(token.MUL)
		} else if p.currentToken.Type == token.DIV {
			p.eat(token.DIV)
		} else if p.currentToken.Type == token.MOD {
			p.eat(token.MOD)
		}
		l = BinOp{
			Left:  l,
			Op:    currentToken,
			Right: p.factor(),
		}
	}

	return l
}

func (p *Parser) expr() interface{} {

	l := p.term()
	tokens := []token.TokenType{
		token.PLUS,
		token.SUB,
	}

	for p.in(p.currentToken.Type, tokens) {
		currentToken := p.currentToken
		if p.currentToken.Type == token.PLUS {
			p.eat(token.PLUS)
		} else if p.currentToken.Type == token.SUB {
			p.eat(token.SUB)
		}
		l = BinOp{
			Left:  l,
			Op:    currentToken,
			Right: p.term(),
		}
	}
	return l
}

func (p *Parser) in(tokenType token.TokenType, tokenTypes []token.TokenType) bool {
	for _, tt := range tokenTypes {
		if tokenType == tt {
			return true
		}
	}

	return false
}

func (p *Parser) Parse() interface{} {
	p.currentToken = p.lexer.GetNextToken()
	return p.expr()
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer:        l,
		currentToken: nil,
	}
}
