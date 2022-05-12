package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_INDEX_END_BRACKET = "Expect index expression to end with ']'."

func (p *Parser) parseIndexExpr(left ast.Expression) ast.Expression {
	idx := &ast.IndexExpr{
		Token: p.currToken,
		Left:  left,
		Index: nil,
	}

	p.nextToken()

	idx.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET, ERR_INDEX_END_BRACKET) {
		return nil
	}

	return idx
}
