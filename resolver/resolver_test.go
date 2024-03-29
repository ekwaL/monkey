package resolver_test

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"monkey/resolver"
	"reflect"
	"testing"
)

func TestResolver(t *testing.T) {
	tt := []struct {
		source string
		want   map[string]int
	}{
		{source: "x;", want: map[string]int{}},
		{source: "let x; x;", want: map[string]int{"x": 0}},
		{source: "let x = 10; x;", want: map[string]int{"x": 0}},
		{source: "let x = 10; { let f = x; }", want: map[string]int{"x": 1}},
		{source: "let x = 10; { let x = x; }", want: map[string]int{"x": 1}},
		{source: "let x = 10; x = x;", want: map[string]int{"x": 0}},
		{source: "let x = 10; x = x + 1;", want: map[string]int{"x": 0}},
		{source: "let x = 10; { x; }", want: map[string]int{"x": 1}},
		{source: "let x = 10; { [x, 20, 30]; }", want: map[string]int{"x": 1}},
		{source: "let x = 1; { [10, 20, 30][x]; }", want: map[string]int{"x": 1}},
		{source: "let x = 1; { {| 10: x, 20: 30 |}; }", want: map[string]int{"x": 1}},
		{source: "let x = 10; { {| 10: 20, 30: 40 |}[x]; }", want: map[string]int{"x": 1}},
		{source: "let x = 10; x = 20;", want: map[string]int{"x": 0}},
		{source: "{ let x = true; return !x; }", want: map[string]int{"x": 0}},
		{source: "let a = 1; let b = 2; return a + b;", want: map[string]int{"a": 0, "b": 0}},
		{
			source: "let x = true; let a = 1; let b = 2; if (x) a else b;",
			want:   map[string]int{"x": 0, "a": 0, "b": 0},
		},
		{
			source: "let f = fn(x) { x }; f(10);",
			want:   map[string]int{"x": 0, "f": 0},
		},
		{
			source: "let x = 10; let f = fn() { x }; { f() }",
			want:   map[string]int{"f": 1, "x": 1},
		},
		{
			source: "class A < B {}",
			want:   map[string]int{},
		},
		{
			source: "class B {} class A < B {}",
			want:   map[string]int{"B": 0},
		},
		{
			source: `class A {
					fn init(n) { this.x = n; }
					fn method() {
						class B < A { fn m() { super.init(10) } }
						B();
					}
				  }`,
			want: map[string]int{"A": 2, "this": 1, "B": 0, "n": 0, "super": 2},
		},
		{
			source: "let x = 1; class B {} { class A < B { fn f() { x = 20; }} }",
			want:   map[string]int{"B": 1, "x": 4},
		},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			locals, errors := resolve(t, tc.source)

			if len(errors) != 0 {
				for _, err := range errors {
					t.Error(err)
				}
				t.Fatal("Unexpected errors while resolving.")
			}

			got := make(map[string]int)
			for name, depth := range locals {
				switch local := name.(type) {
				case *ast.IdentifierExpr:
					got[local.Value] = depth
				case *ast.ThisExpr:
					got[local.TokenLiteral()] = depth
				case *ast.SuperExpr:
					got[local.TokenLiteral()] = depth
				default:
					t.Fatalf("Expect local to be an Identifier, This or Super, got %T.", name)
				}
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResolverErrorHandling(t *testing.T) {
	tt := []struct {
		source string
		want   []string
	}{
		{source: "let x = x;", want: []string{"Can not read variable 'x' in it's own initializer."}},
		{
			source: "let x = x; { let y = y; }",
			want: []string{
				"Can not read variable 'x' in it's own initializer.",
				"Can not read variable 'y' in it's own initializer.",
			},
		},
		{source: "let x = 10; let x = x;", want: []string{"Variable 'x' is already declared in current scope."}},
		{source: "let x = 10; let x = y;", want: []string{"Variable 'x' is already declared in current scope."}},
		{source: "let x = 10; { let x = y; }", want: []string{}},
		// {
		// 	source: "return;",
		// 	want:   []string{resolver.ERR_TOP_LEVEL_RETURN},
		// },
		{
			source: `
					class X {
						fn init() {
							return 10;
						}
					}
			`,
			want: []string{resolver.ERR_INITIALIZER_VAL_RETURN},
		},
		{
			source: "class A < A {}",
			want:   []string{"Class 'A' can not inherit from itself."},
		},
		{
			source: "this;",
			want:   []string{resolver.ERR_THIS_OUTSIDE_OF_CLASS},
		},
		{
			source: "super.method;",
			want:   []string{resolver.ERR_SUPER_OUTSIDE_OF_CLASS},
		},
		{
			source: "class A { fn f() { super.method; } }",
			want:   []string{resolver.ERR_SUPER_WITHOUT_SUPERCLASS},
		},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			_, errors := resolve(t, tc.source)

			i := 0
			for ; i < len(tc.want); i++ {
				if i >= len(errors) {
					t.Errorf("Want error %q, got none.", tc.want[i])
				}

				if tc.want[i] != errors[i] {
					t.Errorf("Wrong error message. Want %q, got %q.", tc.want[i], errors[i])
				}
			}

			for ; i < len(errors); i++ {
				t.Errorf("Got error %q, expected none.", errors[i])
			}
		})
	}
}

func resolve(t testing.TB, source string) (map[ast.Expression]int, []string) {
	t.Helper()

	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, e := range p.Errors() {
			t.Error(e)
		}
		t.Fatalf("Error while parsing %q.", source)
	}

	r := resolver.New()
	r.Resolve(program)

	return r.Locals(), r.Errors()
}
