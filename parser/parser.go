package parser

import (
	"fmt"
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
		fmt.Println(p.currentToken, tokenType)
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
	// program => compound_statement DOT

	node := p.compoundStatement()
	p.eat(token.DOT)
	return node
}

func (p *Parser) compoundStatement() interface{} {
	// compound_statement => BEGIN statement_list END

	p.eat(token.BEGIN)
	statements := p.statementList()
	p.eat(token.END)
	return Compound{
		Children: statements,
	}
}

func (p *Parser) statementList() []interface{} {
	/*
			statement_list => statement
		    	| statement SEMI statement_list
	*/
	node := p.statement()
	results := []interface{}{node}
	for p.currentToken.Type == token.SEMI {
		p.eat(token.SEMI)
		results = append(results, p.statement())
	}

	if p.currentToken.Type == token.ID {
		p.error()
	}

	return results
}

func (p *Parser) statement() interface{} {
	if p.currentToken.Type == token.BEGIN {
		return p.compoundStatement()
	} else if p.currentToken.Type == token.ID {
		return p.assignmentStatement()
	}
	return p.empty()
}

func (p *Parser) assignmentStatement() interface{} {
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
