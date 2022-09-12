package token

import "fmt"

type TokenType string

const (
	INTEGER       TokenType = "INTEGER"
	PLUS          TokenType = "PLUS"
	SUB           TokenType = "SUB"
	MUL           TokenType = "MUL"
	DIV           TokenType = "DIV"
	MOD           TokenType = "MOD"
	LPAREN        TokenType = "LPAREN"
	RPAREN        TokenType = "RPAREN"
	LBRACE        TokenType = "LBRACE"
	RBRACE        TokenType = "RBRACE"
	SEMI          TokenType = "SEMI"
	ASSIGN        TokenType = "ASSGIN"
	VAR           TokenType = "VAR"
	ID            TokenType = "ID"
	BEGIN         TokenType = "BEGIN"
	END           TokenType = "END"
	DOT           TokenType = "DOT"
	PROGRAM       TokenType = "PROGRAM"
	INTEGER_DIV   TokenType = "INTERGER_DIV"
	FLOAT_DIV     TokenType = "FLOAT_DIV"
	REAL          TokenType = "REAL"
	INTEGER_CONST TokenType = "INTEGER_CONST"
	FLOAT_CONST   TokenType = "FLOAT_CONST"
	COLON         TokenType = "COLON"
	COMMA         TokenType = "COMMA"
	EOF           TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v, %v)", t.Type, t.Value)
}

func New(tokenType TokenType, value string) *Token {
	t := &Token{
		Type:  tokenType,
		Value: value,
	}

	return t
}
