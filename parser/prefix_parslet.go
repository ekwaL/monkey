package parser

import "monkey/ast"

func (p *Parser) parsePrefixExpr() ast.Expression {
	expr := &ast.PrefixExpr{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}
