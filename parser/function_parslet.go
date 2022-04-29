package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

const ERR_FN_PARAMETERS_START_LPAREN = "Expect function parameters list to start with '('."
const ERR_FN_PARAMETERS_END_RPAREN = "Expect function parameters list to end with ')'."
const ERR_FN_PARAMETER_SHOULD_BE_IDENTIFIER = "Wrong function parameter %q. Expect function parameters to be an identifiers."
const ERR_FN_BODY_START_LBRACE = "Expect function body to start with '{'."
const ERR_FN_BODY_END_RBRACE = "Expect function body to end with '}'."

func (p *Parser) parseFunctionExpr() ast.Expression {
	fn := &ast.FunctionExpr{
		Token: p.currToken,
	}

	if !p.expectPeek(token.LPAREN, ERR_FN_PARAMETERS_START_LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.RPAREN, ERR_FN_PARAMETERS_END_RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE, ERR_FN_BODY_START_LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStmt()

	if p.currToken.Type != token.RBRACE {
		p.error(ERR_FN_BODY_END_RBRACE)
		return nil
	}

	return fn
}

func (p *Parser) parseFunctionDefinition() *ast.LetStmt {
	fn := &ast.FunctionExpr{
		Token: p.currToken,
	}

	p.nextToken()
	name := p.parseIdentifierExpr().(*ast.IdentifierExpr)

	if !p.expectPeek(token.LPAREN, ERR_FN_PARAMETERS_START_LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.RPAREN, ERR_FN_PARAMETERS_END_RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE, ERR_FN_BODY_START_LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStmt()

	if p.currToken.Type != token.RBRACE {
		p.error(ERR_FN_BODY_END_RBRACE)
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &ast.LetStmt{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
		},
		Name:  name,
		Value: fn,
	}
}

func (p *Parser) parseFunctionParameters() (params []*ast.IdentifierExpr) {
	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RPAREN) {
		p.nextToken()

		if p.currToken.Type != token.IDENTIFIER {
			p.error(fmt.Sprintf(ERR_FN_PARAMETER_SHOULD_BE_IDENTIFIER, p.currToken.Literal))
			return nil
		} else {
			params = append(params, p.parseIdentifierExpr().(*ast.IdentifierExpr))
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return
}
