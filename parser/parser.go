package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParslet func() ast.Expression
	infixParslet  func(left ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParslets map[token.TokenType]prefixParslet
	infixParslets  map[token.TokenType]infixParslet
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return &p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}

		p.nextToken()
	}

	return prog
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return nil
	}
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.currToken.Type {
	case token.IDENTIFIER:
		return p.parseIdentifierExpr()
	case token.INT:
		return p.parseIntLiteralExpr()
	default:
		return nil
	}
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// func (p *Parser) currTokenIs(tt token.TokenType) bool {
// 	return p.currToken.Type == tt
// }

func (p *Parser) peekTokenIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(tt token.TokenType, errMsg string) bool {
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	} else {
		p.error(errMsg)
		return false
	}
}
