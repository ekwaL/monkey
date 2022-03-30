package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// operator precedence
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > || <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X || !x
	CALL        // function()
)

type (
	prefixParslet func() ast.Expression
	infixParslet  func(left ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors   []string
	needSync bool

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

	p.prefixParslets = make(map[token.TokenType]prefixParslet)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifierExpr)
	p.registerPrefix(token.INT, p.parseIntLiteralExpr)

	return &p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}

	for p.currToken.Type != token.EOF {

		stmt := p.parseStatement()

		if p.needSync {
			p.synchronize()
		} else if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)

			p.nextToken()
		}
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
		return p.parseExpressionStmt()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParslets[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	leftExpr := prefix()
	return leftExpr
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
	p.needSync = true
}

func (p *Parser) synchronize() {
	p.needSync = false
	p.nextToken()
	for p.currToken.Type != token.EOF {
		switch p.currToken.Type {
		case token.SEMICOLON:
			p.nextToken()
			return
		case token.LET, token.FUNCTION, token.RETURN, token.IF:
			return
		default:
			p.nextToken()
		}
	}
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

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParslet) {
	p.prefixParslets[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParslet) {
	p.infixParslets[tt] = fn

}
