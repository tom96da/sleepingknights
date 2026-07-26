package lexer

import (
	"testing"

	"github.com/tom96da/sleepingknights/pkg/token"
)

type expectedToken struct {
	typeValue token.TokenType
	literal   string
	line      int
	column    int
}

func collectTokens(input string) []token.Token {
	l := New(input)
	tokens := make([]token.Token, 0, 32)

	for range 512 {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}

func assertTokens(t *testing.T, got []token.Token, want []expectedToken) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d", len(got), len(want))
	}

	for i := range want {
		if got[i].Type != want[i].typeValue || got[i].Literal != want[i].literal {
			t.Fatalf("token[%d] mismatch: got=(%s,%q) want=(%s,%q)", i, got[i].Type, got[i].Literal, want[i].typeValue, want[i].literal)
		}
		if want[i].line > 0 && got[i].Line != want[i].line {
			t.Fatalf("token[%d] line mismatch: got=%d want=%d", i, got[i].Line, want[i].line)
		}
		if want[i].column > 0 && got[i].Column != want[i].column {
			t.Fatalf("token[%d] column mismatch: got=%d want=%d", i, got[i].Column, want[i].column)
		}
	}
}

func TestLexerOperatorsAndDelimiters(t *testing.T) {
	input := "+ - * / % = == != ! < <= > >= ( ) { } , : ;"
	got := collectTokens(input)

	want := []expectedToken{
		{token.PLUS, "+", 0, 0},
		{token.MINUS, "-", 0, 0},
		{token.STAR, "*", 0, 0},
		{token.SLASH, "/", 0, 0},
		{token.PERCENT, "%", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.EQ, "==", 0, 0},
		{token.NOT_EQ, "!=", 0, 0},
		{token.ILLEGAL, "!", 0, 0},
		{token.LT, "<", 0, 0},
		{token.LE, "<=", 0, 0},
		{token.GT, ">", 0, 0},
		{token.GE, ">=", 0, 0},
		{token.LPAREN, "(", 0, 0},
		{token.RPAREN, ")", 0, 0},
		{token.LBRACE, "{", 0, 0},
		{token.RBRACE, "}", 0, 0},
		{token.COMMA, ",", 0, 0},
		{token.COLON, ":", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.EOF, "", 0, 0},
	}

	assertTokens(t, got, want)
}

func TestLexerKeywordsAndIdentifiers(t *testing.T) {
	input := "let const fn if else for while return break continue true false int float bool string void and or not foo bar_1"
	got := collectTokens(input)

	want := []expectedToken{
		{token.LET, "let", 0, 0},
		{token.CONST, "const", 0, 0},
		{token.FN, "fn", 0, 0},
		{token.IF, "if", 0, 0},
		{token.ELSE, "else", 0, 0},
		{token.FOR, "for", 0, 0},
		{token.WHILE, "while", 0, 0},
		{token.RETURN, "return", 0, 0},
		{token.BREAK, "break", 0, 0},
		{token.CONTINUE, "continue", 0, 0},
		{token.TRUE, "true", 0, 0},
		{token.FALSE, "false", 0, 0},
		{token.INT_TYPE, "int", 0, 0},
		{token.FLOAT_TYPE, "float", 0, 0},
		{token.BOOL_TYPE, "bool", 0, 0},
		{token.STRING_TYPE, "string", 0, 0},
		{token.VOID_TYPE, "void", 0, 0},
		{token.AND, "and", 0, 0},
		{token.OR, "or", 0, 0},
		{token.NOT, "not", 0, 0},
		{token.IDENT, "foo", 0, 0},
		{token.IDENT, "bar_1", 0, 0},
		{token.EOF, "", 0, 0},
	}

	assertTokens(t, got, want)
}

func TestLexerNumberLiterals(t *testing.T) {
	t.Run("valid numbers", func(t *testing.T) {
		input := "123 45.67 1e10 2E-3 0x1Af 0b10101"
		got := collectTokens(input)

		want := []expectedToken{
			{token.INT, "123", 0, 0},
			{token.FLOAT, "45.67", 0, 0},
			{token.FLOAT, "1e10", 0, 0},
			{token.FLOAT, "2E-3", 0, 0},
			{token.INT, "0x1Af", 0, 0},
			{token.INT, "0b10101", 0, 0},
			{token.EOF, "", 0, 0},
		}

		assertTokens(t, got, want)
	})

	t.Run("invalid numbers", func(t *testing.T) {
		input := "0x 0b 1e+"
		got := collectTokens(input)

		want := []expectedToken{
			{token.ILLEGAL, "0x", 0, 0},
			{token.ILLEGAL, "0b", 0, 0},
			{token.ILLEGAL, "1e+", 0, 0},
			{token.EOF, "", 0, 0},
		}

		assertTokens(t, got, want)
	})
}

func TestLexerStringLiterals(t *testing.T) {
	t.Run("valid string with escapes", func(t *testing.T) {
		input := "\"a\\nb\\tc\\rd\\be\\ff\\\"g\\\\\""
		got := collectTokens(input)

		want := []expectedToken{
			{token.STRING, "a\nb\tc\rd\be\ff\"g\\", 0, 0},
			{token.EOF, "", 0, 0},
		}

		assertTokens(t, got, want)
	})

	t.Run("unterminated string", func(t *testing.T) {
		input := "\"unterminated"
		got := collectTokens(input)

		want := []expectedToken{
			{token.ILLEGAL, "unterminated string", 0, 0},
			{token.EOF, "", 0, 0},
		}

		assertTokens(t, got, want)
	})

	t.Run("invalid escape", func(t *testing.T) {
		input := "\"bad\\q\""
		got := collectTokens(input)

		if got[0].Type != token.ILLEGAL || got[0].Literal != "invalid escape sequence" {
			t.Fatalf("first token mismatch: got=(%s,%q)", got[0].Type, got[0].Literal)
		}
	})
}

func TestLexerComments(t *testing.T) {
	t.Run("line and block comments", func(t *testing.T) {
		input := "// hello\n/* block */ /** doc */"
		got := collectTokens(input)

		want := []expectedToken{
			{token.LINE_COMMENT, "// hello", 0, 0},
			{token.NEWLINE, "\\n", 0, 0},
			{token.BLOCK_COMMENT, "/* block */", 0, 0},
			{token.DOC_COMMENT, "/** doc */", 0, 0},
			{token.EOF, "", 0, 0},
		}

		assertTokens(t, got, want)
	})

	t.Run("unterminated block comment", func(t *testing.T) {
		input := "/* missing end"
		got := collectTokens(input)

		want := []expectedToken{
			{token.ILLEGAL, "unterminated block comment", 0, 0},
			{token.EOF, "", 0, 0},
		}

		assertTokens(t, got, want)
	})
}

func TestLexerPositionsAndWhitespace(t *testing.T) {
	input := "let x = 1\n\tconst y=2\n"
	got := collectTokens(input)

	want := []expectedToken{
		{token.LET, "let", 1, 1},
		{token.IDENT, "x", 1, 5},
		{token.ASSIGN, "=", 1, 7},
		{token.INT, "1", 1, 9},
		{token.NEWLINE, "\\n", 1, 10},
		{token.CONST, "const", 2, 2},
		{token.IDENT, "y", 2, 8},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "2", 2, 10},
		{token.NEWLINE, "\\n", 2, 11},
		{token.EOF, "", 3, 1},
	}

	assertTokens(t, got, want)
}

func TestLexerIllegalRune(t *testing.T) {
	input := "@"
	got := collectTokens(input)

	want := []expectedToken{
		{token.ILLEGAL, "@", 0, 0},
		{token.EOF, "", 0, 0},
	}

	assertTokens(t, got, want)
}
