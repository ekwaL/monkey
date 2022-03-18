package parser_test

import (
	"monkey/lexer"
	"monkey/parser"
	"testing"
)

func TestParserError(t *testing.T) {
	source := `
		let x 5;
		let y = 10
		let z = 10;
		let = 10;
		let 838383;
		`

	l := lexer.New(source)
	p := parser.New(l)

	p.ParseProgram()

	expect := []string{
		parser.ERR_LET_NO_ASSIGN_AFTER_IDENTIFIER,
		parser.ERR_LET_NO_SEMI_AFTER_LET_STMT,
		parser.ERR_LET_NO_IDENTIFIER_AFTER_LET,
		parser.ERR_LET_NO_IDENTIFIER_AFTER_LET,
	}

	errors := p.Errors()
	if len(errors) != len(expect) {
		t.Fatalf("Wrong parser error count. Got %d, want %d", len(errors), len(expect))
	}

	for i, msg := range errors {
		if want := expect[i]; want != msg {
			t.Errorf("Wrong parser error message. Got %q, want %q", msg, want)
		}
	}
}
