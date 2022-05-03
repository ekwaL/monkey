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

// comment
      // comment
  // comment // comment

let result = add(five, ten);

!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
10 <= 11;
11 >= 10;

"string indeed";
let x = "str";

class A {}
instance.field;
this;
super;
`

func TestNextToken(t *testing.T) {
	input := source

	tt := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{

		// let five = 5 ;
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// let ten = 10;
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// let add = fn(x, y) {
		// 	x + y;
		// };
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

		// let result = add(five, ten);
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

		// !-/*5;
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// 5 < 10 > 5;
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.GREATER, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// if (5 < 10) {
		// 	return true;
		// } else {
		// 	return false;
		// }
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		// 10 == 10;
		{token.INT, "10"},
		{token.EQUAL_EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// 10 != 9;
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		// 10 <= 11;
		{token.INT, "10"},
		{token.LESS_EQUAL, "<="},
		{token.INT, "11"},
		{token.SEMICOLON, ";"},

		// 11 >= 10;
		{token.INT, "11"},
		{token.GREATER_EQUAL, ">="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// "string indeed";
		{token.STRING, "string indeed"},
		{token.SEMICOLON, ";"},

		// let x = "str";
		{token.LET, "let"},
		{token.IDENTIFIER, "x"},
		{token.ASSIGN, "="},
		{token.STRING, "str"},
		{token.SEMICOLON, ";"},

		// class A {}
		{token.CLASS, "class"},
		{token.IDENTIFIER, "A"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		// instance.field;
		{token.IDENTIFIER, "instance"},
		{token.DOT, "."},
		{token.IDENTIFIER, "field"},
		{token.SEMICOLON, ";"},
		// this;
		{token.THIS, "this"},
		{token.SEMICOLON, ";"},
		// super;
		{token.SUPER, "super"},
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
				"tests[%d] - literal wrong. expected %q, got %q",
				i, tc.expectedLiteral, tok.Literal,
			)
		}
	}
}
