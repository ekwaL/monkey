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
		super;
		super.;
		object.123;
		a.b.c() = 10;
		let x 5;
		let y = 10
		let z = 10;
		let = 10;
		let 838383;
		if (a > b) true;
		if a > b true;
		if (a > b true;
		if (a > b) { true; } else { 10; };
		fn { 1 };
		fn (i, { true };
		fn (x) { true; };
		fn (1) { true; };
		fn (x) true;
		fn x { 1 }
		fn x (i, { true };
		fn x(x) { true; };
		fn x(1) { true; };
		fn x(x) true
		class {}
		class X < {}
		class X < Y }
		class X < Y {
			if (a > b) true;
		}
		[];
		[)];
		[,];
		class X < Y {
		fn x(x) { !true;
		fn (x) { !true;
		[1,)
		`

	l := lexer.New(source)
	p := parser.New(l)

	p.ParseProgram()

	expect := []string{
		parser.ERR_SUPER_NO_DOT,
		parser.ERR_SUPER_NO_IDENTIFIER_AFTER_DOT,
		parser.ERR_GET_NO_PROP_NAME,
		parser.ERR_WRONG_ASSIGNMENT_TARGET,
		parser.ERR_LET_NO_ASSIGN_AFTER_IDENTIFIER,
		parser.ERR_LET_NO_SEMI_AFTER_LET_STMT,
		parser.ERR_LET_NO_IDENTIFIER_AFTER_LET,
		parser.ERR_LET_NO_IDENTIFIER_AFTER_LET,
		parser.ERR_IF_CONDITION_START_LPAREN,
		parser.ERR_IF_CONDITION_END_RPAREN,
		parser.ERR_FN_PARAMETERS_START_LPAREN,
		fmt.Sprintf(parser.ERR_FN_PARAMETER_SHOULD_BE_IDENTIFIER, "{"),
		parser.ERR_FN_PARAMETERS_END_RPAREN,
		fmt.Sprintf(parser.ERR_FN_PARAMETER_SHOULD_BE_IDENTIFIER, "1"),
		parser.ERR_FN_BODY_START_LBRACE,
		parser.ERR_FN_PARAMETERS_START_LPAREN,
		fmt.Sprintf(parser.ERR_FN_PARAMETER_SHOULD_BE_IDENTIFIER, "{"),
		parser.ERR_FN_PARAMETERS_END_RPAREN,
		fmt.Sprintf(parser.ERR_FN_PARAMETER_SHOULD_BE_IDENTIFIER, "1"),
		parser.ERR_FN_BODY_START_LBRACE,
		parser.ERR_CLASS_NO_CLASSNAME,
		parser.ERR_CLASS_NO_SUPER_NAME,
		parser.ERR_CLASS_BODY_START_LBRACE,
		parser.ERR_CLASS_WRONG_DEFINITION,
		fmt.Sprintf(parser.ERR_NO_PREFIX_PARSLET_FOUND, ")"),
		fmt.Sprintf(parser.ERR_NO_PREFIX_PARSLET_FOUND, ","),
		fmt.Sprintf(parser.ERR_NO_PREFIX_PARSLET_FOUND, ")"),
		parser.ERR_ARR_LITERAL_END_BRACKET,
		parser.ERR_FN_BODY_END_RBRACE,
		parser.ERR_FN_BODY_END_RBRACE,
		parser.ERR_CLASS_WRONG_DEFINITION,
		parser.ERR_CLASS_BODY_END_RBRACE,
	}

	errors := p.Errors()
	if len(errors) != len(expect) {
		// t.Fatalf("Wrong parser error count. Got %d, want %d", len(errors), len(expect))
		t.Errorf("Wrong parser error count. Got %d, want %d", len(errors), len(expect))
	}

	i := 0
	for _, msg := range errors {
		if i >= len(expect) {
			t.Errorf("%d: Wrong parser error message. \nGot %q, \nwant nothing", i, msg)
		} else {
			if want := expect[i]; want != msg {
				t.Errorf("%d: Wrong parser error message. \nGot %q, \nwant %q", i, msg, want)
			}
		}
		i++
	}

	for ; i < len(expect); i++ {
		t.Errorf("Want %q, \ngot nothing", expect[i])
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
