package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseExpressionStmt() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{
		Token: p.currToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
