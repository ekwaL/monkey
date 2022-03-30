package parser

import (
	"fmt"
	"monkey/ast"
	"strconv"
)

const ERR_COULD_NOT_PARSE_INT = "Could not parse %q as integer: %v"

func (p *Parser) parseIntLiteralExpr() ast.Expression {
	if val, err := strconv.ParseInt(p.currToken.Literal, 0, 64); err == nil {
		return &ast.IntLiteralExpr{
			Token: p.currToken,
			Value: val,
		}
	} else {
		msg := fmt.Sprintf(ERR_COULD_NOT_PARSE_INT, p.currToken.Literal, err.Error())
		p.error(msg)
		return nil
	}
}
