package parser

import (
	"monkey/ast"
)

const ERR_WRONG_ASSIGNMENT_TARGET = "Assignment target should be an identifier or field."

func (p *Parser) parseAssignExpr(left ast.Expression) ast.Expression {
	switch node := left.(type) {
	case *ast.IdentifierExpr:
		expr := &ast.AssignExpr{
			Token:      p.currToken,
			Identifier: node,
		}
		p.nextToken()
		expr.Expression = p.parseExpression(LOWEST)

		return expr

	case *ast.GetExpr:
		expr := &ast.SetExpr{
			Token:      p.currToken,
			Expression: node.Expression,
			Field:      node.Field,
		}
		p.nextToken()
		expr.Value = p.parseExpression(LOWEST)

		return expr

	default:
		p.error(ERR_WRONG_ASSIGNMENT_TARGET)
		return nil
	}
}
