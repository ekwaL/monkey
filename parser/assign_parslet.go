package parser

import (
	"monkey/ast"
)

const ERR_WRONG_ASSIGNMENT_TARGET = "Assignment target should be an identifier."

func (p *Parser) parseAssignExpr(left ast.Expression) ast.Expression {
	identifier, ok := left.(*ast.IdentifierExpr)
	if !ok {
		p.error(ERR_WRONG_ASSIGNMENT_TARGET)
	}

	expr := &ast.AssignExpr{
		Token:      p.currToken,
		Identifier: identifier,
	}

	p.nextToken()
	// expr.Expression= p.parseExpression(ASSIGN)
	expr.Expression = p.parseExpression(LOWEST)

	return expr
}
