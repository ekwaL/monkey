package ast_test

import (
	"monkey/ast"
	"monkey/token"
	"testing"
)

func TestAstString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStmt{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.IdentifierExpr{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "num",
					},
					Value: "num",
				},
				Value: &ast.IntLiteralExpr{
					Token: token.Token{
						Type:    token.INT,
						Literal: "10",
					},
					Value: 10,
				},
			},
			&ast.ReturnStmt{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
				},
				Value: &ast.IntLiteralExpr{
					Token: token.Token{
						Type:    token.INT,
						Literal: "123",
					},
					Value: 123,
				},
			},
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.INT,
					Literal: "999",
				},
				Expression: &ast.IntLiteralExpr{
					Token: token.Token{
						Type:    token.INT,
						Literal: "999",
					},
					Value: 999,
				},
			},
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.INT,
					Literal: "false",
				},
				Expression: &ast.BoolLiteralExpr{
					Token: token.Token{
						Type:    token.INT,
						Literal: "false",
					},
					Value: false,
				},
			},
		},
	}

	want :=
		`let num = 10;
return 123;
999;
false;
`
	if program.String() != want {
		t.Errorf("Wrong program.String() output. Want %q, got %q.", want, program.String())
	}
}
