package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_ARR_LITERAL_END_BRACKET = "Expect array literal elements list to end with ']'."

func (p *Parser) parseArrayLiteralExpr() ast.Expression {
	arr := &ast.ArrayLiteralExpr{
		Token:    p.currToken,
		Elements: []ast.Expression{},
	}

	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RBRACKET) {
		p.nextToken()

		arr.Elements = append(arr.Elements, p.parseExpression(LOWEST))

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.RBRACKET, ERR_ARR_LITERAL_END_BRACKET) {
		return nil
	}

	return arr
}
