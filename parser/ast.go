package parser

import (
	"strconv"

	"github.com/takeru56/tcompiler/token"
)

//
// Interface
//

// Node abstract Stmt and Expr
type Node interface {
	string() string
}

// Expr abstructs expression
type Expr interface {
	Node
	nodeExpr()
}

// Stmt abstructs some kinds of statements
type Stmt interface {
	Node
	nodeStmt()
}

// For Debugging
func (i InfixExpr) nodeExpr()           {}
func (i IntegerLiteral) nodeExpr()      {}
func (i IntegerRangeLiteral) nodeExpr() {}
func (b BoolLiteral) nodeExpr()         {}
func (i IdentExpr) nodeExpr()           {}
func (c CallExpr) nodeExpr()            {}
func (i InstantiationExpr) nodeExpr()   {}
func (c CallMethodExpr) nodeExpr()      {}
func (l LoopStmt) nodeStmt()            {}
func (a AssignStmt) nodeStmt()          {}
func (b BlockStmt) nodeStmt()           {}
func (i IfStmt) nodeStmt()              {}
func (w WhileStmt) nodeStmt()           {}
func (f FunctionDef) nodeStmt()         {}
func (r ReturnStmt) nodeStmt()          {}
func (c ClassDef) nodeStmt()            {}

//
// Expr
//

// OpKind express kind of operands as enum
type OpKind int

const (
	Add OpKind = iota
	Sub
	Mul
	Div
	EQ
	NEQ
	Less
	Greater
)

// InfixExpr has a operand and two nodes.
type InfixExpr struct {
	tok   token.Token
	Op    OpKind
	Left  Node
	Right Node
}

func (i InfixExpr) string() string {
	return "(" + i.Left.string() + " " + i.tok.Literal + " " + i.Right.string() + ")"
}

// IntegerLiteral express unsigned number
type IntegerLiteral struct {
	Tok token.Token
	Val int
}

func (i IntegerLiteral) string() string {
	return strconv.Itoa(i.Val)
}

type IntegerRangeLiteral struct {
	From IntegerLiteral
	To   IntegerLiteral
}

func (i IntegerRangeLiteral) string() string {
	return strconv.Itoa(i.From.Val) + ".." + strconv.Itoa(i.To.Val)
}

type BoolLiteral struct {
	Tok token.Token
}

func (b BoolLiteral) string() string {
	return b.Tok.Literal
}

// IdentKind show kind of the Identifier as enum
type IdentKind int

const (
	variable IdentKind = iota
	fn
)

// IdentExpr has kind and name
type IdentExpr struct {
	kind     IdentKind
	Name     string
	FSelf    bool
	ValType  IdentValType
	ValLimit IntegerRangeLiteral
}

func (i IdentExpr) string() string {
	if i.FSelf {
		if i.ValType == Exclude || i.ValType == Include {
			return "self." + i.Name + ": {" + valTypeDefinition[i.ValType] + ": " + i.ValLimit.string() + "}"
		}
		if i.ValType != Any {
			return "self." + i.Name + ": " + valTypeDefinition[i.ValType]
		}
		return "self." + i.Name
	}
	return i.Name
}

type IdentValType int

// 格納できる値の型を保持する
const (
	Num IdentValType = iota
	Bool
	Nil
	Range
	Any
	Include
	Exclude
)

var valTypeDefinition = map[IdentValType]string{
	Num:     "number",
	Bool:    "bool",
	Nil:     "nil",
	Include: "include",
	Exclude: "exclude",
}

func ValTypeToInt(vt IdentValType) int {
	switch vt {
	case Num:
		return 0
	case Bool:
		return 1
	case Nil:
		return 2
	case Include:
		return 5
	case Exclude:
		return 6
	}
	return 4
}

type CallExpr struct {
	Ident IdentExpr
	Args  []Node
}

func (c CallExpr) string() string {
	args := ""
	for _, arg := range c.Args {
		args += arg.string()
	}
	return c.Ident.Name + "(" + args + ")"
}

type InstantiationExpr struct {
	Ident IdentExpr
	Args  []Node
}

func (i InstantiationExpr) string() string {
	args := ""
	for _, arg := range i.Args {
		args += arg.string()
	}
	return i.Ident.Name + "(" + args + ")"
}

type CallMethodExpr struct {
	Receiver Node
	Method   Node
}

func (c CallMethodExpr) string() string {
	return c.Receiver.string() + "." + c.Method.string()
}

//
// Stmt
//

type AssignStmt struct {
	Ident IdentExpr
	Expr  Node
}

func (a AssignStmt) string() string {
	return a.Ident.string() + " = " + a.Expr.string()
}

type BlockStmt struct {
	Nodes []Node
}

func (b BlockStmt) string() string {
	s := "do\n"
	for _, stmt := range b.Nodes {
		s += "  " + stmt.string() + "\n"
	}
	return s + "end"
}

type IfStmt struct {
	Block     BlockStmt
	Condition Node
}

func (i IfStmt) string() string {
	s := "if " + i.Condition.string() + " then\n"
	for _, node := range i.Block.Nodes {
		s += "  " + node.string() + "\n"
	}
	return s + "end"
}

type ReturnStmt struct {
	Expr Node
}

func (r ReturnStmt) string() string {
	return "return " + r.Expr.string()
}

type WhileStmt struct {
	Block     BlockStmt
	Condition Node
}

func (w WhileStmt) string() string {
	s := "while " + w.Condition.string() + " do\n"
	for _, node := range w.Block.Nodes {
		s += "  " + node.string() + "\n"
	}
	return s + "end"
}

// LoopStmt has a block
type LoopStmt struct {
	block []Stmt
}

func (l LoopStmt) string() string {
	s := "loop {"
	for _, b := range l.block {
		s += " " + b.string()
	}
	return s + " }"
}

type FunctionDef struct {
	Ident      IdentExpr
	Block      BlockStmt
	Args       []IdentExpr
	FlagMethod bool
}

func (f FunctionDef) string() string {
	s := "def " + f.Ident.Name + "("
	for i, arg := range f.Args {
		if i > 0 {
			s += ", "
		}
		s += arg.Name
	}
	s += ")\n"
	for _, b := range f.Block.Nodes {
		s += "  " + b.string() + "\n"
	}
	return s + "end\n"
}

type ClassDef struct {
	Ident   IdentExpr
	Methods []FunctionDef
}

func (c ClassDef) string() string {
	s := "class " + c.Ident.Name + "\n"
	for _, m := range c.Methods {
		s += m.string()
	}
	s += "end"
	return s
}
