package vm_test

import (
	"monkey/ast"
	"monkey/compiler"
	"monkey/vm"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type vmTestCase struct {
	source string
	want   interface{}
}

func TestVM(t *testing.T) {
	tt := []vmTestCase{
		{source: "1", want: int64(1)},
		{source: "2", want: int64(2)},
		{source: "1 + 2", want: int64(3)},
	}

	for _, tc := range tt {
		runVMTest(t, tc)
	}
}

func runVMTest(t testing.TB, tc vmTestCase) {
	t.Helper()

	program := parse(t, tc.source)

	c := compiler.New()

	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Error while compiling: %q.", err)
	}

	vm := vm.New(c.Bytecode())

	err = vm.Run()
	if err != nil {
		t.Fatalf("VM error: %q.", err)
	}

	top := vm.StackTop()

	testObject(t, top, tc.want)
}

func parse(t testing.TB, source string) *ast.Program {
	t.Helper()
	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	ensureNoParserErrors(t, p)

	return program
}

func ensureNoParserErrors(t testing.TB, p *parser.Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q.", msg)
	}
	t.FailNow()
}

func testObject(t testing.TB, obj object.Object, want interface{}) {
	if obj == nil {
		t.Fatal("Got nil, want object.Object")
	}

	switch obj.Type() {
	case object.INTEGER_OBJ:
		o, ok := obj.(*object.Integer)
		if !ok {
			t.Errorf("Object is not an Integer, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(int64)
		if !ok {
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
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
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
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
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if w != o.Value {
			t.Errorf("Wrong object value. Got %v, want %v.", o.Value, w)
		}
	case object.ARRAY_OBJ:
		a, ok := obj.(*object.Array)
		if !ok {
			t.Errorf("Object is not an Array, got %T. (%+v)", obj, obj)
		}
		w, ok := want.([]interface{})
		if !ok {
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if len(w) != len(a.Elements) {
			t.Fatalf("Wrong array value, got %s, want %v", a.Inspect(), w)
		}

		for i, el := range a.Elements {
			testObject(t, el, w[i])
		}
	case object.HASH_OBJ:
		h, ok := obj.(*object.Hash)
		if !ok {
			t.Errorf("Object is not an Hash, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(map[interface{}]interface{})
		if !ok {
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if len(w) != len(h.Pairs) {
			t.Fatalf("Wrong hash value, got %s, want %v", h.Inspect(), w)
		}

		for _, pair := range h.Pairs {
			var ok bool
			var want interface{}

			switch k := pair.Key.(type) {
			case *object.String:
				want, ok = w[k.Value]
			case *object.Boolean:
				want, ok = w[k.Value]
			case *object.Integer:
				want, ok = w[k.Value]
			default:
				t.Fatalf("Unknown key type.")
			}

			if !ok {
				t.Errorf("Want %v, got %+v.", w, h.Pairs)
				t.Errorf("Didn't expect key '%v' in pair '%v: %v'", pair.Key, pair.Key, pair.Value)
			} else {
				testObject(t, pair.Value, want)
			}
		}
	case object.NULL_OBJ:
		_, ok := obj.(*object.Null)
		if !ok {
			t.Errorf("Object is not an Null, got %T. (%+v)", obj, obj)
		}

		if want != nil {
			t.Errorf("Object is Null, but want %v.", want)
		}
	case object.ERROR_OBJ:
		t.Errorf("Got unexpected error object: %v.", obj)
	default:
		t.Errorf("Unknown object type %q.", obj.Type())
	}
}
