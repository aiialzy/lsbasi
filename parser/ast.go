package parser

type AST struct {
}

type BinOp struct {
	AST
	Left  interface{}
	Op    interface{}
	Right interface{}
}

type Num struct {
	AST
	Value interface{}
}
