package eval_test

import (
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/resolver"
	"testing"
)

type expectFn struct {
	params []string
	body   string
}

type expectClass struct {
	name  string
	super string
	props []string
}

type expectInstance struct {
	class string
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
		{source: "10; return;", want: nil},
		{source: "return true; 10;", want: true},
		{source: "if (1 > 2) 10;", want: nil},
		{source: "10; return 2 == 3; 20;", want: false},
		{source: "if (2 > 1) { if (3 > 2) { return 10; }; return 1;};", want: int64(10)},
		{source: "let a = 10;", want: int64(10)},
		{source: "let a;", want: nil},
		{source: "let a; a;", want: nil},
		{source: "let a; a = 10;", want: int64(10)},
		{source: "let a = 10; let b = a = 20;", want: int64(20)},
		// {source: "{ 20; let a = 10; }", want: int64(20)},
		{source: "let a = 10; a;", want: int64(10)},
		{source: "let a = 10 * 5; a;", want: int64(50)},
		{source: "let a = 10; let b = a; b;", want: int64(10)},
		{source: "let a = 10; let b = a; let c = a + b + 20; c;", want: int64(40)},
		{source: "let a = 10; a = 20;", want: int64(20)},
		{source: "let a = 10; let b = 20; a = b = 30;", want: int64(30)},
		{source: "let a = 1; let b = 2; let c = 3; a = b = c;", want: int64(3)},
		{source: "fn(a, b) { a + b }", want: &expectFn{[]string{"a", "b"}, "{ (a + b); }"}},
		{source: "let x = fn(a, b) { a + b }; x;", want: &expectFn{[]string{"a", "b"}, "{ (a + b); }"}},
		{source: "let i = fn(x) { x }; i(10);", want: int64(10)},
		{source: "let i = fn(x) { return x; }; i(10);", want: int64(10)},
		{source: "fn f(x) { x }; f(10);", want: int64(10)},
		{source: "fn f(x) { return x } f(10);", want: int64(10)},
		{source: `len("Hello, World!")`, want: int64(13)},
		{source: `len("")`, want: int64(0)},
		{source: `len("Hello, World!"); { let len = 10; len; }`, want: int64(10)},
		{source: `len("Hello, World!"); { let len = 10; len; } len("Hello, World!")`, want: int64(13)},
		{source: "let x = 10; let y = 10; { let x = x; x = 20; y = x; } x;", want: int64(10)},
		{source: "let x = 10; let y = 10; { let x = x; x = 20; y = x; } y;", want: int64(20)},
		{source: "let x = 10; let f = fn() { x }; { f() }", want: int64(10)},
		{source: "let x = 10; { let f = fn() { x }; let x = 20; f(); }", want: int64(10)},
		{source: "let x = 10; { let f = fn() { x }; let x = 20; f(); } x;", want: int64(10)},
		{source: "let x = 10; { fn f() { x }; let x = 20; f(); } x;", want: int64(10)},
		// classes
		{source: "class A {}", want: &expectClass{name: "A", super: "", props: []string{}}},

		{source: "class B {} class A < B {}", want: &expectClass{name: "A", super: "B", props: []string{}}},
		{source: "class A {} A();", want: &expectInstance{class: "A"}},
		{source: `class A {} A().field = 10;`, want: int64(10)},
		{source: "class A {} let obj = A(); obj.field = 10; obj.field;", want: int64(10)},
		{source: "class A { let method = fn() { 10; }; } let obj = A(); obj.method()", want: int64(10)},
		{source: "class A { fn method() { 10; }} let obj = A(); obj.method()", want: int64(10)},
		{source: "class A { fn method() { this.x; }} let obj = A(); obj.x = 1; obj.method()", want: int64(1)},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					let obj = A(); obj.method()`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					let obj = A(); obj.x = 50; obj.init()`,
			want: &expectInstance{class: "A"},
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					let obj = A(); obj.x = 50; obj.init(); obj.x`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					class B < A{}
					let obj = B();`,
			want: &expectInstance{class: "B"},
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					class B < A{}
					let obj = B(); obj.x;`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					class B < A {
						fn init() { super.init() }
					}
					let obj = B(); obj.x;`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x * 2; }
					}
					class B < A {
						fn init() { super.init() }
						fn method() { this.X }
						fn doubleX() { super.method() }
					}
					class C < B {}
					let obj = C(); obj.doubleX();`,
			want: int64(40),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x * 2; }
					}
					class B < A {
						fn init() { super.init() }
						fn method() { this.X }
						fn doubleX() { super.method }
					}
					class C < B {}
					let obj = C(); let d = obj.doubleX(); obj.x = 10; d()`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					class B < A {}
					let obj = B(); obj.x = 50; obj.init(); obj.method()`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					class B {
						fn init() {this.x = 10; }
					}
					let a = A(); let b = B(); b.method = a.method; b.method();`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init(n) { this.x = n; }
						fn method() { this.x; }
					}
					class B < A{}
					let obj = B(10); obj.x`,
			want: int64(10),
		},
		{
			source: `class A {
						fn init(n) { this.x = n; }
						fn method() {
							class B < A {}
							B(20);
						}
					}
					let obj = A(10).method(); obj.x`,
			want: int64(20),
		},
		{
			source: `class A {
						fn init(n) { this.x = n; }
						fn new(x) {
							A(x);
						}
					}
					let obj = A(10).new(20); obj.x`,
			want: int64(20),
		},
		{
			source: "let x = 1; class B {} { class A < B { fn f() { x = 20; } } A().f() } x",
			want:   int64(20),
		},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			got := evalSource(t, tc.source)

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
		{source: "let x = fn() { true }; x(1)", want: "wrong arguments count: expect 0, got 1"},
		{source: "len(10)", want: "type mismatch: len(INTEGER)"},
		{source: `len("one", "two")`, want: "wrong arguments count: expect 1, got 2"},
		{source: "if (10 > 0) { return true + false; }; 6;", want: "unknown operator: BOOLEAN + BOOLEAN"},
		{source: "(true - false) * 1", want: "unknown operator: BOOLEAN - BOOLEAN"},
		{source: "x;", want: "identifier not found: 'x'"},
		{source: "let a = 10; return y;", want: "identifier not found: 'y'"},
		{source: "let a = 10; b = 20;", want: "identifier not found: 'b'"},
		{source: "let a = 10; a = b;", want: "identifier not found: 'b'"},
		{source: "let function = 1; function(false)", want: "not a function: INTEGER '1'"},
		{source: "let x = 10; { let f = x; } f;", want: "identifier not found: 'f'"},
		{source: `len("Hello, World!"); { let len = 10; len; len("Hello, World!")}`, want: "not a function: INTEGER '10'"},
		{source: `class A {} A(1, "b");`, want: "wrong arguments count: expect 0, got 2"},
		{source: `class A {} A.field;`, want: "only instances have properties: CLASS.field"},
		{source: `class A {} A.field = 10;`, want: "only instances have fields: CLASS.field"},
		{source: `class A {} A().field;`, want: "undefined property: 'field'"},
		{source: "class A { fn method() { this.x; }} let obj = A(); obj.method()", want: "undefined property: 'x'"},
		{source: `"hi".field;`, want: "only instances have properties: STRING.field"},
		{source: `"hi".field = 10;`, want: "only instances have fields: STRING.field"},
		{source: "class A < B {}", want: "identifier not found: 'B'"},
		{source: "let B = 10; class A < B {}", want: "superclass must be a class: 'A < INTEGER'"},
		{
			source: `class A {
						fn init() { this.x = 20; }
						fn method() { this.x; }
					}
					class B < A{
						fn init() {}
					}
					let obj = B(); obj.x;`,
			want: "undefined property: 'x'",
		},
		{
			source: `class A {
						fn init(n) { this.x = n; }
						fn method() { this.x; }
					}
					class B < A{}
					let obj = B(); obj.x`,
			want: "wrong arguments count: expect 1, got 0",
		},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			got := evalSource(t, tc.source)

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
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
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
	case object.CLASS_OBJ:
		class, ok := obj.(*object.Class)
		if !ok {
			t.Errorf("Object is not an Class, got %T. (%+v)", obj, obj)
		}
		w, ok := want.(*expectClass)
		if !ok {
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if class.Name.Value != w.name {
			t.Errorf("Wrong class name, want %q, got %q.", w.name, class.Name.Value)
		}

		gotSuper := ""
		if class.Super != nil {
			gotSuper = class.Super.Name.Value
		}
		if gotSuper != w.super {
			t.Errorf("Wrong superclass name, want %q, got %q.", w.super, gotSuper)
		}

		if len(class.Methods) != len(w.props) {
			t.Errorf("Wrong properties number, got %d, want %d.", len(class.Methods), len(w.props))
		}

		methods := []string{}
		for m := range class.Methods {
			methods = append(methods, m)
		}
		for i, m := range methods {
			if m != w.props[i] {
				t.Errorf("Wrong property name, prop index=%d, got %q, want %q.", i, m, w.props[i])
			}
		}

		// if class.Body.String() != w.body {
		// 	t.Errorf("Wrong function body, got %q, want %q.", class.Body.String(), w.body)
		// }
	case object.INSTANCE_OBJ:
		inst, ok := obj.(*object.Instance)
		if !ok {
			t.Errorf("Object is not an Instance, got %T. (%+v)", obj, obj)
		}

		w, ok := want.(*expectInstance)
		if !ok {
			t.Fatalf("Can not compare %q value with %T .", obj.Type(), want)
		}

		if inst.Class.Name.Value != w.class {
			t.Errorf("Wrong instance class, want %q, got %q.", w.class, inst.Class.Name.Value)
		}

	case object.ERROR_OBJ:
		t.Errorf("Got unexpected error object: %v.", obj)
	default:
		t.Errorf("Unknown object type %q.", obj.Type())
	}
}

func evalSource(t testing.TB, source string) object.Object {
	t.Helper()

	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("Error while parsing %q.", source)
	}

	r := resolver.New()
	r.Resolve(program)

	eval.Locals = r.Locals()

	return eval.Eval(program, object.NewEnvironment())
}
