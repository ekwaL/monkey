package parser

import "monkey/ast"

func (p *Parser) parseNullExpr() ast.Expression {
	return &ast.NullExpr{
		Token: p.currToken,
	}
}
