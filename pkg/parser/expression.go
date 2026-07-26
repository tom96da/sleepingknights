package parser

import (
	"fmt"

	"github.com/tom96da/sleepingknights/pkg/ast"
	"github.com/tom96da/sleepingknights/pkg/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	left := p.parsePrefixExpression()
	if left == nil {
		return nil
	}

	for !p.isExpressionTerminator(p.current().Type) && precedence < p.currentPrecedence() {
		op := p.current()
		p.advance()
		left = p.parseInfixExpression(left, op)
		if left == nil {
			return nil
		}
	}

	return left
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	tok := p.current()

	switch tok.Type {
	case token.IDENT:
		p.advance()
		return &ast.Identifier{Name: tok.Literal, Line: tok.Line, Column: tok.Column}
	case token.INT:
		p.advance()
		return &ast.IntegerLiteral{Value: tok.Literal, Line: tok.Line, Column: tok.Column}
	case token.FLOAT:
		p.advance()
		return &ast.FloatLiteral{Value: tok.Literal, Line: tok.Line, Column: tok.Column}
	case token.STRING:
		p.advance()
		return &ast.StringLiteral{Value: tok.Literal, Line: tok.Line, Column: tok.Column}
	case token.TRUE:
		p.advance()
		return &ast.BooleanLiteral{Value: true, Line: tok.Line, Column: tok.Column}
	case token.FALSE:
		p.advance()
		return &ast.BooleanLiteral{Value: false, Line: tok.Line, Column: tok.Column}
	case token.MINUS, token.NOT:
		p.advance()
		right := p.parseExpression(precUnary)
		if right == nil {
			return nil
		}
		return &ast.UnaryExpression{Operator: tok.Type, Operand: right, Line: tok.Line, Column: tok.Column}
	case token.LPAREN:
		p.advance()
		expr := p.parseExpression(precLowest)
		if expr == nil {
			return nil
		}
		if !p.match(token.RPAREN) {
			p.addError(p.current(), "expected ')' to close grouped expression")
			return nil
		}
		return expr
	default:
		p.addError(tok, fmt.Sprintf("unexpected token in expression: %s", tok.Type))
		p.advance()
		return nil
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression, op token.Token) ast.Expression {
	switch op.Type {
	case token.LPAREN:
		ident, ok := left.(*ast.Identifier)
		if !ok {
			p.addError(op, "function call target must be an identifier")
			p.skipCallArguments()
			return nil
		}
		args := make([]ast.Expression, 0, 4)
		if !p.check(token.RPAREN) {
			for {
				arg := p.parseExpression(precLowest)
				if arg == nil {
					return nil
				}
				args = append(args, arg)
				if p.match(token.COMMA) {
					continue
				}
				break
			}
		}
		if !p.match(token.RPAREN) {
			p.addError(p.current(), "expected ')' after call arguments")
			return nil
		}
		return &ast.FunctionCall{Function: ident, Arguments: args, Line: op.Line, Column: op.Column}
	case token.ASSIGN:
		name, ok := left.(*ast.Identifier)
		if !ok {
			p.addError(op, "left side of assignment must be an identifier")
			return nil
		}
		right := p.parseExpression(precAssign - 1)
		if right == nil {
			return nil
		}
		return &ast.AssignmentExpression{Name: name, Value: right, Line: op.Line, Column: op.Column}
	case token.OR, token.AND,
		token.EQ, token.NOT_EQ,
		token.LT, token.LE, token.GT, token.GE,
		token.PLUS, token.MINUS,
		token.STAR, token.SLASH, token.PERCENT:
		prec := precedences[op.Type]
		right := p.parseExpression(prec)
		if right == nil {
			return nil
		}
		return &ast.BinaryExpression{Left: left, Operator: op.Type, Right: right, Line: op.Line, Column: op.Column}
	default:
		p.addError(op, fmt.Sprintf("unexpected infix operator: %s", op.Type))
		return nil
	}
}

func (p *Parser) skipCallArguments() {
	depth := 1
	for !p.check(token.EOF) {
		if p.match(token.LPAREN) {
			depth++
			continue
		}
		if p.match(token.RPAREN) {
			depth--
			if depth == 0 {
				return
			}
			continue
		}
		p.advance()
	}
}
