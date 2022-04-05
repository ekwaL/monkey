package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_GROUPING_RIGHT_PAREN_MISSING = "Expecting ')'."

func (p *Parser) parseGroupingExpr() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN, ERR_GROUPING_RIGHT_PAREN_MISSING) {
		return nil
	}

	return expr
}
