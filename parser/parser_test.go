package parser_test

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"strconv"
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
		testIntLiteralExpression(t, stmt.(*ast.LetStmt).Value, tc.expectedValue)
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

		testIntLiteralExpression(t, returnStmt.Value, expectedValue)
	}
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

	testIdentifierExpression(t, expressionStmt.Expression, want)
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

	testIntLiteralExpression(t, expressionStmt.Expression, wantInt)
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
