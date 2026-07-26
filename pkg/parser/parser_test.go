package parser

import (
	"testing"

	"github.com/tom96da/sleepingknights/pkg/ast"
	"github.com/tom96da/sleepingknights/pkg/token"
)

func TestNewSkipsCommentsAndCollectsIllegalTokens(t *testing.T) {
	src := "// line\n/* block */\n@\nlet x = 1\n"
	p := New(src)
	program, errs := p.ParseProgram()
	if len(errs) == 0 {
		t.Fatal("expected at least one error from illegal token")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("unexpected statement count after skipping comments: got=%d", len(program.Statements))
	}
}

func TestParseProgramFunctionDeclWithDocAndTypes(t *testing.T) {
	src := "/** doc */\nfn add(a: int, b: float): void {\nreturn\n}\n"
	p := New(src)
	program, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}
	if len(program.Functions) != 1 {
		t.Fatalf("unexpected function count: got=%d", len(program.Functions))
	}
	if program.Functions[0].DocComment == "" {
		t.Fatal("expected doc comment")
	}
}

func TestParseProgramFunctionOrderConstraint(t *testing.T) {
	src1 := "let x = 1\nfn later() {\nreturn\n}\n"
	p1 := New(src1)
	_, errs := p1.ParseProgram()
	if len(errs) == 0 {
		t.Fatal("expected function order error")
	}

	src2 := "let x = 1\n/** doc */\nfn later() {\nreturn\n}\n"
	p2 := New(src2)
	_, errs = p2.ParseProgram()
	if len(errs) == 0 {
		t.Fatal("expected function order error with doc comment")
	}
}

func TestParseFunctionDeclErrorBranches(t *testing.T) {
	cases := []string{
		"/** doc */\nlet x = 1\n",     // doc comment then missing fn
		"fn\n",                        // missing function name
		"fn f\n",                      // missing '('
		"fn f(\n",                     // missing parameter name
		"fn f(a int) {\nreturn\n}\n",  // missing ':' in parameter
		"fn f(a: bad) {\nreturn\n}\n", // invalid parameter type
		"fn f(a: int {\n}\n",          // missing ')' after parameters
		"fn f(): bad {\nreturn\n}\n",  // invalid return type
		"fn f()\n",                    // missing function block
	}

	for i, src := range cases {
		t.Run("case", func(t *testing.T) {
			p := New(src)
			_, errs := p.ParseProgram()
			if len(errs) == 0 {
				t.Fatalf("case %d expected errors", i)
			}
		})
	}
}

func TestParserHelperBranches(t *testing.T) {
	p := &Parser{}
	if p.current().Type != token.EOF {
		t.Fatal("current() at empty parser must be EOF")
	}
	if p.previous().Type != token.EOF {
		t.Fatal("previous() at index 0 must be EOF")
	}

	p.tokens = []token.Token{{Type: token.IDENT, Literal: "x", Line: 2, Column: 3}, {Type: token.EOF, Line: 2, Column: 4}}
	p.index = 1
	if p.previous().Type != token.IDENT {
		t.Fatal("previous() should return prior token")
	}
	p.index = 99
	if p.current().Type != token.EOF {
		t.Fatal("current() out-of-range must be EOF")
	}

	p2 := &Parser{tokens: []token.Token{{Type: token.RBRACE, Line: 1, Column: 1}, {Type: token.EOF, Line: 1, Column: 2}}}
	p2.synchronize()
	if p2.index != 0 {
		t.Fatalf("synchronize() should stop on RBRACE: got index=%d", p2.index)
	}
}

func TestParseIfLet(t *testing.T) {
	src := "if (let x = 1) {\nprint(x)\n}\n"
	p := New(src)
	program, errs := p.ParseProgram()
	if len(errs) > 0 {
		t.Fatalf("unexpected parse errors: %+v", errs)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("unexpected statement count: got=%d", len(program.Statements))
	}
	ifStmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("expected if statement, got %T", program.Statements[0])
	}
	if ifStmt.LetBinding == nil {
		t.Fatal("expected let binding in if condition")
	}
}

func TestParseErrorsRecovered(t *testing.T) {
	src := "let x =\nlet y = 1\n"
	p := New(src)
	_, errs := p.ParseProgram()
	if len(errs) == 0 {
		t.Fatal("expected parse errors")
	}
}
