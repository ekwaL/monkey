package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_GET_NO_PROP_NAME = "Expect property name after '.' ."

func (p *Parser) parseGetExpr(left ast.Expression) ast.Expression {
	expr := &ast.GetExpr{
		Token:      p.currToken,
		Expression: left,
	}

	if !p.expectPeek(token.IDENTIFIER, ERR_GET_NO_PROP_NAME) {
		return nil
	}
	expr.Field = p.parseIdentifierExpr().(*ast.IdentifierExpr)

	return expr
}
