package parser

type AST struct {
}

type Program struct {
	AST
	Name  string
	Block interface{}
}

type Block struct {
	AST
	Declarations      []interface{}
	CompoundStatement interface{}
}

type VarDecl struct {
	AST
	VarNode  interface{}
	TypeNode interface{}
}

type Type struct {
	AST
	Token interface{}
}

type Compound struct {
	AST
	Children []interface{}
}

type Assign struct {
	AST
	Left  interface{}
	Op    interface{}
	Right interface{}
}

type ID struct {
	AST
	Token interface{}
}

type BinOp struct {
	AST
	Left  interface{}
	Op    interface{}
	Right interface{}
}

type UnaryOp struct {
	AST
	Op    interface{}
	Right interface{}
}

type Num struct {
	AST
	Value interface{}
}

type NoOp struct {
	AST
}
