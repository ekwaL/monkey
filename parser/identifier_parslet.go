package parser

import "monkey/ast"

func (p *Parser) parseIdentifierExpr() ast.Expression {
	identifier := p.currToken

	return &ast.IdentifierExpr{
		Token: identifier,
		Value: identifier.Literal,
	}
}
