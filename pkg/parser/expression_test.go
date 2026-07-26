package parser

import (
	"testing"

	"github.com/tom96da/sleepingknights/pkg/ast"
	"github.com/tom96da/sleepingknights/pkg/token"
)

func TestParseExpressionPrimitivesAndUnary(t *testing.T) {
	src := "(1)\n-2\nnot false\n1.25\n\"str\"\n"
	p := New(src)
	_, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}
}

func TestParseAssignmentRightAssociative(t *testing.T) {
	src := "let a = 0\nlet b = 0\na = b = 3\n"
	p := New(src)
	program, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}

	exprStmt, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement[2] type mismatch: %T", program.Statements[2])
	}
	outer, ok := exprStmt.Expr.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("expected outer assignment, got %T", exprStmt.Expr)
	}
	if _, ok := outer.Value.(*ast.AssignmentExpression); !ok {
		t.Fatalf("expected nested assignment on RHS, got %T", outer.Value)
	}
}

func TestParseExpressionErrorBranches(t *testing.T) {
	cases := []string{
		"(1 + 2\n",      // grouped expression missing ')'
		"(\n",           // grouped expression inner parse nil
		"1(2)\n",        // call target not identifier
		"foo(1\n",       // call missing ')'
		"foo(,)\n",      // call argument parse nil
		"(x + 1) = 2\n", // assignment left must be identifier
		"x = )\n",       // assignment rhs parse nil
		"1 + )\n",       // binary rhs parse nil
		"not )\n",       // unary rhs parse nil
		"(let x = 1)\n", // unexpected token in expression
	}

	for i, src := range cases {
		t.Run("case", func(t *testing.T) {
			p := New(src)
			_, errs := p.ParseProgram()
			if len(errs) == 0 {
				t.Fatalf("case %d expected parse errors", i)
			}
		})
	}
}

func TestSkipCallArgumentsBranches(t *testing.T) {
	p := &Parser{tokens: []token.Token{
		{Type: token.LPAREN},
		{Type: token.INT, Literal: "1"},
		{Type: token.RPAREN},
		{Type: token.EOF},
	}}
	p.skipCallArguments()
	if p.current().Type != token.EOF {
		t.Fatalf("expected EOF after skipCallArguments, got=%s", p.current().Type)
	}

	p2 := &Parser{tokens: []token.Token{
		{Type: token.LPAREN},
		{Type: token.LPAREN},
		{Type: token.INT, Literal: "1"},
		{Type: token.EOF},
	}}
	p2.skipCallArguments()
	if p2.current().Type != token.EOF {
		t.Fatalf("expected EOF on unterminated args, got=%s", p2.current().Type)
	}
}

func TestParseInfixDefaultBranch(t *testing.T) {
	p := &Parser{}
	left := &ast.Identifier{Name: "x", Line: 1, Column: 1}
	op := token.Token{Type: token.NEWLINE, Line: 1, Column: 2}
	got := p.parseInfixExpression(left, op)
	if got != nil {
		t.Fatal("expected nil for unexpected infix operator")
	}
	if len(p.errors) == 0 {
		t.Fatal("expected error for unexpected infix operator")
	}
}
