package ast

import (
	"testing"

	"github.com/tom96da/sleepingknights/pkg/token"
)

func TestProgramPosBranches(t *testing.T) {
	p1 := &Program{Functions: []*FunctionDecl{{Line: 2, Column: 3}}}
	l, c := p1.Pos()
	if l != 2 || c != 3 {
		t.Fatalf("unexpected function-first pos: %d:%d", l, c)
	}

	p2 := &Program{Statements: []Statement{&ExpressionStatement{Line: 4, Column: 5}}}
	l, c = p2.Pos()
	if l != 4 || c != 5 {
		t.Fatalf("unexpected statement-first pos: %d:%d", l, c)
	}

	p3 := &Program{}
	l, c = p3.Pos()
	if l != 1 || c != 1 {
		t.Fatalf("unexpected default pos: %d:%d", l, c)
	}
}

func TestAllNodePosAndMarkerMethods(t *testing.T) {
	ident := &Identifier{Name: "x", Line: 1, Column: 2}
	intLit := &IntegerLiteral{Value: "1", Line: 2, Column: 3}
	floatLit := &FloatLiteral{Value: "1.2", Line: 3, Column: 4}
	strLit := &StringLiteral{Value: "s", Line: 4, Column: 5}
	boolLit := &BooleanLiteral{Value: true, Line: 5, Column: 6}
	unary := &UnaryExpression{Operator: token.NOT, Operand: boolLit, Line: 6, Column: 7}
	binary := &BinaryExpression{Left: intLit, Operator: token.PLUS, Right: floatLit, Line: 7, Column: 8}
	assign := &AssignmentExpression{Name: ident, Value: intLit, Line: 8, Column: 9}
	call := &FunctionCall{Function: ident, Arguments: []Expression{strLit}, Line: 9, Column: 10}

	block := &BlockStatement{Statements: []Statement{}, Line: 10, Column: 11}
	varDecl := &VariableDecl{IsConst: false, Name: ident, Value: intLit, Line: 11, Column: 12}
	ifStmt := &IfStatement{Condition: boolLit, Then: block, Line: 12, Column: 13}
	whileStmt := &WhileStatement{Condition: boolLit, Body: block, Line: 13, Column: 14}
	forStmt := &ForStatement{Init: varDecl, Condition: boolLit, Update: assign, Body: block, Line: 14, Column: 15}
	retStmt := &ReturnStatement{Value: intLit, Line: 15, Column: 16}
	breakStmt := &BreakStatement{Line: 16, Column: 17}
	continueStmt := &ContinueStatement{Line: 17, Column: 18}
	exprStmt := &ExpressionStatement{Expr: call, Line: 18, Column: 19}

	typeName := &TypeName{Name: "int", Line: 19, Column: 20}
	param := &Parameter{Name: "a", Type: typeName, Line: 20, Column: 21}
	fn := &FunctionDecl{Name: ident, Parameters: []*Parameter{param}, ReturnType: typeName, Body: block, Line: 21, Column: 22}

	assertPos(t, ident, 1, 2)
	assertPos(t, intLit, 2, 3)
	assertPos(t, floatLit, 3, 4)
	assertPos(t, strLit, 4, 5)
	assertPos(t, boolLit, 5, 6)
	assertPos(t, unary, 6, 7)
	assertPos(t, binary, 7, 8)
	assertPos(t, assign, 8, 9)
	assertPos(t, call, 9, 10)
	assertPos(t, block, 10, 11)
	assertPos(t, varDecl, 11, 12)
	assertPos(t, ifStmt, 12, 13)
	assertPos(t, whileStmt, 13, 14)
	assertPos(t, forStmt, 14, 15)
	assertPos(t, retStmt, 15, 16)
	assertPos(t, breakStmt, 16, 17)
	assertPos(t, continueStmt, 17, 18)
	assertPos(t, exprStmt, 18, 19)
	assertPos(t, typeName, 19, 20)
	assertPos(t, param, 20, 21)
	assertPos(t, fn, 21, 22)

	_ = []Statement{block, varDecl, ifStmt, whileStmt, forStmt, retStmt, breakStmt, continueStmt, exprStmt}
	_ = []Expression{ident, intLit, floatLit, strLit, boolLit, unary, binary, assign, call}

	// Call marker methods explicitly so coverage includes these empty method bodies.
	block.statementNode()
	varDecl.statementNode()
	ifStmt.statementNode()
	whileStmt.statementNode()
	forStmt.statementNode()
	retStmt.statementNode()
	breakStmt.statementNode()
	continueStmt.statementNode()
	exprStmt.statementNode()

	assign.expressionNode()
	binary.expressionNode()
	unary.expressionNode()
	call.expressionNode()
	ident.expressionNode()
	intLit.expressionNode()
	floatLit.expressionNode()
	strLit.expressionNode()
	boolLit.expressionNode()
}

func assertPos(t *testing.T, n Node, wantL int, wantC int) {
	t.Helper()
	l, c := n.Pos()
	if l != wantL || c != wantC {
		t.Fatalf("unexpected pos: got=%d:%d want=%d:%d", l, c, wantL, wantC)
	}
}
