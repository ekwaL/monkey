package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntLiteralExpr(t *testing.T) {
	tt := []struct {
		source string
		want   interface{}
	}{
		{source: "5;", want: int64(5)},
		{source: "228322;", want: int64(228322)},
		{source: "true;", want: true},
		{source: "false;", want: false},
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
		// o, ok := obj.(*object.Null)
		// if !ok {
		// 	t.Errorf("Object is not an Null, got %T. (%+v)", obj, obj)
		// }
		// w, ok := want.(int64)
		// if !ok {
		// 	t.Errorf("Can not compare %q value with %T .", obj.Type(), want)
		// }

		// if w != o.Value {
		// 	t.Errorf("Wrong object value. Got %v, want %v.", o.Value, w)
		// }
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
