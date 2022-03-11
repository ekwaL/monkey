package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '+':
		tok = makeToken(token.PLUS, l.ch)
	case '=':
		tok = makeToken(token.ASSIGN, l.ch)
	case '(':
		tok = makeToken(token.LPAREN, l.ch)
	case ')':
		tok = makeToken(token.RPAREN, l.ch)
	case '{':
		tok = makeToken(token.LBRACE, l.ch)
	case '}':
		tok = makeToken(token.RBRACE, l.ch)
	case ',':
		tok = makeToken(token.COMMA, l.ch)
	case ';':
		tok = makeToken(token.SEMICOLON, l.ch)
	case 0:
		tok = makeToken(token.EOF, 0)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readInt()
			tok.Type = token.INT
			return tok
		} else {
			tok = makeToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readInt() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func makeToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
