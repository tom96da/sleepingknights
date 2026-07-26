package parser

import (
	"github.com/tom96da/sleepingknights/pkg/ast"
	"github.com/tom96da/sleepingknights/pkg/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.current().Type {
	case token.LET:
		return p.parseVariableDecl(false, true)
	case token.CONST:
		return p.parseVariableDecl(true, true)
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.BREAK:
		tok := p.advance()
		stmt := &ast.BreakStatement{Line: tok.Line, Column: tok.Column}
		p.requireStatementTerminator()
		return stmt
	case token.CONTINUE:
		tok := p.advance()
		stmt := &ast.ContinueStatement{Line: tok.Line, Column: tok.Column}
		p.requireStatementTerminator()
		return stmt
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVariableDecl(isConst bool, requireTerminator bool) *ast.VariableDecl {
	start := p.advance()
	if !p.match(token.IDENT) {
		p.addError(p.current(), "expected identifier in declaration")
		p.synchronize()
		return nil
	}
	nameTok := p.previous()

	if !p.match(token.ASSIGN) {
		p.addError(p.current(), "expected '=' in declaration")
		p.synchronize()
		return nil
	}

	value := p.parseExpression(precLowest)
	if value == nil {
		p.synchronize()
		return nil
	}

	decl := &ast.VariableDecl{
		IsConst: isConst,
		Name:    &ast.Identifier{Name: nameTok.Literal, Line: nameTok.Line, Column: nameTok.Column},
		Value:   value,
		Line:    start.Line,
		Column:  start.Column,
	}

	if requireTerminator {
		p.requireStatementTerminator()
	}
	return decl
}

func (p *Parser) parseIfStatement() ast.Statement {
	start := p.advance()
	if !p.match(token.LPAREN) {
		p.addError(p.current(), "expected '(' after 'if'")
		p.synchronize()
		return nil
	}

	var cond ast.Expression
	var letBinding *ast.VariableDecl
	if p.check(token.LET) {
		letBinding = p.parseVariableDecl(false, false)
		if letBinding == nil {
			return nil
		}
	} else {
		cond = p.parseExpression(precLowest)
		if cond == nil {
			return nil
		}
	}

	if !p.match(token.RPAREN) {
		p.addError(p.current(), "expected ')' after if condition")
		p.synchronize()
		return nil
	}

	thenBlock := p.parseBlockStatement()
	if thenBlock == nil {
		return nil
	}

	p.skipNewlines()
	var elseStmt ast.Statement
	if p.match(token.ELSE) {
		p.skipNewlines()
		if p.check(token.IF) {
			elseStmt = p.parseIfStatement()
		} else {
			elseStmt = p.parseBlockStatement()
		}
	}

	return &ast.IfStatement{
		Condition:  cond,
		LetBinding: letBinding,
		Then:       thenBlock,
		Else:       elseStmt,
		Line:       start.Line,
		Column:     start.Column,
	}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	start := p.advance()
	if !p.match(token.LPAREN) {
		p.addError(p.current(), "expected '(' after 'while'")
		p.synchronize()
		return nil
	}

	cond := p.parseExpression(precLowest)
	if cond == nil {
		return nil
	}
	if !p.match(token.RPAREN) {
		p.addError(p.current(), "expected ')' after while condition")
		p.synchronize()
		return nil
	}

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	return &ast.WhileStatement{Condition: cond, Body: body, Line: start.Line, Column: start.Column}
}

func (p *Parser) parseForStatement() ast.Statement {
	start := p.advance()
	if !p.match(token.LPAREN) {
		p.addError(p.current(), "expected '(' after 'for'")
		p.synchronize()
		return nil
	}

	if !p.check(token.LET) {
		p.addError(p.current(), "for initializer must start with 'let'")
		p.synchronize()
		return nil
	}
	init := p.parseVariableDecl(false, false)
	if init == nil {
		return nil
	}
	if !p.match(token.SEMICOLON) {
		p.addError(p.current(), "expected ';' after for initializer")
		p.synchronize()
		return nil
	}

	cond := p.parseExpression(precLowest)
	if cond == nil {
		return nil
	}
	if !p.match(token.SEMICOLON) {
		p.addError(p.current(), "expected ';' after for condition")
		p.synchronize()
		return nil
	}

	update := p.parseExpression(precLowest)
	if update == nil {
		return nil
	}
	if !p.match(token.RPAREN) {
		p.addError(p.current(), "expected ')' after for update")
		p.synchronize()
		return nil
	}

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	return &ast.ForStatement{Init: init, Condition: cond, Update: update, Body: body, Line: start.Line, Column: start.Column}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	start := p.advance()
	if p.check(token.NEWLINE) || p.check(token.RBRACE) || p.check(token.EOF) {
		stmt := &ast.ReturnStatement{Value: nil, Line: start.Line, Column: start.Column}
		p.requireStatementTerminator()
		return stmt
	}

	value := p.parseExpression(precLowest)
	if value == nil {
		p.synchronize()
		return nil
	}

	stmt := &ast.ReturnStatement{Value: value, Line: start.Line, Column: start.Column}
	p.requireStatementTerminator()
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	start := p.current()
	expr := p.parseExpression(precLowest)
	if expr == nil {
		p.synchronize()
		return nil
	}
	stmt := &ast.ExpressionStatement{Expr: expr, Line: start.Line, Column: start.Column}
	p.requireStatementTerminator()
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' to start block")
		p.synchronize()
		return nil
	}
	start := p.previous()
	block := &ast.BlockStatement{Statements: make([]ast.Statement, 0, 8), Line: start.Line, Column: start.Column}

	p.skipNewlines()
	for !p.check(token.RBRACE) && !p.check(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.skipNewlines()
	}

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}' to close block")
		return nil
	}

	return block
}

func (p *Parser) requireStatementTerminator() {
	if p.match(token.NEWLINE) {
		p.skipNewlines()
		return
	}
	if p.check(token.EOF) || p.check(token.RBRACE) {
		return
	}
	p.addError(p.current(), "expected newline after statement")
	p.synchronize()
}
