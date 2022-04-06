package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_CALL_ARGUMENTS_END_RPAREN = "Expect function call arguments list to end with ')'."

func (p *Parser) parseCallExpr(left ast.Expression) ast.Expression {
	expr := &ast.CallExpr{
		Token:     p.currToken,
		Function:  left,
	}

	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RPAREN) {
		p.nextToken()

		expr.Arguments = append(expr.Arguments, p.parseExpression(LOWEST))

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.RPAREN, ERR_CALL_ARGUMENTS_END_RPAREN) {
		return nil
	}

	return expr
}
