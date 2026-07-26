package lexer

import (
	"testing"

	"github.com/tom96da/sleepingknights/pkg/token"
)

func TestPeekEOFBranches(t *testing.T) {
	l := New("")
	if l.peekChar() != 0 {
		t.Fatalf("peekChar expected EOF")
	}
	if l.peekSecondChar() != 0 {
		t.Fatalf("peekSecondChar expected EOF")
	}
}

func TestPeekSecondNonEOFBranch(t *testing.T) {
	l := New("abc")
	if l.peekSecondChar() != 'c' {
		t.Fatalf("peekSecondChar mismatch: got=%q", l.peekSecondChar())
	}
}

func TestLexSimpleEOF(t *testing.T) {
	l := New("+")
	tok := l.NextToken()
	if tok.Type != token.PLUS {
		t.Fatalf("unexpected token: %s", tok.Type)
	}
	if l.NextToken().Type != token.EOF {
		t.Fatal("expected EOF")
	}
}
