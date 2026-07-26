package ast

import "github.com/tom96da/sleepingknights/pkg/token"

// Node is the base interface for all AST nodes.
type Node interface {
	Pos() (line int, column int)
}

// Statement is the base interface for statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression is the base interface for expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root AST node.
type Program struct {
	Functions  []*FunctionDecl
	Statements []Statement
}

func (p *Program) Pos() (int, int) {
	if len(p.Functions) > 0 {
		return p.Functions[0].Pos()
	}
	if len(p.Statements) > 0 {
		return p.Statements[0].Pos()
	}
	return 1, 1
}

// TypeName represents a declared type such as int/float/bool/string/void.
type TypeName struct {
	Name   string
	Line   int
	Column int
}

func (t *TypeName) Pos() (int, int) { return t.Line, t.Column }

// Parameter is a function parameter declaration.
type Parameter struct {
	Name   string
	Type   *TypeName
	Line   int
	Column int
}

func (p *Parameter) Pos() (int, int) { return p.Line, p.Column }

// FunctionDecl is a top-level function declaration.
type FunctionDecl struct {
	DocComment string
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *TypeName
	Body       *BlockStatement
	Line       int
	Column     int
}

func (f *FunctionDecl) Pos() (int, int) { return f.Line, f.Column }

// BlockStatement is a braced sequence of statements.
type BlockStatement struct {
	Statements []Statement
	Line       int
	Column     int
}

func (b *BlockStatement) Pos() (int, int) { return b.Line, b.Column }
func (b *BlockStatement) statementNode()  {}

// VariableDecl is let/const declaration.
type VariableDecl struct {
	IsConst bool
	Name    *Identifier
	Value   Expression
	Line    int
	Column  int
}

func (v *VariableDecl) Pos() (int, int) { return v.Line, v.Column }
func (v *VariableDecl) statementNode()  {}

// IfStatement models if/else and if-let.
type IfStatement struct {
	Condition  Expression
	LetBinding *VariableDecl
	Then       *BlockStatement
	Else       Statement
	Line       int
	Column     int
}

func (i *IfStatement) Pos() (int, int) { return i.Line, i.Column }
func (i *IfStatement) statementNode()  {}

// WhileStatement models a while loop.
type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
	Line      int
	Column    int
}

func (w *WhileStatement) Pos() (int, int) { return w.Line, w.Column }
func (w *WhileStatement) statementNode()  {}

// ForStatement models a for loop.
type ForStatement struct {
	Init      *VariableDecl
	Condition Expression
	Update    Expression
	Body      *BlockStatement
	Line      int
	Column    int
}

func (f *ForStatement) Pos() (int, int) { return f.Line, f.Column }
func (f *ForStatement) statementNode()  {}

// ReturnStatement returns from function.
type ReturnStatement struct {
	Value  Expression
	Line   int
	Column int
}

func (r *ReturnStatement) Pos() (int, int) { return r.Line, r.Column }
func (r *ReturnStatement) statementNode()  {}

// BreakStatement exits loop.
type BreakStatement struct {
	Line   int
	Column int
}

func (b *BreakStatement) Pos() (int, int) { return b.Line, b.Column }
func (b *BreakStatement) statementNode()  {}

// ContinueStatement jumps to next loop iteration.
type ContinueStatement struct {
	Line   int
	Column int
}

func (c *ContinueStatement) Pos() (int, int) { return c.Line, c.Column }
func (c *ContinueStatement) statementNode()  {}

// ExpressionStatement wraps an expression used as a statement.
type ExpressionStatement struct {
	Expr   Expression
	Line   int
	Column int
}

func (e *ExpressionStatement) Pos() (int, int) { return e.Line, e.Column }
func (e *ExpressionStatement) statementNode()  {}

// AssignmentExpression assigns a value to an identifier.
type AssignmentExpression struct {
	Name   *Identifier
	Value  Expression
	Line   int
	Column int
}

func (a *AssignmentExpression) Pos() (int, int) { return a.Line, a.Column }
func (a *AssignmentExpression) expressionNode() {}

// BinaryExpression models a binary operator.
type BinaryExpression struct {
	Left     Expression
	Operator token.TokenType
	Right    Expression
	Line     int
	Column   int
}

func (b *BinaryExpression) Pos() (int, int) { return b.Line, b.Column }
func (b *BinaryExpression) expressionNode() {}

// UnaryExpression models a unary operator.
type UnaryExpression struct {
	Operator token.TokenType
	Operand  Expression
	Line     int
	Column   int
}

func (u *UnaryExpression) Pos() (int, int) { return u.Line, u.Column }
func (u *UnaryExpression) expressionNode() {}

// FunctionCall models a function call.
type FunctionCall struct {
	Function  *Identifier
	Arguments []Expression
	Line      int
	Column    int
}

func (f *FunctionCall) Pos() (int, int) { return f.Line, f.Column }
func (f *FunctionCall) expressionNode() {}

// Identifier is an identifier expression.
type Identifier struct {
	Name   string
	Line   int
	Column int
}

func (i *Identifier) Pos() (int, int) { return i.Line, i.Column }
func (i *Identifier) expressionNode() {}

// IntegerLiteral is an integer literal.
type IntegerLiteral struct {
	Value  string
	Line   int
	Column int
}

func (i *IntegerLiteral) Pos() (int, int) { return i.Line, i.Column }
func (i *IntegerLiteral) expressionNode() {}

// FloatLiteral is a float literal.
type FloatLiteral struct {
	Value  string
	Line   int
	Column int
}

func (f *FloatLiteral) Pos() (int, int) { return f.Line, f.Column }
func (f *FloatLiteral) expressionNode() {}

// StringLiteral is a string literal.
type StringLiteral struct {
	Value  string
	Line   int
	Column int
}

func (s *StringLiteral) Pos() (int, int) { return s.Line, s.Column }
func (s *StringLiteral) expressionNode() {}

// BooleanLiteral is a bool literal.
type BooleanLiteral struct {
	Value  bool
	Line   int
	Column int
}

func (b *BooleanLiteral) Pos() (int, int) { return b.Line, b.Column }
func (b *BooleanLiteral) expressionNode() {}
