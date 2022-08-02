package compiler_test

import (
	"monkey/ast"
	"monkey/code"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	source           string
	wantConstants    []interface{}
	wantInstructions []code.Instructions
}

func TestCompiler(t *testing.T) {
	tt := []compilerTestCase{
		{
			source:        "1 + 2;",
			wantConstants: []interface{}{int64(1), int64(2)},
			wantInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	for _, tc := range tt {
		runCompilerTest(t, tc)
	}
}

func runCompilerTest(t testing.TB, tc compilerTestCase) {
	t.Helper()

	program := parse(t, tc.source)

	c := compiler.New()

	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Error while compiling: %q.", err)
	}

	bytecode := c.Bytecode()

	testInstructions(t, tc.wantInstructions, bytecode.Instructions)
	testConstants(t, tc.wantConstants, bytecode.Constants)
}

func testInstructions(t testing.TB, wantInstructions []code.Instructions, got code.Instructions) {
	t.Helper()

	want := code.Instructions{}
	for _, ins := range wantInstructions {
		want = append(want, ins...)
	}

	if len(want) != len(got) {
		t.Errorf("Wrong instructions length, want %d, got %d.", len(want), len(got))
	}

	i := 0
	for ; i < len(want); i++ {
		if i >= len(got) {
			t.Errorf("Want byte %d at position %d, got nothing.", want[i], i)
			continue
		}

		if want[i] != got[i] {
			t.Errorf("Want byte %d at position %d, got %d.", want[i], i, got[i])
		}
	}

	for ; i < len(got); i++ {
		t.Errorf("Got byte %d at position %d, want nothing.", got[i], i)
	}
}

func testConstants(t testing.TB, want []interface{}, got []object.Object) {
	if len(want) != len(got) {
		t.Errorf("Wrong constants count, want %d, got %d.", len(want), len(got))
	}

	i := 0
	for ; i < len(want); i++ {
		if i >= len(got) {
			t.Errorf("Want constant %v at position %d, got nothing.", want[i], i)
			continue
		}

		testObject(t, got[i], want[i])
	}

	for ; i < len(got); i++ {
		t.Errorf("Got constant %q at position %d, want nothing.", got[i].Inspect(), i)
	}
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
	// case object.FUNCTION_OBJ:
	// 	fn, ok := obj.(*object.Function)
	// 	if !ok {
	// 		t.Errorf("Object is not an Function, got %T. (%+v)", obj, obj)
	// 	}
	// 	w, ok := want.(*expectFn)
	// 	if !ok {
	// 		t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
	// 	}

	// 	if len(fn.Parameters) != len(w.params) {
	// 		t.Errorf("Wrong parameters number, got %d, want %d.", len(fn.Parameters), len(w.params))
	// 	}

	// 	for i, p := range fn.Parameters {
	// 		if p.Value != w.params[i] {
	// 			t.Errorf("Wrong parameter name, parameter index=%d, got %q, want %q.", i, p.Value, w.params[i])
	// 		}
	// 	}

	// 	if fn.Body.String() != w.body {
	// 		t.Errorf("Wrong function body, got %q, want %q.", fn.Body.String(), w.body)
	// 	}
	// case object.CLASS_OBJ:
	// 	class, ok := obj.(*object.Class)
	// 	if !ok {
	// 		t.Errorf("Object is not an Class, got %T. (%+v)", obj, obj)
	// 	}

	// 	w, ok := want.(string)
	// 	if !ok {
	// 		t.Fatalf("Can not compare class name value with %T .", want)
	// 	}

	// 	if class.Name.Value != w {
	// 		t.Errorf("Wrong class name, want %q, got %q.", w, class.Name.Value)
	// 	}
	// case object.INSTANCE_OBJ:
	// 	inst, ok := obj.(*object.Instance)
	// 	if !ok {
	// 		t.Errorf("Object is not an Instance, got %T. (%+v)", obj, obj)
	// 	}

	// 	w, ok := want.(string)
	// 	if !ok {
	// 		t.Fatalf("Can not compare instance class name value with %T .", want)
	// 	}

	// 	if inst.Class.Name.Value != w {
	// 		t.Errorf("Wrong instance class, want %q, got %q.", w, inst.Class.Name.Value)
	// 	}
	case object.ERROR_OBJ:
		t.Errorf("Got unexpected error object: %v.", obj)
	default:
		t.Errorf("Unknown object type %q.", obj.Type())
	}
}
