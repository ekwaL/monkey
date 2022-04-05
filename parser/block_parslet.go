package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	block := &ast.BlockStmt{
		Token:      p.currToken,
		Statements: []ast.Statement{},
	}

	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
	}

	p.nextToken()

	return block
}
