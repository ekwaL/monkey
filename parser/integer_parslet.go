package parser

import (
	"monkey/ast"
	"strconv"
)

func (p *Parser) parseIntLiteralExpr() *ast.IntLiteralExpr {
	if val, err := strconv.Atoi(p.currToken.Literal); err == nil {
		tok := p.currToken
		p.nextToken()
		return &ast.IntLiteralExpr{
			Token: tok,
			Value: val,
		}
	}
	panic("Fatal: Can not parse Int literal token value.")
}
