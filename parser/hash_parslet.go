package parser

import (
	"monkey/ast"
	"monkey/token"
)

const (
	ERR_HASH_LITERAL_END_BRACE = "Expect hash literal pairs list to end with '|}'."
	ERR_HASH_COLON_AFTER_KEY   = "Expect ':' after hash key."
	ERR_HASH_NO_COMMA          = "Expect pairs to be separated with ','."
)

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteralExpr{
		Token: p.currToken,
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RHASHBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)
		if key == nil {
			return nil
		}

		if !p.expectPeek(token.COLON, ERR_HASH_COLON_AFTER_KEY) {
			return nil
		}
		p.nextToken()

		value := p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RHASHBRACE) && !p.expectPeek(token.COMMA, ERR_HASH_NO_COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RHASHBRACE, ERR_HASH_LITERAL_END_BRACE) {
		return nil
	}

	return hash
}
