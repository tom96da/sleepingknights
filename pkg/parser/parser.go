package parser

import (
	"fmt"

	"github.com/tom96da/sleepingknights/pkg/ast"
	"github.com/tom96da/sleepingknights/pkg/lexer"
	"github.com/tom96da/sleepingknights/pkg/token"
)

const (
	precLowest = iota
	precAssign
	precOr
	precAnd
	precEquality
	precComparison
	precAdditive
	precMultiplicative
	precUnary
	precCall
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:  precAssign,
	token.OR:      precOr,
	token.AND:     precAnd,
	token.EQ:      precEquality,
	token.NOT_EQ:  precEquality,
	token.LT:      precComparison,
	token.LE:      precComparison,
	token.GT:      precComparison,
	token.GE:      precComparison,
	token.PLUS:    precAdditive,
	token.MINUS:   precAdditive,
	token.STAR:    precMultiplicative,
	token.SLASH:   precMultiplicative,
	token.PERCENT: precMultiplicative,
	token.LPAREN:  precCall,
}

// Parser parses source tokens into an AST.
type Parser struct {
	tokens []token.Token
	index  int
	errors []ParseError
}

// New builds a parser from source text.
func New(source string) *Parser {
	l := lexer.New(source)
	toks := make([]token.Token, 0, 256)
	errs := make([]ParseError, 0)

	for {
		tok := l.NextToken()
		if tok.Type == token.LINE_COMMENT || tok.Type == token.BLOCK_COMMENT {
			continue
		}
		if tok.Type == token.ILLEGAL {
			errs = append(errs, ParseError{Line: tok.Line, Column: tok.Column, Message: tok.Literal})
		}
		toks = append(toks, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return &Parser{tokens: toks, errors: errs}
}

// ParseProgram parses all function declarations followed by statements.
func (p *Parser) ParseProgram() (*ast.Program, []ParseError) {
	program := &ast.Program{}
	seenStatement := false

	p.skipNewlines()
	for !p.check(token.EOF) {
		if p.check(token.DOC_COMMENT) {
			fn := p.parseFunctionDecl()
			if fn != nil {
				if seenStatement {
					p.addError(p.current(), "function declarations must appear before top-level statements")
				}
				program.Functions = append(program.Functions, fn)
			}
			p.skipNewlines()
			continue
		}

		if p.check(token.FN) {
			fn := p.parseFunctionDecl()
			if fn != nil {
				if seenStatement {
					p.addError(p.current(), "function declarations must appear before top-level statements")
				}
				program.Functions = append(program.Functions, fn)
			}
			p.skipNewlines()
			continue
		}

		seenStatement = true
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.skipNewlines()
	}

	return program, p.errors
}

func (p *Parser) parseFunctionDecl() *ast.FunctionDecl {
	doc := ""
	if p.match(token.DOC_COMMENT) {
		doc = p.previous().Literal
		p.skipNewlines()
	}

	if !p.match(token.FN) {
		p.addError(p.current(), "expected 'fn' after doc comment")
		p.synchronize()
		return nil
	}
	start := p.previous()

	if !p.match(token.IDENT) {
		p.addError(p.current(), "expected function name")
		p.synchronize()
		return nil
	}
	nameTok := p.previous()

	if !p.match(token.LPAREN) {
		p.addError(p.current(), "expected '(' after function name")
		p.synchronize()
		return nil
	}

	params := make([]*ast.Parameter, 0, 4)
	if !p.check(token.RPAREN) {
		for {
			if !p.match(token.IDENT) {
				p.addError(p.current(), "expected parameter name")
				p.synchronize()
				return nil
			}
			paramName := p.previous()

			if !p.match(token.COLON) {
				p.addError(p.current(), "expected ':' after parameter name")
				p.synchronize()
				return nil
			}

			paramType := p.parseTypeName()
			if paramType == nil {
				p.synchronize()
				return nil
			}

			params = append(params, &ast.Parameter{
				Name:   paramName.Literal,
				Type:   paramType,
				Line:   paramName.Line,
				Column: paramName.Column,
			})

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	if !p.match(token.RPAREN) {
		p.addError(p.current(), "expected ')' after parameters")
		p.synchronize()
		return nil
	}

	var returnType *ast.TypeName
	if p.match(token.COLON) {
		returnType = p.parseTypeName()
		if returnType == nil {
			p.synchronize()
			return nil
		}
	}

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	return &ast.FunctionDecl{
		DocComment: doc,
		Name:       &ast.Identifier{Name: nameTok.Literal, Line: nameTok.Line, Column: nameTok.Column},
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
		Line:       start.Line,
		Column:     start.Column,
	}
}

func (p *Parser) parseTypeName() *ast.TypeName {
	tok := p.current()
	switch tok.Type {
	case token.INT_TYPE, token.FLOAT_TYPE, token.BOOL_TYPE, token.STRING_TYPE, token.VOID_TYPE:
		p.advance()
		return &ast.TypeName{Name: tok.Literal, Line: tok.Line, Column: tok.Column}
	default:
		p.addError(tok, fmt.Sprintf("expected type name, got %s", tok.Type))
		return nil
	}
}

func (p *Parser) current() token.Token {
	if p.index >= len(p.tokens) {
		return token.Token{Type: token.EOF, Line: 1, Column: 1}
	}
	return p.tokens[p.index]
}

func (p *Parser) previous() token.Token {
	if p.index == 0 {
		return token.Token{Type: token.EOF, Line: 1, Column: 1}
	}
	return p.tokens[p.index-1]
}

func (p *Parser) check(tt token.TokenType) bool {
	return p.current().Type == tt
}

func (p *Parser) advance() token.Token {
	if p.index < len(p.tokens) {
		p.index++
	}
	return p.previous()
}

func (p *Parser) match(tt token.TokenType) bool {
	if p.check(tt) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) addError(tok token.Token, msg string) {
	p.errors = append(p.errors, ParseError{Line: tok.Line, Column: tok.Column, Message: msg})
}

func (p *Parser) skipNewlines() {
	for p.match(token.NEWLINE) {
	}
}

func (p *Parser) synchronize() {
	for !p.check(token.EOF) {
		if p.check(token.NEWLINE) {
			p.advance()
			p.skipNewlines()
			return
		}
		if p.check(token.RBRACE) {
			return
		}
		p.advance()
	}
}

func (p *Parser) isExpressionTerminator(tt token.TokenType) bool {
	switch tt {
	case token.NEWLINE, token.EOF, token.RPAREN, token.RBRACE, token.SEMICOLON, token.COMMA:
		return true
	default:
		return false
	}
}

func (p *Parser) currentPrecedence() int {
	if prec, ok := precedences[p.current().Type]; ok {
		return prec
	}
	return precLowest
}
