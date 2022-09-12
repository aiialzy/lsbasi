package interpreter

import (
	"fmt"
	"lsbasi/lexer"
	"lsbasi/parser"
	"lsbasi/token"
	"strconv"
	"strings"
)

var globalScope = map[string]interface{}{}

type Interpreter struct {
	parser *parser.Parser
}

func (i *Interpreter) error() {
	panic("Invalid syntax")
}

func (i *Interpreter) visit(node interface{}) interface{} {
	// fmt.Printf("%T\n", node)
	// fmt.Println(node)
	switch n := node.(type) {
	case parser.Num:
		return i.visitNum(n)

	case parser.Program:
		return i.visitProgram(n)

	case parser.Block:
		return i.visitBlock(n)

	case parser.Type:
		return i.visitType(n)

	case parser.VarDecl:
		return i.visitVarDecl(n)

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
	var result interface{}
	var err error
	if strings.Contains(value.Value, ".") {
		result, err = strconv.ParseFloat(numStr, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid float: %v", value.Value))
		}
	} else {
		result, err = strconv.ParseInt(numStr, base, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid integer: %v", value.Value))
		}
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

func (i *Interpreter) visitProgram(node parser.Program) interface{} {
	i.visit(node.Block)
	return nil
}

func (i *Interpreter) visitBlock(node parser.Block) interface{} {
	for _, declaration := range node.Declarations {
		i.visit(declaration)
	}
	i.visit(node.CompoundStatement)
	return nil
}

func (i *Interpreter) visitVarDecl(node parser.VarDecl) interface{} {
	return nil
}

func (i *Interpreter) visitType(node parser.Type) interface{} {
	return nil
}

func (i *Interpreter) visitBinOp(node parser.BinOp) interface{} {
	op := node.Op.(*token.Token)

	l := i.visit(node.Left)
	r := i.visit(node.Right)

	if li, ok := l.(int64); ok && op.Type != token.FLOAT_DIV {
		ri := r.(int64)
		switch op.Type {
		case token.PLUS:
			return li + ri

		case token.SUB:
			return li - ri

		case token.MUL:
			return li * ri

		case token.MOD:
			return li % ri

		case token.INTEGER_DIV:
			return li / ri
		}
	} else {
		var lf, rf float64
		if _, ok := l.(float64); !ok {
			lf = float64(l.(int64))
			rf = float64(r.(int64))
		} else {
			lf = l.(float64)
			rf = r.(float64)
		}

		switch op.Type {
		case token.PLUS:
			return lf + rf

		case token.SUB:
			return lf - rf

		case token.MUL:
			return lf * rf

		case token.FLOAT_DIV:
			return lf / rf
		}
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
