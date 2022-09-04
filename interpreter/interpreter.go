package interpreter

import (
	"fmt"
	"lsbasi/lexer"
	"lsbasi/parser"
	"lsbasi/token"
	"strconv"
)

type Interpreter struct {
	parser *parser.Parser
}

func (i *Interpreter) error() {
	panic("Invalid syntax")
}

func (i *Interpreter) visit(node interface{}) int64 {
	switch n := node.(type) {
	case parser.Num:
		return i.visitNum(n)

	case parser.BinOp:
		return i.visitBinOp(n)
	}

	i.error()
	return 0
}

func (i *Interpreter) visitNum(node parser.Num) int64 {
	value := node.Value.(*token.Token)
	result, _ := strconv.ParseInt(value.Value.(string), 10, 64)
	fmt.Printf("%v ", result)
	return result
}

func (i *Interpreter) visitBinOp(node parser.BinOp) int64 {
	op := node.Op.(*token.Token)
	fmt.Printf("(%v ", op.Value.(string))

	l := i.visit(node.Left)
	r := i.visit(node.Right)
	fmt.Printf(")")

	// fmt.Printf("(%v ", op.Value.(string))

	switch op.Type {
	case token.PLUS:
		return l + r

	case token.SUB:
		return l - r

	case token.MUL:
		return l * r

	case token.DIV:
		return l / r

	case token.MOD:
		return l % r
	}

	i.error()
	return 0
}

func (i *Interpreter) Interprete() int64 {
	root := i.parser.Parse()
	result := i.visit(root)
	return result
}

func New(text []rune) *Interpreter {
	l := lexer.New(text)
	return &Interpreter{
		parser: parser.New(l),
	}
}
