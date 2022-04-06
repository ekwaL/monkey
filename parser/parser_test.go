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
		`
	program := parse(t, source)

	if program == nil {
		t.Fatal("ParseProgram() returned 'nil'.")
	}

	wantLen := 3
	if len(program.Statements) != wantLen {
		t.Fatalf("program.Statements len is %d, want %d .", len(program.Statements), wantLen)
	}

	tt := []struct {
		expectedIdentifier string
		expectedValue      int64
	}{
		{expectedIdentifier: "x", expectedValue: 5},
		{expectedIdentifier: "y", expectedValue: 10},
		{expectedIdentifier: "foobarbaz", expectedValue: 10},
	}

	for i, tc := range tt {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tc.expectedIdentifier) {
			return
		}
		testIdentifierOrLiteralExpr(t, stmt.(*ast.LetStmt).Value, tc.expectedValue)
	}

}

func TestReturnStatement(t *testing.T) {
	source := `
		return 5;
		return 9999;
		`
	program := parse(t, source)
	expect := []int64{5, 9999}

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

		testIdentifierOrLiteralExpr(t, returnStmt.Value, expectedValue)
	}
}

func TestBlockStatement(t *testing.T) {

}

func TestIfExpression(t *testing.T) {
	source := `
		if (5 > 10) {
			20;
		} else true;`
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
	if elseBranch.TokenLiteral() != "true" {
		t.Errorf("Wrong TokenLiteral, want 'true', got %q.", elseBranch.TokenLiteral())
	}
	testIdentifierOrLiteralExpr(t, elseBranch.Expression, true)
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
		t.Errorf("Wrong parameters count: want 1, got %d", len(fnExpr.Parameters))
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
		testIdentifierExpression(t, expr, string(w))
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
