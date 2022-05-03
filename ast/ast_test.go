package ast_test

import (
	"monkey/ast"
	"monkey/token"
	"strings"
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
			// num = "hello"
			&ast.ExpressionStmt{
				Token: token.Token{},
				Expression: &ast.AssignExpr{
					Token: token.Token{
						Type:    token.ASSIGN,
						Literal: "=",
					},
					Identifier: &ast.IdentifierExpr{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "num",
						},
						Value: "num",
					},
					Expression: &ast.StringLiteralExpr{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "hello",
						},
						Value: "hello",
					},
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

			//fn(x, y) { "string literal"; };
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
									Type:    token.STRING,
									Literal: "string literal",
								},
								Expression: &ast.StringLiteralExpr{
									Token: token.Token{
										Type:    token.STRING,
										Literal: "string literal",
									},
									Value: "string literal",
								},
							},
						},
					},
				},
			},

			// hello(1, 2 + 3);
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.IDENTIFIER,
					Literal: "sum",
				},
				Expression: &ast.CallExpr{
					Token: token.Token{
						Type:    token.LPAREN,
						Literal: "(",
					},
					Function: &ast.IdentifierExpr{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "sum",
						},
						Value: "sum",
					},
					Arguments: []ast.Expression{
						&ast.IntLiteralExpr{
							Token: token.Token{
								Type:    token.INT,
								Literal: "1",
							},
							Value: 1,
						},
						&ast.InfixExpr{
							Token: token.Token{
								Type:    token.PLUS,
								Literal: "+",
							},
							Left: &ast.IntLiteralExpr{
								Token: token.Token{
									Type:    token.INT,
									Literal: "2",
								},
								Value: 2,
							},
							Operator: "+",
							Right: &ast.IntLiteralExpr{
								Token: token.Token{
									Type:    token.INT,
									Literal: "3",
								},
								Value: 3,
							},
						},
					},
				},
			},

			// class Hello < World {
			// 	let method = fn(x) { this;super.method;x; };
			// }
			&ast.ClassStmt{
				Token: token.Token{
					Type:    token.CLASS,
					Literal: "class",
				},
				Name: &ast.IdentifierExpr{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "Hello",
					},
					Value: "Hello",
				},
				Superclass: &ast.IdentifierExpr{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "World",
					},
					Value: "World",
				},
				Methods: []*ast.LetStmt{
					{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.IdentifierExpr{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "method",
							},
							Value: "method",
						},
						Value: &ast.FunctionExpr{
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
							},
							Body: &ast.BlockStmt{
								Token: token.Token{
									Type:    token.LBRACE,
									Literal: "{",
								},
								Statements: []ast.Statement{
									&ast.ExpressionStmt{
										Token: token.Token{
											Type:    token.THIS,
											Literal: "this",
										},
										Expression: &ast.ThisExpr{
											Token: token.Token{
												Type:    token.THIS,
												Literal: "this",
											},
										},
									},
									&ast.ExpressionStmt{
										Token: token.Token{
											Type:    token.SUPER,
											Literal: "super",
										},
										Expression: &ast.SuperExpr{
											Token: token.Token{
												Type:    token.SUPER,
												Literal: "super",
											},
											Method: &ast.IdentifierExpr{
												Token: token.Token{
													Type:    token.IDENTIFIER,
													Literal: "method",
												},
												Value: "method",
											},
										},
									},
									&ast.ExpressionStmt{
										Token: token.Token{
											Type:    token.IDENTIFIER,
											Literal: "x",
										},
										Expression: &ast.IdentifierExpr{
											Token: token.Token{
												Type:    token.IDENTIFIER,
												Literal: "x",
											},
											Value: "x",
										},
									},
								},
							},
						},
					},
				},
			},
			// (obj.field);
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.IDENTIFIER,
					Literal: "obj",
				},
				Expression: &ast.GetExpr{
					Token: token.Token{
						Type:    token.DOT,
						Literal: ".",
					},
					Expression: &ast.IdentifierExpr{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "obj",
						},
						Value: "obj",
					},
					Field: &ast.IdentifierExpr{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "field",
						},
						Value: "field",
					},
				},
			},
			// (obj.field = 10);
			&ast.ExpressionStmt{
				Token: token.Token{
					Type:    token.IDENTIFIER,
					Literal: "obj",
				},
				Expression: &ast.SetExpr{
					Token: token.Token{
						Type:    token.DOT,
						Literal: "=",
					},
					Expression: &ast.IdentifierExpr{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "obj",
						},
						Value: "obj",
					},
					Field: &ast.IdentifierExpr{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "field",
						},
						Value: "field",
					},
					Value: &ast.IntLiteralExpr{
						Token: token.Token{
							Type:    token.INT,
							Literal: "10",
						},
						Value: 10,
					},
				},
			},
		},
	}

	want :=
		`let num = 10;
num = "hello";
return 123;
999;
(!(false == true));
if (!true) { 10; } else { false; };
fn(x, y) { "string literal"; };
sum(1, (2 + 3));
class Hello < World {
	let method = fn(x) { this;super.method;x; };
}
(obj.field);
(obj.field = 10);
`
	if program.String() != want {
		t.Errorf("Wrong program.String() output. Want %q, got %q.", want, program.String())
		wantLines := strings.Split(want, "\n")
		gotLines := strings.Split(program.String(), "\n")

		i := 0
		for ; i < len(wantLines); i++ {
			if i >= len(gotLines) {
				t.Errorf("Want %q, got nothing.", wantLines[i])
				continue
			}

			if gotLines[i] != wantLines[i] {
				t.Errorf("Want %q, got %q.", wantLines[i], gotLines[i])
			}
		}

		for ; i < len(gotLines); i++ {
			t.Errorf("Got %q, want nothing.", gotLines[i])
		}
	}

}
