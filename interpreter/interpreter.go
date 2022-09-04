package interpreter

import (
	"fmt"
	"lsbasi/lexer"
	"lsbasi/parser"
	"lsbasi/token"
	"strconv"
)

var globalScope = map[string]interface{}{}

type Interpreter struct {
	parser *parser.Parser
}

func (i *Interpreter) error() {
	panic("Invalid syntax")
}

func (i *Interpreter) visit(node interface{}) interface{} {
	switch n := node.(type) {
	case parser.Num:
		return i.visitNum(n)

	case parser.UnaryOp:
		return i.visitUnaryOp(n)

	case parser.BinOp:
		return i.visitBinOp(n)

	case parser.Compound:
		return i.visitCompound(n)

	case parser.Assign:
		return i.visitAssign(n)

	case parser.ID:
		return i.visitID(n)

	case parser.NoOp:
		return nil
	}

	i.error()
	return 0
}

func (i *Interpreter) visitNum(node parser.Num) interface{} {
	value := node.Value.(*token.Token)
	numStr := value.Value
	base := 10
	if len(numStr) > 1 && numStr[0] == '0' {
		if numStr[1] == 'x' || numStr[1] == 'X' {
			base = 16
			numStr = numStr[2:]
		} else if numStr[1] == 'b' || numStr[1] == 'B' {
			base = 2
			numStr = numStr[2:]
		} else if numStr[1] == 'o' || numStr[1] == 'O' {
			base = 8
			numStr = numStr[2:]
		} else {
			base = 8
			numStr = numStr[1:]
		}
	}
	result, err := strconv.ParseInt(numStr, base, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid integer: %v", value.Value))
	}
	return result
}

func (i *Interpreter) visitUnaryOp(node parser.UnaryOp) interface{} {
	op := node.Op.(*token.Token)

	if op.Type == token.PLUS {
		return i.visit(node.Right)
	} else if op.Type == token.SUB {
		return -i.visit(node.Right).(int64)
	}

	i.error()
	return 0
}

func (i *Interpreter) visitBinOp(node parser.BinOp) interface{} {
	op := node.Op.(*token.Token)

	l := i.visit(node.Left).(int64)
	r := i.visit(node.Right).(int64)

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

func (i *Interpreter) visitCompound(node parser.Compound) interface{} {
	for _, child := range node.Children {
		i.visit(child)
	}

	return nil
}

func (i *Interpreter) visitAssign(node parser.Assign) interface{} {
	l := node.Left.(parser.ID).Token.(*token.Token).Value
	r := i.visit(node.Right)
	globalScope[l] = r

	return nil
}

func (i *Interpreter) visitID(node parser.ID) interface{} {
	t := node.Token.(*token.Token)
	if v, exists := globalScope[t.Value]; !exists {
		fmt.Printf("undefined variable %v\n", t.Value)
		i.error()
	} else {
		return v
	}

	return nil
}

func (i *Interpreter) Interprete() {
	root := i.parser.Parse()
	i.visit(root)
}

func (i *Interpreter) PrintValues() {
	fmt.Println(globalScope)
}

func New(text []rune) *Interpreter {
	l := lexer.New(text)
	return &Interpreter{
		parser: parser.New(l),
	}
}
