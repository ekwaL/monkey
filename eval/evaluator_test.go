package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEval(t *testing.T) {
	tt := []struct {
		source string
		want   interface{}
	}{
		{source: "5;", want: int64(5)},
		{source: "228322;", want: int64(228322)},
		{source: "true;", want: true},
		{source: "false;", want: false},
		// {source: "null;", want: nil},
		{source: "!true;", want: false},
		{source: "!false;", want: true},
		{source: "!5;", want: false},
		{source: "!!true;", want: true},
		{source: "!!false;", want: false},
		{source: "!!5;", want: true},
		{source: "-5;", want: int64(-5)},
		{source: "--5;", want: int64(5)},
		{source: "-true;", want: nil},
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
		{source: "5 > true;", want: nil},
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
	case object.NULL_OBJ:
		_, ok := obj.(*object.Null)
		if !ok {
			t.Errorf("Object is not an Null, got %T. (%+v)", obj, obj)
		}

		if want != nil {
			t.Errorf("Object is Null, but want %v.", want)
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

	return Eval(program)
}
