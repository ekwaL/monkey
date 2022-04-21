package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type expectFn struct {
	params []string
	body   string
}

func TestEval(t *testing.T) {
	tt := []struct {
		source string
		want   interface{}
	}{
		{source: "5;", want: int64(5)},
		{source: "228322;", want: int64(228322)},
		{source: "true;", want: true},
		{source: "false;", want: false},
		{source: `"string literal";`, want: "string literal"},
		// {source: "null;", want: nil},
		{source: "!true;", want: false},
		{source: "!false;", want: true},
		{source: "!5;", want: false},
		{source: "!!true;", want: true},
		{source: "!!false;", want: false},
		{source: "!!5;", want: true},
		{source: "-5;", want: int64(-5)},
		{source: "--5;", want: int64(5)},
		{source: "--5;", want: int64(5)},
		{source: "5 + 5;", want: int64(10)},
		{source: "5 - 5;", want: int64(0)},
		{source: "5 * 5;", want: int64(25)},
		{source: "5 / 5;", want: int64(1)},
		{source: "5 > 5;", want: false},
		{source: "5 < 5;", want: false},
		// {source: "5 >= 5;", want: true},
		// {source: "5 <= 5;", want: true},
		// {source: "null == null;", want: false},
		{source: "5 == 5;", want: true},
		{source: "5 != 5;", want: false},
		{source: `"string" + " " + "concatenation"`, want: "string concatenation"},
		{source: `"str" == "str";`, want: true},
		{source: `"str1" == "str2";`, want: false},
		{source: "(2 + 2) * 2 == 8;", want: true},
		{source: "2 + 2 * 2 == 6;", want: true},
		{source: "(2 + 2) * 2 > 2 + 2 * 2;", want: true},
		{source: "(5 < 5) == false;", want: true},
		{source: "2 > 3 != 3 > 4;", want: false},
		{source: "if (true) 10;", want: int64(10)},
		{source: "if (1) 10;", want: int64(10)},
		{source: "if (true) { 10; } else { 20; };", want: int64(10)},
		{source: "if (false) { 10; } else { 20; };", want: int64(20)},
		{source: "if (1 > 2) 10;", want: nil},
		{source: "if (2 > 1) 10;", want: int64(10)},
		{source: "return 10;", want: int64(10)},
		{source: "return true; 10;", want: true},
		{source: "if (1 > 2) 10;", want: nil},
		{source: "10; return 2 == 3; 20;", want: false},
		{source: "if (2 > 1) { if (3 > 2) { return 10; }; return 1;};", want: int64(10)},
		{source: "let a = 10;", want: int64(10)},
		// {source: "{ 20; let a = 10; }", want: int64(20)},
		{source: "let a = 10; a;", want: int64(10)},
		{source: "let a = 10 * 5; a;", want: int64(50)},
		{source: "let a = 10; let b = a; b;", want: int64(10)},
		{source: "let a = 10; let b = a; let c = a + b + 20; c;", want: int64(40)},
		{source: "fn(a, b) { a + b }", want: &expectFn{[]string{"a", "b"}, "{ (a + b); }"}},
		{source: "let x = fn(a, b) { a + b }; x;", want: &expectFn{[]string{"a", "b"}, "{ (a + b); }"}},
		{source: "let i = fn(x) { x }; i(10);", want: int64(10)},
		{source: "let i = fn(x) { return x; }; i(10);", want: int64(10)},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			got := eval(t, tc.source)

			if got == nil {
				t.Errorf("Error while evaluating %q, got nil, want '%v'.", tc.source, tc.want)
			}

			testObject(t, got, tc.want)
		})
	}
}

func TestRuntimeErrorHandling(t *testing.T) {
	tt := []struct {
		source string
		want   string
	}{
		{source: "-true;", want: "unknown operator: -BOOLEAN"},
		{source: "true - 5;", want: "type mismatch: BOOLEAN - INTEGER"},
		{source: "5 > true;", want: "type mismatch: INTEGER > BOOLEAN"},
		{source: "5 + true; 5;", want: "type mismatch: INTEGER + BOOLEAN"},
		{source: "true + false;", want: "unknown operator: BOOLEAN + BOOLEAN"},
		{source: "true + false;", want: "unknown operator: BOOLEAN + BOOLEAN"},
		{source: `"str1" < "str2";`, want: "unknown operator: STRING < STRING"},
		{source: `"str1" - "str2";`, want: "unknown operator: STRING - STRING"},
		{source: "5; true - 1; 5;", want: "type mismatch: BOOLEAN - INTEGER"},
		{source: "if (10 > 0) { return true + false; }; 6;", want: "unknown operator: BOOLEAN + BOOLEAN"},
		{source: "(true - false) * 1", want: "unknown operator: BOOLEAN - BOOLEAN"},
		{source: "x;", want: "identifier not found: 'x'"},
		{source: "let a = 10; y;", want: "identifier not found: 'y'"},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			got := eval(t, tc.source)

			err, ok := got.(*object.Error)
			if !ok {
				t.Errorf("No error object returned, got %T (%+v).", got, got)
			}
			if err.Message != tc.want {
				t.Errorf("Wrong error message, got %q, want %q.", err.Message, tc.want)
			}
		})
	}
}

func testObject(t testing.TB, obj object.Object, want interface{}) {
	switch obj.Type() {
	case object.INTEGER_OBJ:
		o, ok := obj.(*object.Integer)
		if !ok {
			t.Errorf("Object is not an Integer, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(int64)
		if !ok {
			t.Errorf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if w != o.Value {
			t.Errorf("Wrong object value. Got %v, want %v.", o.Value, w)
		}
	case object.BOOLEAN_OBJ:
		o, ok := obj.(*object.Boolean)
		if !ok {
			t.Errorf("Object is not an Boolean, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(bool)
		if !ok {
			t.Errorf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if w != o.Value {
			t.Errorf("Wrong object value. Got %v, want %v.", o.Value, w)
		}
	case object.STRING_OBJ:
		o, ok := obj.(*object.String)
		if !ok {
			t.Errorf("Object is not an String, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(string)
		if !ok {
			t.Errorf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if w != o.Value {
			t.Errorf("Wrong object value. Got %v, want %v.", o.Value, w)
		}
	case object.NULL_OBJ:
		_, ok := obj.(*object.Null)
		if !ok {
			t.Errorf("Object is not an Null, got %T. (%+v)", obj, obj)
		}

		if want != nil {
			t.Errorf("Object is Null, but want %v.", want)
		}
	case object.FUNCTION_OBJ:
		fn, ok := obj.(*object.Function)
		if !ok {
			t.Errorf("Object is not an Function, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(*expectFn)
		if !ok {
			t.Errorf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if len(fn.Parameters) != len(w.params) {
			t.Errorf("Wrong parameters number, got %d, want %d.", len(fn.Parameters), len(w.params))
		}

		for i, p := range fn.Parameters {
			if p.Value != w.params[i] {
				t.Errorf("Wrong parameter name, parameter index=%d, got %q, want %q.", i, p.Value, w.params[i])
			}
		}

		if fn.Body.String() != w.body {
			t.Errorf("Wrong function body, got %q, want %q.", fn.Body.String(), w.body)
		}
	default:
		t.Errorf("Unknown object type %q.", obj.Type())
	}
}

func eval(t testing.TB, source string) object.Object {
	t.Helper()

	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("Error while parsing %q.", source)
	}

	return Eval(program, object.NewEnvironment())
}
