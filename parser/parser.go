package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression

	PrefixParslet interface {
		PrefixParse() ast.Expression
	}

	InfixParslet interface {
		InfixParse(ast.Expression) ast.Expression
	}
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	currToken token.Token
	peekToken token.Token
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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tt token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead.",
		tt, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
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

func (p *Parser) currTokenIs(tt token.TokenType) bool {
	return p.currToken.Type == tt
}

func (p *Parser) peekTokenIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) expectPeek(tt token.TokenType, errMsg string) bool {
	// func (p *Parser) expectPeek(tt token.TokenType) bool {
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	} else {
		p.error(errMsg)
		// p.peekError(tt)
		return false
	}
}
