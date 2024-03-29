package parser

import (
	"monkey/ast"
	"monkey/token"
)

const ERR_LET_NO_IDENTIFIER_AFTER_LET = "Expect identifier after 'let' keyword."
const ERR_LET_NO_ASSIGN_AFTER_IDENTIFIER = "Expect '=' after identifier in 'let' statement."
const ERR_LET_NO_SEMI_AFTER_LET_STMT = "Expect ';' after 'let' statement."

func (p *Parser) parseLetStmt() *ast.LetStmt {
	tok := p.currToken

	if !p.expectPeek(token.IDENTIFIER, ERR_LET_NO_IDENTIFIER_AFTER_LET) {
		return nil
	}

	name := p.parseIdentifierExpr().(*ast.IdentifierExpr)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return &ast.LetStmt{Token: tok, Name: name, Value: nil}
	}

	if !p.expectPeek(token.ASSIGN, ERR_LET_NO_ASSIGN_AFTER_IDENTIFIER) {
		return nil
	}

	p.nextToken()
	value := p.parseExpression(LOWEST)

	if value == nil {
		return nil
	}

	if !p.expectPeek(token.SEMICOLON, ERR_LET_NO_SEMI_AFTER_LET_STMT) {
		return nil
	}

	return &ast.LetStmt{Token: tok, Name: name, Value: value}
}
