package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// operator precedence
const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	OR          // ||
	AND         // &&
	EQUALS      // ==
	LESSGREATER // > || <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X || !x
	CALL        // function()
	GET         // obj.field || arr[]
	//  isn't call/get/index expressions should have the same precedence
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:      ASSIGN,
	token.OR:          OR,
	token.AND:         AND,
	token.EQUAL_EQUAL: EQUALS,
	token.NOT_EQUAL:   EQUALS,
	token.LESS:        LESSGREATER,
	token.GREATER:     LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.STAR:        PRODUCT,
	token.LPAREN:      CALL,
	token.DOT:         GET,
	token.LBRACKET:    GET,
}

const (
	ERR_NO_PREFIX_PARSLET_FOUND = "No prefix parslet found for %q."
	ERR_ILLEGAL_TOKEN           = "Illegal token: '%s'."
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
	p.registerPrefix(token.TRUE, p.parseBoolLiteralExpr)
	p.registerPrefix(token.FALSE, p.parseBoolLiteralExpr)
	p.registerPrefix(token.NULL, p.parseNullExpr)
	p.registerPrefix(token.STRING, p.parseStringLiteralExpr)
	p.registerPrefix(token.BANG, p.parsePrefixExpr)
	p.registerPrefix(token.MINUS, p.parsePrefixExpr)
	p.registerPrefix(token.LPAREN, p.parseGroupingExpr)
	p.registerPrefix(token.IF, p.parseIfExpr)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpr)
	p.registerPrefix(token.THIS, p.parseThisExpr)
	p.registerPrefix(token.SUPER, p.parseSuperExpr)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteralExpr)

	p.infixParslets = make(map[token.TokenType]infixParslet)
	p.registerInfix(token.PLUS, p.parseInfixExpr)
	p.registerInfix(token.MINUS, p.parseInfixExpr)
	p.registerInfix(token.STAR, p.parseInfixExpr)
	p.registerInfix(token.SLASH, p.parseInfixExpr)
	p.registerInfix(token.GREATER, p.parseInfixExpr)
	p.registerInfix(token.LESS, p.parseInfixExpr)
	p.registerInfix(token.EQUAL_EQUAL, p.parseInfixExpr)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpr)
	p.registerInfix(token.OR, p.parseInfixExpr)
	p.registerInfix(token.AND, p.parseInfixExpr)
	p.registerInfix(token.LPAREN, p.parseCallExpr)
	p.registerInfix(token.ASSIGN, p.parseAssignExpr)
	p.registerInfix(token.DOT, p.parseGetExpr)
	p.registerInfix(token.LBRACKET, p.parseIndexExpr)

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
	case token.LBRACE:
		return p.parseBlockStmt()
	case token.CLASS:
		return p.parseClassStmt()
	default:
		if p.currToken.Type == token.FUNCTION && p.peekTokenIs(token.IDENTIFIER) {
			return p.parseFunctionDefinition()
		}
		return p.parseExpressionStmt()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParslets[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParsletError(p.currToken.Type)
		return nil
	}

	leftExpr := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParslets[p.peekToken.Type]

		if infix == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infix(leftExpr)
	}

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

	if p.currToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	for p.currToken.Type != token.EOF {
		switch p.currToken.Type {
		case token.SEMICOLON:
			p.nextToken()
			return
		case token.LET, token.FUNCTION, token.RETURN, token.IF, token.CLASS:
			return
		default:
			p.nextToken()
		}
	}
}

func (p *Parser) noPrefixParsletError(tt token.TokenType) {
	var msg string
	if tt == token.ILLEGAL {
		msg = fmt.Sprintf(ERR_ILLEGAL_TOKEN, p.currToken.Literal)
	} else {
		msg = fmt.Sprintf(ERR_NO_PREFIX_PARSLET_FOUND, tt)
	}
	p.error(msg)
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParslet) {
	p.prefixParslets[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParslet) {
	p.infixParslets[tt] = fn

}
