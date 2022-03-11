package lexer_test

import (
	"monkey/lexer"
	"monkey/token"
	"testing"
)

const source = `
let five = 5 ;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`

func TestNextToken(t *testing.T) {
	input := source

	tt := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.EOF, "\x00"},
	}

	l := lexer.New(input)

	for i, tc := range tt {
		tok := l.NextToken()

		if tok.Type != tc.expectedType {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected %q, got %q",
				i, tc.expectedType, tok.Type,
			)
		}

		if tok.Literal != tc.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expedted %q, got %q",
				i, tc.expectedLiteral, tok.Literal,
			)
		}
	}
}
