package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string { return "" }
func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}

	return out.String()
}

// Statements

type LetStmt struct {
	Token token.Token
	Name  *IdentifierExpr
	Value Expression
}

func (ls *LetStmt) statementNode()       {}
func (ls *LetStmt) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStmt) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())
	out.WriteString(";")

	return out.String()
}

type ReturnStmt struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStmt) statementNode()       {}
func (r *ReturnStmt) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStmt) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(r.Value.String())
	out.WriteString(";")

	return out.String()
}

type ExpressionStmt struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStmt) statementNode()       {}
func (e *ExpressionStmt) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStmt) String() string {
	var out bytes.Buffer

	out.WriteString(e.Expression.String())
	out.WriteString(";")

	return out.String()
}

// Expressions

type IdentifierExpr struct {
	Token token.Token
	Value string
}

func (i *IdentifierExpr) expressionNode()      {}
func (i *IdentifierExpr) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierExpr) String() string {
	return i.Value
}

type IntLiteralExpr struct {
	Token token.Token
	Value int
}

func (i *IntLiteralExpr) expressionNode()      {}
func (i *IntLiteralExpr) TokenLiteral() string { return i.Token.Literal }
func (i *IntLiteralExpr) String() string {
	return i.TokenLiteral() // value.String()?
}
