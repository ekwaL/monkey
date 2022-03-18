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

	name := p.parseIdentifierExpr()

	if !p.expectPeek(token.ASSIGN, ERR_LET_NO_ASSIGN_AFTER_IDENTIFIER) {
		return nil
	}

	p.nextToken()
	value := p.parseExpression()

	if value == nil {
		return nil
	}

	// println(p.currToken.Literal, p.peekToken.Literal)
	// if p.expectPeek(token.SEMICOLON, ERR_LET_NO_SEMI_AFTER_LET_STMT) {
	if p.currToken.Type == token.SEMICOLON {
		// p.nextToken()
	} else {
		p.error(ERR_LET_NO_SEMI_AFTER_LET_STMT)
		return nil
	}

	return &ast.LetStmt{Token: tok, Name: name, Value: value}
}
