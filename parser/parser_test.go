package parser_test

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"strconv"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	source := `
		let x = 5;
		let y = 10;
		let foobarbaz = 10;
		let x;
		`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 4
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	tt := []struct {
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{expectedIdentifier: "x", expectedValue: 5},
		{expectedIdentifier: "y", expectedValue: 10},
		{expectedIdentifier: "foobarbaz", expectedValue: 10},
		{expectedIdentifier: "x", expectedValue: nil},
	}

	for i, tc := range tt {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tc.expectedIdentifier) {
			return
		}
		if stmt.(*ast.LetStmt).Value == nil {
			if tc.expectedValue != nil {
				t.Errorf("Want Let.Value to be nil, got %v.", tc.expectedValue)
			}
		} else {
			testIdentifierOrLiteralExpr(t, stmt.(*ast.LetStmt).Value, tc.expectedValue)
		}
	}

}

func TestReturnStatement(t *testing.T) {
	source := `
		return;
		return 5;
		`
	program := parse(t, source)
	expect := []interface{}{nil, int64(5)}

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := len(expect)
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	for i, expectedValue := range expect {
		stmt := program.Statements[i]
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Errorf("stmt is not *ast.ReturnStmt. Got %T.", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Wrong TokenLiteral, want 'return', got %q.", returnStmt.TokenLiteral())
		}

		if expectedValue != nil {
			testIdentifierOrLiteralExpr(t, returnStmt.Value, expectedValue)
		} else if returnStmt.Value != nil {
			t.Errorf("Wrong Value, want nil, got %v.", returnStmt.Value)
		}
	}
}

func TestAssignExpression(t *testing.T) {
	source := `
		x = 5;`

	program := parse(t, source)
	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	expr := expressionStmt.Expression
	assignExpr, ok := expr.(*ast.AssignExpr)
	if !ok {
		t.Fatalf("expr is not *ast.AssignExpr. Got %T.", expr)
	}

	if assignExpr.TokenLiteral() != "=" {
		t.Errorf("Wrong TokenLiteral, want '=', got %q.", assignExpr.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, assignExpr.Identifier, "x")
	testIdentifierOrLiteralExpr(t, assignExpr.Expression, 5)
}

func TestIfExpression(t *testing.T) {
	source := `
		if (5 > 10) {
			20;
		} else null;`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if exprStmt.TokenLiteral() != "if" {
		t.Errorf("Wrong TokenLiteral, want 'if', got %q.", exprStmt.TokenLiteral())
	}

	expr := exprStmt.Expression
	ifExpr, ok := expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("expr is not *ast.IfExpr. Got %T.", expr)
	}
	if ifExpr.TokenLiteral() != "if" {
		t.Errorf("Wrong TokenLiteral, want 'if', got %q.", ifExpr.TokenLiteral())
	}
	testInfixExpr(t, ifExpr.Condition, 5, ">", 10)

	thenBranch, ok := ifExpr.Then.(*ast.BlockStmt)
	if !ok {
		t.Fatalf("Then-branch is not a BlockStmt, got %T", ifExpr.Then)
	}
	if len(thenBranch.Statements) != 1 {
		t.Errorf("Then-branch block should have 1 statement, got %d",
			len(thenBranch.Statements))
	}

	cons, ok := thenBranch.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("First then-branch statement is not a ExpressionStmt, got %T",
			thenBranch.Statements[0])
	}
	testIdentifierOrLiteralExpr(t, cons.Expression, 20)

	elseBranch, ok := ifExpr.Else.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Else-branch is not *ast.ExpressionStmt. Got %T.", ifExpr.Else)
	}
	if elseBranch.TokenLiteral() != "null" {
		t.Errorf("Wrong TokenLiteral, want 'null', got %q.", elseBranch.TokenLiteral())
	}
	testIdentifierOrLiteralExpr(t, elseBranch.Expression, nil)
}

func TestFunctionExpression(t *testing.T) {
	source := `
		fn(i, j) { i - j; };`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if exprStmt.TokenLiteral() != "fn" {
		t.Errorf("Wrong TokenLiteral, want 'fn', got %q.", exprStmt.TokenLiteral())
	}

	expr := exprStmt.Expression
	fnExpr, ok := expr.(*ast.FunctionExpr)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionExpr. Got %T.", expr)
	}
	if fnExpr.TokenLiteral() != "fn" {
		t.Errorf("Wrong TokenLiteral, want 'fn', got %q.", fnExpr.TokenLiteral())
	}

	if len(fnExpr.Parameters) != 2 {
		t.Errorf("Wrong parameters count: want 2, got %d", len(fnExpr.Parameters))
	}
	testIdentifierOrLiteralExpr(t, fnExpr.Parameters[0], "i")
	testIdentifierOrLiteralExpr(t, fnExpr.Parameters[1], "j")

	if len(fnExpr.Body.Statements) != 1 {
		t.Errorf("Wrong body statements count: want 1, got %d", len(fnExpr.Parameters))
	}

	cons, ok := fnExpr.Body.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("First body statement is not an ExpressionStmt, got %T",
			fnExpr.Body.Statements[0])
	}
	testInfixExpr(t, cons.Expression, "i", "-", "j")
}

func TestFunctionDefinition(t *testing.T) {
	source := `
		fn minus(i, j) { i - j; }`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	testLetStatement(t, stmt, "minus")
	letStmt, ok := stmt.(*ast.LetStmt)
	expr := letStmt.Value
	fnExpr, ok := expr.(*ast.FunctionExpr)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionExpr. Got %T.", expr)
	}
	if fnExpr.TokenLiteral() != "fn" {
		t.Errorf("Wrong TokenLiteral, want 'fn', got %q.", fnExpr.TokenLiteral())
	}

	if len(fnExpr.Parameters) != 2 {
		t.Errorf("Wrong parameters count: want 2, got %d", len(fnExpr.Parameters))
	}
	testIdentifierOrLiteralExpr(t, fnExpr.Parameters[0], "i")
	testIdentifierOrLiteralExpr(t, fnExpr.Parameters[1], "j")

	if len(fnExpr.Body.Statements) != 1 {
		t.Errorf("Wrong body statements count: want 1, got %d", len(fnExpr.Parameters))
	}

	cons, ok := fnExpr.Body.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("First body statement is not an ExpressionStmt, got %T",
			fnExpr.Body.Statements[0])
	}
	testInfixExpr(t, cons.Expression, "i", "-", "j")
}

func TestClassDefinition(t *testing.T) {
	source := `
		class Hello < World {
			fn minus(i, j) { this; super.method; i - j; }
		}`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	classStmt, ok := stmt.(*ast.ClassStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ClassStmt. Got %T.", stmt)
	}
	testIdentifierExpression(t, classStmt.Name, "Hello")
	testIdentifierExpression(t, classStmt.Superclass, "World")

	wantMethodsLen := 1
	if len(classStmt.Methods) != wantMethodsLen {
		t.Fatalf("class.Props len is %d, want %d .", len(classStmt.Methods), wantMethodsLen)
	}

	method := classStmt.Methods[0]
	expr := method.Value
	fnExpr, ok := expr.(*ast.FunctionExpr)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionExpr. Got %T.", expr)
	}
	if fnExpr.TokenLiteral() != "fn" {
		t.Errorf("Wrong TokenLiteral, want 'fn', got %q.", fnExpr.TokenLiteral())
	}

	if len(fnExpr.Parameters) != 2 {
		t.Errorf("Wrong parameters count: want 2, got %d", len(fnExpr.Parameters))
	}
	testIdentifierOrLiteralExpr(t, fnExpr.Parameters[0], "i")
	testIdentifierOrLiteralExpr(t, fnExpr.Parameters[1], "j")

	if len(fnExpr.Body.Statements) != 3 {
		t.Errorf("Wrong body statements count: want 3, got %d", len(fnExpr.Parameters))
	}

	exprStmt, ok := fnExpr.Body.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("First body statement is not an ExpressionStmt, got %T",
			fnExpr.Body.Statements[0])
	}
	this, ok := exprStmt.Expression.(*ast.ThisExpr)
	if !ok {
		t.Fatalf("First body expression is not an ThisExpr, got %T", exprStmt.Expression)
	}
	if this.TokenLiteral() != "this" {
		t.Errorf("Wrong TokenLiteral, want 'this', got %q.", this.TokenLiteral())
	}

	exprStmt, ok = fnExpr.Body.Statements[1].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("Second body statement is not an ExpressionStmt, got %T",
			fnExpr.Body.Statements[1])
	}
	super, ok := exprStmt.Expression.(*ast.SuperExpr)
	if !ok {
		t.Fatalf("Second body expression is not an SuperExpr, got %T", exprStmt.Expression)
	}
	if super.TokenLiteral() != "super" {
		t.Errorf("Wrong TokenLiteral, want 'super', got %q.", super.TokenLiteral())
	}
	testIdentifierOrLiteralExpr(t, super.Method, "method")

	cons, ok := fnExpr.Body.Statements[2].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("First body statement is not an ExpressionStmt, got %T",
			fnExpr.Body.Statements[2])
	}
	testInfixExpr(t, cons.Expression, "i", "-", "j")
}

func TestCallExpression(t *testing.T) {
	source := `
		fun(1, true == false);`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if exprStmt.TokenLiteral() != "fun" {
		t.Errorf("Wrong TokenLiteral, want 'fun', got %q.", exprStmt.TokenLiteral())
	}

	expr := exprStmt.Expression
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.CallExpr. Got %T.", expr)
	}
	if callExpr.TokenLiteral() != "(" {
		t.Errorf("Wrong TokenLiteral, want '(', got %q.", callExpr.TokenLiteral())
	}

	if len(callExpr.Arguments) != 2 {
		t.Errorf("Wrong parameters count: want 1, got %d", len(callExpr.Arguments))
	}
	testIdentifierOrLiteralExpr(t, callExpr.Arguments[0], 1)
	testInfixExpr(t, callExpr.Arguments[1], true, "==", false)
	testIdentifierOrLiteralExpr(t, callExpr.Function, "fun")
}

func TestGetExpression(t *testing.T) {
	source := `
		obj.field;`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if exprStmt.TokenLiteral() != "obj" {
		t.Errorf("Wrong TokenLiteral, want 'obj', got %q.", exprStmt.TokenLiteral())
	}

	expr := exprStmt.Expression
	getExpr, ok := expr.(*ast.GetExpr)
	if !ok {
		t.Fatalf("expr is not *ast.GetExpr. Got %T.", expr)
	}
	if getExpr.TokenLiteral() != "." {
		t.Errorf("Wrong TokenLiteral, want '.', got %q.", getExpr.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, getExpr.Expression, "obj")
	testIdentifierOrLiteralExpr(t, getExpr.Field, "field")
}

func TestSetExpression(t *testing.T) {
	source := `
		obj.field = 10;`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 1
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	exprStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if exprStmt.TokenLiteral() != "obj" {
		t.Errorf("Wrong TokenLiteral, want 'obj', got %q.", exprStmt.TokenLiteral())
	}

	expr := exprStmt.Expression
	setExpr, ok := expr.(*ast.SetExpr)
	if !ok {
		t.Fatalf("expr is not *ast.SetExpr. Got %T.", expr)
	}
	if setExpr.TokenLiteral() != "=" {
		t.Errorf("Wrong TokenLiteral, want '=', got %q.", setExpr.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, setExpr.Expression, "obj")
	testIdentifierOrLiteralExpr(t, setExpr.Field, "field")
	testIdentifierOrLiteralExpr(t, setExpr.Value, 10)
}

func TestIdentifierExpressionStatement(t *testing.T) {
	source := "foobar;"

	program := parse(t, source)

	want := "foobar"
	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if expressionStmt.TokenLiteral() != want {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", want, expressionStmt.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, expressionStmt.Expression, want)
}

func TestIntLiteral(t *testing.T) {
	source := "10;"

	program := parse(t, source)

	want := "10"
	var wantInt int64 = 10
	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if expressionStmt.TokenLiteral() != want {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", want, expressionStmt.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, expressionStmt.Expression, wantInt)
}

func TestBoolLiteral(t *testing.T) {
	source := `
		true;
		false;`

	program := parse(t, source)

	tt := []struct {
		wantLiteral string
		wantValue   bool
	}{
		{wantLiteral: "true", wantValue: true},
		{wantLiteral: "false", wantValue: false},
	}
	wantLen := len(tt)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	for i, tc := range tt {
		stmt := program.Statements[i]

		expressionStmt, ok := stmt.(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
		}
		if expressionStmt.TokenLiteral() != tc.wantLiteral {
			t.Errorf("Wrong TokenLiteral, want %q, got %q.", tc.wantLiteral, expressionStmt.TokenLiteral())
		}

		testIdentifierOrLiteralExpr(t, stmt.(*ast.ExpressionStmt).Expression, tc.wantValue)
	}
}

func TestStringLiteral(t *testing.T) {
	source := `"string literal";`

	program := parse(t, source)

	want := "string literal"
	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	if expressionStmt.TokenLiteral() != want {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", want, expressionStmt.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, expressionStmt.Expression, want)
}

func TestArrayLiteral(t *testing.T) {
	source := `[1, "two", x, 2 * 3, fn () {}];`

	program := parse(t, source)

	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	wantLiteral := "["
	if expressionStmt.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, expressionStmt.TokenLiteral())
	}

	arrExpr, ok := expressionStmt.Expression.(*ast.ArrayLiteralExpr)
	if !ok {
		t.Fatalf("Expression is not *ast.ArrayLiteralExpr. Got %T.", expressionStmt.Expression)
	}
	if arrExpr.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, arrExpr.TokenLiteral())
	}

	want := []interface{}{1, "two", "x"}
	if len(arrExpr.Elements) != len(want)+2 {
		t.Errorf("Wrong Elements length, got %d, want %d.", len(arrExpr.Elements), len(want)+2)
	}
	for i, w := range want {
		testIdentifierOrLiteralExpr(t, arrExpr.Elements[i], w)
	}
	testInfixExpr(t, arrExpr.Elements[len(want)], 2, "*", 3)
	fn := arrExpr.Elements[len(want)+1]
	if _, ok := fn.(*ast.FunctionExpr); !ok {
		t.Errorf("Expect last array Element to be a function expression, got %T.", fn)
	}
}

func TestHashLiteral(t *testing.T) {
	source := `{| "one": 1, "two": 1 + 1, 3: "three" |};`

	program := parse(t, source)

	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	wantLiteral := "{|"
	if expressionStmt.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, expressionStmt.TokenLiteral())
	}

	hashExpr, ok := expressionStmt.Expression.(*ast.HashLiteralExpr)
	if !ok {
		t.Fatalf("Expression is not *ast.HashLiteralExpr. Got %T.", expressionStmt.Expression)
	}
	if hashExpr.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, hashExpr.TokenLiteral())
	}

	want := map[interface{}]interface{}{
		"one":    1,
		"two":    []interface{}{1, "+", 1},
		int64(3): "three",
	}
	if len(hashExpr.Pairs) != len(want) {
		t.Errorf("Wrong Pairs length, got %d, want %d.", len(hashExpr.Pairs), len(want))
	}
	for key, val := range hashExpr.Pairs {
		switch k := key.(type) {
		case *ast.StringLiteralExpr:
			wantVal, ok := want[k.Value]
			if !ok {
				t.Errorf("Wrong key, got %q.", k.Value)
			}
			if infix, ok := wantVal.([]interface{}); ok {
				testInfixExpr(t, val, infix[0], infix[1].(string), infix[2])
			} else {
				testIdentifierOrLiteralExpr(t, val, wantVal)
			}
		case *ast.IntLiteralExpr:
			wantVal, ok := want[k.Value]
			if !ok {
				t.Errorf("Wrong key, got %d.", k.Value)
			}
			if infix, ok := wantVal.([]interface{}); ok {
				testInfixExpr(t, val, infix[0], infix[1].(string), infix[2])
			} else {
				testIdentifierOrLiteralExpr(t, val, wantVal)
			}
		default:
			t.Errorf("Unknown key type: %T", key)
		}
	}
}

func TestEmptyHashLiteral(t *testing.T) {
	source := "{||}"

	program := parse(t, source)

	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	wantLiteral := "{|"
	if expressionStmt.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, expressionStmt.TokenLiteral())
	}

	hashExpr, ok := expressionStmt.Expression.(*ast.HashLiteralExpr)
	if !ok {
		t.Fatalf("Expression is not *ast.HashLiteralExpr. Got %T.", expressionStmt.Expression)
	}
	if hashExpr.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, hashExpr.TokenLiteral())
	}
	if len(hashExpr.Pairs) != 0 {
		t.Errorf("Wrong Pairs length, got %d, want %d.", len(hashExpr.Pairs), 0)
	}

}


func TestIndexExpression(t *testing.T) {
	source := "arr[1 + 1];"

	program := parse(t, source)

	wantLen := 1

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	stmt := program.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
	}
	wantLiteral := "arr"
	if expressionStmt.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, expressionStmt.TokenLiteral())
	}

	wantLiteral = "["
	idxExpr, ok := expressionStmt.Expression.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("Expression is not *ast.IndexExpr. Got %T.", expressionStmt.Expression)
	}
	if idxExpr.TokenLiteral() != wantLiteral {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", wantLiteral, idxExpr.TokenLiteral())
	}

	testIdentifierOrLiteralExpr(t, idxExpr.Left, "arr")
	testInfixExpr(t, idxExpr.Index, 1, "+", 1)
}

func TestParseingPrefixExpressions(t *testing.T) {
	tt := []struct {
		source   string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-20", "-", 20},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tc := range tt {
		program := parse(t, tc.source)

		if program == nil {
			t.Fatal("ParseProgram() returned 'nil'.")
		}

		wantLen := 1
		if len(program.Statements) != wantLen {
			t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
		}

		stmt := program.Statements[0]
		expressionStmt, ok := stmt.(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
		}
		if expressionStmt.TokenLiteral() != tc.operator {
			t.Errorf("Wrong TokenLiteral, want %q, got %q.", tc.operator, expressionStmt.TokenLiteral())
		}

		expr := expressionStmt.Expression
		prefixExpr, ok := expr.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("expr is not *ast.PrefixExpr. Got %T.", expr)
		}
		if prefixExpr.TokenLiteral() != tc.operator {
			t.Errorf("Wrong TokenLiteral, want %q, got %q.", tc.operator, expressionStmt.TokenLiteral())
		}
		if prefixExpr.Operator != tc.operator {
			t.Errorf("Wrong Operator, want %q, got %q.", tc.operator, prefixExpr.Operator)
		}

		testIdentifierOrLiteralExpr(t, prefixExpr.Right, tc.value)
	}
}

func TestParseingInfixExpressions(t *testing.T) {
	tt := []struct {
		source   string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 6;", 5, "+", 6},
		{"5 - 6;", 5, "-", 6},
		{"5 * 6;", 5, "*", 6},
		{"5 / 6;", 5, "/", 6},
		{"5 > 6;", 5, ">", 6},
		{"5 < 6;", 5, "<", 6},
		{"5 == 6;", 5, "==", 6},
		{"5 != 6;", 5, "!=", 6},
		{"true == true;", true, "==", true},
		{"false == false;", false, "==", false},
		{"false != true;", false, "!=", true},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			program := parse(t, tc.source)

			if program == nil {
				t.Fatal("ParseProgram() returned 'nil'.")
			}

			wantLen := 1
			if len(program.Statements) != wantLen {
				t.Fatalf("program.Statements len is %d, want %d .",
					len(program.Statements), wantLen)
			}

			stmt := program.Statements[0]
			expressionStmt, ok := stmt.(*ast.ExpressionStmt)
			if !ok {
				t.Fatalf("stmt is not *ast.ExpressionStmt. Got %T.", stmt)
			}

			expr := expressionStmt.Expression
			testInfixExpr(t, expr, tc.left, tc.operator, tc.right)
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tt := []struct {
		source string
		want   string
	}{
		{"-a + b;", "((-a) + b)"},
		{"-a * b;", "((-a) * b)"},
		{"!-a;", "(!(-a))"},
		{"a + b + c;", "((a + b) + c)"},
		{"a + b - c;", "((a + b) - c)"},
		{"a * b * c;", "((a * b) * c)"},
		{"a * b / c;", "((a * b) / c)"},
		{"a + b / c;", "(a + (b / c))"},
		{"a * b / c;", "((a * b) / c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4);\n((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 / 4 - 5 != -2 / 5 / -6 ", "(((3 / 4) - 5) != (((-2) / 5) / (-6)))"},
		{"true;", "true"},
		{"false;", "false"},
		{"3 > 5 == false;", "((3 > 5) == false)"},
		{"true != 1 > 2;", "(true != (1 > 2))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (1 + 1)", "(2 / (1 + 1))"},
		{"-(1 + 1)", "(-(1 + 1))"},
		{"!(true != false)", "(!(true != false))"},
		{"1 + a(b * c) - 3", "((1 + a((b * c))) - 3)"},
		{"a(1, 2, b, a(b * c, 3))", "a(1, 2, b, a((b * c), 3))"},
		{"a(1 + 2 / 3 - 4)", "a(((1 + (2 / 3)) - 4))"},
		{"a.b.c()", "((a.b).c)()"},
		{"a.b().c()", "((a.b)().c)()"},
		{"a.b.c = 10;", "((a.b).c = 10)"},
		{"a.b().c = 10;", "((a.b)().c = 10)"},
		{"a || b && c", "(a || (b && c))"},
		{
			"x - y || a * b + c || d && e * !f;",
			"(((x - y) || ((a * b) + c)) || (d && (e * (!f))))",
		},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tc := range tt {
		t.Run(tc.source, func(t *testing.T) {
			program := parse(t, tc.source)

			if program == nil {
				t.Fatal("ParseProgram() returned 'nil'.")
			}

			got := strings.Trim(program.String(), ";\n")
			if got != tc.want {
				t.Errorf("Program.String() is wrong, want %q, got %q.", tc.want, got)
			}
		})
	}
}

// Helpers

func testInfixExpr(
	t testing.TB,
	expr ast.Expression,
	left interface{},
	operator string,
	right interface{}) {
	t.Helper()
	infixExpr, ok := expr.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("expr is not *ast.InfixExpr. Got %T.", expr)
	}
	if infixExpr.TokenLiteral() != operator {
		t.Errorf("Wrong InfixExpr.TokenLiteral, want %q, got %q.",
			operator, infixExpr.TokenLiteral())
	}
	if infixExpr.Operator != operator {
		t.Errorf("Wrong Operator, want %q, got %q.",
			operator, infixExpr.Operator)
	}

	testIdentifierOrLiteralExpr(t, infixExpr.Left, left)
	testIdentifierOrLiteralExpr(t, infixExpr.Right, right)
}

func testIdentifierOrLiteralExpr(t testing.TB, expr ast.Expression, want interface{}) {
	t.Helper()
	switch w := want.(type) {
	case int:
		testIntLiteralExpression(t, expr, int64(w))
	case int64:
		testIntLiteralExpression(t, expr, w)
	case bool:
		testBoolLiteralExpression(t, expr, w)
	case string:
		switch expr.(type) {
		case *ast.IdentifierExpr:
			testIdentifierExpression(t, expr, string(w))
		default:
			testStringLiteralExpression(t, expr, string(w))
		}
	case nil:
		testNullLiteralExpression(t, expr)
	}
}

func testIdentifierExpression(t testing.TB, expr ast.Expression, want string) {
	identExpr, ok := expr.(*ast.IdentifierExpr)
	if !ok {
		t.Errorf("Expression is not IdentifierExpr, got %T.", expr)
	}

	if identExpr.Value != want {
		t.Errorf("Wrong IdentifierExpr Value. Got %q, want %q.", identExpr.Value, want)
	}

	if identExpr.TokenLiteral() != want {
		t.Errorf("Wrong TokenLiteral, want %q, got %q.", want, identExpr.TokenLiteral())
	}

	t.Helper()
}

func testIntLiteralExpression(t testing.TB, expr ast.Expression, want int64) {
	t.Helper()

	intLiteralExpr, ok := expr.(*ast.IntLiteralExpr)
	if !ok {
		t.Errorf("Expression is not IntLiteralExpr, got %T.", expr)
	}

	if intLiteralExpr.Value != want {
		t.Errorf("Wrong IntLiteralExpr Value. Got %d, want %d.", intLiteralExpr.Value, want)
	}

	wantStr := strconv.FormatInt(want, 10)
	if intLiteralExpr.TokenLiteral() != wantStr {
		t.Errorf("Wrong TokenLiteral(). Got %q, want %q.", intLiteralExpr.TokenLiteral(), wantStr)
	}
}

func testBoolLiteralExpression(t testing.TB, expr ast.Expression, want bool) {
	t.Helper()

	boolLiteralExpr, ok := expr.(*ast.BoolLiteralExpr)
	if !ok {
		t.Errorf("Expression is not BoolLiteralExpr, got %T.", expr)
	}

	if boolLiteralExpr.Value != want {
		t.Errorf("Wrong BoolLiteralExpr Value. Got %v, want %v.", boolLiteralExpr.Value, want)
	}

	wantStr := strconv.FormatBool(want)
	if boolLiteralExpr.TokenLiteral() != wantStr {
		t.Errorf("Wrong TokenLiteral(). Got %q, want %q.", boolLiteralExpr.TokenLiteral(), wantStr)
	}
}

func testStringLiteralExpression(t testing.TB, expr ast.Expression, want string) {
	t.Helper()

	stringLiteralExpr, ok := expr.(*ast.StringLiteralExpr)
	if !ok {
		t.Errorf("Expression is not StringLiteralExpr, got %T.", expr)
	}

	if stringLiteralExpr.Value != want {
		t.Errorf("Wrong StringLiteralExpr Value. Got %v, want %v.", stringLiteralExpr.Value, want)
	}

	if stringLiteralExpr.TokenLiteral() != want {
		t.Errorf("Wrong TokenLiteral(). Got %q, want %q.", stringLiteralExpr.TokenLiteral(), want)
	}
}

func testNullLiteralExpression(t testing.TB, expr ast.Expression) {
	t.Helper()

	nullExpr, ok := expr.(*ast.NullExpr)
	if !ok {
		t.Errorf("Expression is not NullExpr, got %T.", expr)
	}
	if nullExpr.TokenLiteral() != "null" {
		t.Errorf("Wrong TokenLiteral(). Got %q, want 'null'.", nullExpr.TokenLiteral())
	}
}

func testLetStatement(t testing.TB, stmt ast.Statement, want string) bool {
	t.Helper()

	if stmt.TokenLiteral() != "let" {
		t.Errorf("Token literal is not 'let', got %q.", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStmt)

	if !ok {
		t.Errorf("Statement is not LetStmt, got %T.", stmt)
		return false
	}

	if letStmt.Name.Value != want {
		t.Errorf("Wrong LetStmt Name Value. Got %q, want %q.", letStmt.Name.Value, want)
		return false
	}

	if letStmt.Name.TokenLiteral() != want {
		t.Errorf("Wrong LetStmt Name. Got %q, want %q.", letStmt.Name.TokenLiteral(), want)
		return false
	}

	return true
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

func parse(t testing.TB, source string) *ast.Program {
	t.Helper()
	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()
	ensureNoParserErrors(t, p)

	return program
}
