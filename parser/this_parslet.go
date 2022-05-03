package parser

import "monkey/ast"

func (p *Parser) parseThisExpr() ast.Expression {
	return &ast.ThisExpr{
		Token: p.currToken,
	}
}
