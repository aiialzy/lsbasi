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
	} else if p.currentToken.Type == token.PLUS || p.currentToken.Type == token.SUB {
		op := p.currentToken
		p.eat(op.Type)
		l = UnaryOp{
			Op:    op,
			Right: p.factor(),
		}
	} else if p.currentToken.Type == token.INTEGER {
		l = Num{
			Value: p.currentToken,
		}
		p.eat(token.INTEGER)
	} else if p.currentToken.Type == token.ID {
		l = ID{
			Token: p.currentToken,
		}
		p.eat(token.ID)
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

func (p *Parser) program() interface{} {
	return p.compoundStatement()
}

func (p *Parser) compoundStatement() interface{} {
	p.eat(token.LBRACE)
	statements := []interface{}{}
	for p.currentToken.Type != token.RBRACE {
		statements = append(statements, p.statement())
	}
	p.eat(token.RBRACE)
	return Compound{
		Children: statements,
	}
}

func (p *Parser) statement() interface{} {
	if p.currentToken.Type == token.LBRACE {
		return p.compoundStatement()
	} else if p.currentToken.Type == token.VAR {
		statement := p.assignmentStatement()
		p.eat(token.SEMI)
		return statement
	} else if p.currentToken.Type == token.SEMI {
		p.eat(token.SEMI)
	}
	return p.empty()
}

func (p *Parser) assignmentStatement() interface{} {
	p.eat(token.VAR)

	id := p.id()

	op := p.currentToken
	p.eat(token.ASSIGN)

	value := p.expr()
	return Assign{
		Left:  id,
		Op:    op,
		Right: value,
	}
}

func (p *Parser) id() interface{} {
	id := ID{
		Token: p.currentToken,
	}
	p.eat(token.ID)
	return id
}

func (p *Parser) empty() interface{} {
	return NoOp{}
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
	return p.program()
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer:        l,
		currentToken: nil,
	}
}
