package parser

import "monkey/ast"

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpr{
		Token:    p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}
