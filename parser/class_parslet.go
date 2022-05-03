package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_CLASS_NO_CLASSNAME = "Expect class name after 'class' keyword."
const ERR_CLASS_NO_SUPER_NAME = "Expect superclass name after '<' keyword."
const ERR_CLASS_BODY_START_LBRACE = "Expect class body to start with '{'."
const ERR_CLASS_BODY_END_RBRACE = "Expect class body to end with '}'."
const ERR_CLASS_WRONG_DEFINITION = "Class definition should contain only methods definitions."

func (p *Parser) parseClassStmt() *ast.ClassStmt {
	class := &ast.ClassStmt{
		Token: p.currToken,
	}

	if !p.expectPeek(token.IDENTIFIER, ERR_CLASS_NO_CLASSNAME) {
		return nil
	}

	class.Name = p.parseIdentifierExpr().(*ast.IdentifierExpr)

	if p.peekTokenIs(token.LESS) {
		p.nextToken()
		if !p.expectPeek(token.IDENTIFIER, ERR_CLASS_NO_SUPER_NAME) {
			return nil
		}

		class.Superclass = p.parseIdentifierExpr().(*ast.IdentifierExpr)
	}

	if !p.expectPeek(token.LBRACE, ERR_CLASS_BODY_START_LBRACE) {
		return nil
	}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()

		stmt := p.parseStatement()
		if method, ok := stmt.(*ast.LetStmt); ok && method != nil {
			if _, ok := method.Value.(*ast.FunctionExpr); ok {
				class.Methods = append(class.Methods, method)
			} else {
				p.error(ERR_CLASS_WRONG_DEFINITION)
			}
		} else {
			p.error(ERR_CLASS_WRONG_DEFINITION)
		}

	}

	if !p.expectPeek(token.RBRACE, ERR_CLASS_BODY_END_RBRACE) {
		return nil
	}

	return class
}
