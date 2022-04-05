package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_IF_CONDITION_START_LPAREN = "Expect if condition to start with '('."
const ERR_IF_CONDITION_END_RPAREN = "Expect if condition to end with ')'."

func (p *Parser) parseIfExpr() ast.Expression {
	ifExpr := &ast.IfExpr{
		Token: p.currToken,
	}

	if !p.expectPeek(token.LPAREN, ERR_IF_CONDITION_START_LPAREN) {
		return nil
	}
	p.nextToken()
	ifExpr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN, ERR_IF_CONDITION_END_RPAREN) {
		return nil
	}
	p.nextToken()

	ifExpr.Then = p.parseStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()
		ifExpr.Else = p.parseStatement()
	}

	return ifExpr
}
