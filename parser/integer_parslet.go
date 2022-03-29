package parser

import (
	"fmt"
	"monkey/ast"
	"strconv"
)

const ERR_COULD_NOT_PARSE_INT = "Could not parse %q as integer: %v"

func (p *Parser) parseIntLiteralExpr() ast.Expression {
	if val, err := strconv.ParseInt(p.currToken.Literal, 0, 64); err == nil {
		tok := p.currToken
		p.nextToken()
		return &ast.IntLiteralExpr{
			Token: tok,
			Value: val,
		}
	} else {
		msg := fmt.Sprintf(ERR_COULD_NOT_PARSE_INT, p.currToken.Literal, err.Error())
		p.error(msg)
		return nil
	}
}
