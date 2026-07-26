package parser

import (
	"testing"

	"github.com/tom96da/sleepingknights/pkg/ast"
)

func TestParseVariableAndAssignmentStatement(t *testing.T) {
	src := "let x = 1\nx = x + 2\n"
	p := New(src)
	program, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("unexpected statement count: got=%d", len(program.Statements))
	}

	decl, ok := program.Statements[0].(*ast.VariableDecl)
	if !ok {
		t.Fatalf("statement[0] type mismatch: %T", program.Statements[0])
	}
	if decl.Name.Name != "x" {
		t.Fatalf("unexpected decl name: %s", decl.Name.Name)
	}

	exprStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement[1] type mismatch: %T", program.Statements[1])
	}
	if _, ok := exprStmt.Expr.(*ast.AssignmentExpression); !ok {
		t.Fatalf("expected assignment expression, got %T", exprStmt.Expr)
	}
}

func TestParseControlFlowStatements(t *testing.T) {
	src := "const c = 1\n" +
		"while (true) {\n" +
		"break\n" +
		"continue\n" +
		"return\n" +
		"return 1\n" +
		"}\n" +
		"if (true and false or not false) {\n" +
		"foo(1, 2)\n" +
		"} else {\n" +
		"return\n" +
		"}\n" +
		"for (let i = 0; i < 3; i = i + 1) {\n" +
		"print(i)\n" +
		"}\n" +
		"{\n" +
		"let z = 1\n" +
		"}\n"
	p := New(src)
	_, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}
}

func TestParseElseIfChainStatement(t *testing.T) {
	src := "if (true) {\n" +
		"return\n" +
		"} else if (false) {\n" +
		"return\n" +
		"}\n"
	p := New(src)
	_, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}
}

func TestParseStatementErrorBranches(t *testing.T) {
	cases := []string{
		"let = 1\n",                                 // decl: missing identifier
		"let x 1\n",                                 // decl: missing '='
		"let x = 1 let y = 2\n",                     // missing statement newline
		"if true) {\n}\n",                           // if: missing '('
		"if (true {\n}\n",                           // if: missing ')'
		"if (true) true\n",                          // if: missing block
		"if (let = 1) {\n}\n",                       // if-let decl parse nil
		"if ()) {\n}\n",                             // if condition parse nil
		"if (let x = 1\n",                           // if condition missing ')'
		"while true) {\n}\n",                        // while: missing '('
		"while (true {\n}\n",                        // while: missing ')'
		"while ()) {\n}\n",                          // while condition parse nil
		"while (true) return\n",                     // while body parse nil
		"for true\n",                                // for: missing '('
		"for (x = 1; true; x = 2) {\n}\n",           // for init must be let
		"for (let = 0; true; i = i + 1) {\n}\n",     // for init parse nil
		"for (let i = 0 true; i = i + 1) {\n}\n",    // for: missing first ';'
		"for (let i = 0; ; i = i + 1) {\n}\n",       // for condition parse nil
		"for (let i = 0; true i = i + 1) {\n}\n",    // for: missing second ';'
		"for (let i = 0; true; ) {\n}\n",            // for update parse nil
		"for (let i = 0; true; i = i + 1 {\n}\n",    // for: missing ')'
		"for (let i = 0; true; i = i + 1) return\n", // for body parse nil
		"return )\n",                                // return expr parse nil
		"{\nlet x = 1\n",                            // block: missing closing brace
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

func TestRequireStatementTerminatorAllowsEOF(t *testing.T) {
	p := New("let x = 1")
	_, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected errors at EOF-terminated statement: %+v", errs)
	}
}
