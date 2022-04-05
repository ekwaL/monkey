package parser_test

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
	"strings"
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

func TestIntParseError(t *testing.T) {
	source := `
		123456789101112131415161718192021222324252627282930;
		`

	l := lexer.New(source)
	p := parser.New(l)

	p.ParseProgram()

	errors := p.Errors()
	wantLen := 1
	want := fmt.Sprintf(parser.ERR_COULD_NOT_PARSE_INT, "123456789101112131415161718192021222324252627282930", "")

	if len(errors) != wantLen {
		t.Fatalf("Wrong parser error count. Got %d, want %d", len(errors), wantLen)
	}

	if msg := errors[0]; !strings.HasPrefix(msg, want) {
		t.Errorf("Wrong parser error message. %q should start with %q", msg, want)
	}
}

// No error should be thrown since any misstyped bool literal is an identifier.
// func TestBoolParseError(t *testing.T) {
// 	source := `
// 		truefal`

// 	l := lexer.New(source)
// 	p := parser.New(l)

// 	p.ParseProgram()

// 	errors := p.Errors()
// 	wantLen := 1
// 	want := fmt.Sprintf(parser.ERR_COULD_NOT_PARSE_BOOL, "trufal", "")

// 	if len(errors) != wantLen {
// 		t.Fatalf("Wrong parser error count. Got %d, want %d", len(errors), wantLen)
// 	}

// 	if msg := errors[0]; !strings.HasPrefix(msg, want) {
// 		t.Errorf("Wrong parser error message. %q should start with %q", msg, want)
// 	}
// }
