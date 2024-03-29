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

	l.skipWhitespaceAndComments()

	switch l.ch {
	case '+':
		tok = makeToken(token.PLUS, l.ch)
	case '-':
		tok = makeToken(token.MINUS, l.ch)
	case '*':
		tok = makeToken(token.STAR, l.ch)
	case '/':
		tok = makeToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.BANG, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQUAL_EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.ASSIGN, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LESS_EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.LESS, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GREATER_EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.GREATER, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.ILLEGAL, l.ch)
		}
	case '|':
		ch := l.ch
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '}' {
			l.readChar()
			tok = token.Token{Type: token.RHASHBRACE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.ILLEGAL, l.ch)
		}
	case '(':
		tok = makeToken(token.LPAREN, l.ch)
	case ')':
		tok = makeToken(token.RPAREN, l.ch)
	case '{':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LHASHBRACE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = makeToken(token.LBRACE, l.ch)
		}
	case '}':
		tok = makeToken(token.RBRACE, l.ch)
	case '[':
		tok = makeToken(token.LBRACKET, l.ch)
	case ']':
		tok = makeToken(token.RBRACKET, l.ch)
	case '.':
		tok = makeToken(token.DOT, l.ch)
	case ',':
		tok = makeToken(token.COMMA, l.ch)
	case ':':
		tok = makeToken(token.COLON, l.ch)
	case ';':
		tok = makeToken(token.SEMICOLON, l.ch)
	case 0:
		tok = makeToken(token.EOF, 0)
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
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
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	l.readChar()
	from := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}

	return l.input[from:l.position]
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespaceAndComments() {
	l.skipWhitespace()
	if l.ch == '/' && l.peekChar() == '/' {
		l.skipComment()
		l.skipWhitespaceAndComments()
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
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
