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
				Expression: &ast.PrefixExpr{
					Token: token.Token{
						Type:    token.BANG,
						Literal: "!",
					},
					Operator: "!",
					Right: &ast.InfixExpr{
						Token: token.Token{
							Type:    token.EQUAL_EQUAL,
							Literal: "==",
						},
						Left: &ast.BoolLiteralExpr{
							Token: token.Token{
								Type:    token.FALSE,
								Literal: "false",
							},
							Value: false,
						},
						Operator: "==",
						Right: &ast.BoolLiteralExpr{
							Token: token.Token{
								Type:    token.TRUE,
								Literal: "true",
							},
							Value: true,
						},
					},
				},
			},
			// if (!true) { 10; } else { false; };
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.IF,
					Literal: "if",
				},
				Expression: &ast.IfExpr{
					Token: token.Token{
						Type:    token.IF,
						Literal: "if",
					},
					Condition: &ast.PrefixExpr{
						Token: token.Token{
							Type:    token.BANG,
							Literal: "!",
						},
						Operator: "!",
						Right: &ast.BoolLiteralExpr{
							Token: token.Token{
								Type:    token.TRUE,
								Literal: "true",
							},
							Value: false,
						},
					},
					Then: &ast.BlockStmt{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []ast.Statement{
							&ast.ExpressionStmt{
								Token: token.Token{
									Type:    token.INT,
									Literal: "10",
								},
								Expression: &ast.IntLiteralExpr{
									Token: token.Token{
										Type:    token.INT,
										Literal: "10",
									},
									Value: 10,
								},
							},
						},
					},
					Else: &ast.BlockStmt{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []ast.Statement{
							&ast.ExpressionStmt{
								Token: token.Token{
									Type:    token.FALSE,
									Literal: "false",
								},
								Expression: &ast.BoolLiteralExpr{
									Token: token.Token{
										Type:    token.FALSE,
										Literal: "false",
									},
									Value: false,
								},
							},
						},
					},
				},
			},

			//fn(x, y) { 10; };
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.FUNCTION,
					Literal: "fn",
				},
				Expression: &ast.FunctionExpr{
					Token: token.Token{
						Type:    token.FUNCTION,
						Literal: "fn",
					},
					Parameters: []*ast.IdentifierExpr{
						{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "x",
							},
							Value: "x",
						},
						{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "y",
							},
							Value: "y",
						},
					},
					Body: &ast.BlockStmt{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []ast.Statement{
							&ast.ExpressionStmt{
								Token: token.Token{
									Type:    token.INT,
									Literal: "10",
								},
								Expression: &ast.IntLiteralExpr{
									Token: token.Token{
										Type:    token.INT,
										Literal: "10",
									},
									Value: 10,
								},
							},
						},
					},
				},
			},
		},
	}

	want :=
		`let num = 10;
return 123;
999;
(!(false == true));
if (!true) { 10; } else { false; };
fn(x, y) { 10; };
`
	if program.String() != want {
		t.Errorf("Wrong program.String() output. Want %q, got %q.", want, program.String())
	}
}
