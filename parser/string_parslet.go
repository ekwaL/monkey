package parser

import "monkey/ast"

func (p *Parser) parseStringLiteralExpr() ast.Expression {
	return &ast.StringLiteralExpr{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}
