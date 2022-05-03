package parser

import (
	"monkey/ast"
	"monkey/token"
)

const (
	ERR_SUPER_NO_DOT                  = "Expect '.' after 'super' keyword."
	ERR_SUPER_NO_IDENTIFIER_AFTER_DOT = "Expect superclass method name after 'super.' ."
)

func (p *Parser) parseSuperExpr() ast.Expression {
	super := &ast.SuperExpr{
		Token: p.currToken,
	}

	if !p.expectPeek(token.DOT, ERR_SUPER_NO_DOT) {
		return nil
	}

	if !p.expectPeek(token.IDENTIFIER, ERR_SUPER_NO_IDENTIFIER_AFTER_DOT) {
		return nil
	}

	super.Method = p.parseIdentifierExpr().(*ast.IdentifierExpr)

	return super
}
