package parser

import (
	"fmt"
	"monkey/ast"
	"strconv"
)

const ERR_COULD_NOT_PARSE_BOOL = "Could not parse %q as boolean: %v"

func (p *Parser) parseBoolLiteralExpr() ast.Expression {
	if val, err := strconv.ParseBool(p.currToken.Literal); err == nil {
		return &ast.BoolLiteralExpr{
			Token: p.currToken,
			Value: val,
		}
	} else {
		msg := fmt.Sprintf(ERR_COULD_NOT_PARSE_BOOL, p.currToken.Literal, err)
		p.error(msg)
		return nil
	}
}
