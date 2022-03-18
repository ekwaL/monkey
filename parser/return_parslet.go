package parser

import "monkey/ast"

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	r := &ast.ReturnStmt{Token: p.currToken }

	p.nextToken()
	r.Value = p.parseExpression()

	return r
}
