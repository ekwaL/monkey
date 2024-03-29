package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	r := &ast.ReturnStmt{Token: p.currToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return r
	}

	p.nextToken()

	r.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return r
}
