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
		if stmt == nil {
			out.WriteString("nil")
		} else {
			out.WriteString(stmt.String())
		}
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
	Token      token.Token // first token of Expression
	Expression Expression
}

func (e *ExpressionStmt) statementNode()       {}
func (e *ExpressionStmt) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStmt) String() string {
	var out bytes.Buffer

	if e.Expression == nil {
		out.WriteString("nil")
	} else {
		out.WriteString(e.Expression.String())
	}
	out.WriteString(";")

	return out.String()
}

type BlockStmt struct {
	Token      token.Token // '{'
	Statements []Statement
}

func (b *BlockStmt) statementNode()       {}
func (b *BlockStmt) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStmt) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")

	return out.String()
}

// Expressions

type IdentifierExpr struct {
	Token token.Token
	Value string
}

func (i *IdentifierExpr) expressionNode()      {}
func (i *IdentifierExpr) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierExpr) String() string       { return i.Value }

type IntLiteralExpr struct {
	Token token.Token
	Value int64
}

func (i *IntLiteralExpr) expressionNode()      {}
func (i *IntLiteralExpr) TokenLiteral() string { return i.Token.Literal }
func (i *IntLiteralExpr) String() string       { return i.TokenLiteral() }

type BoolLiteralExpr struct {
	Token token.Token
	Value bool
}

func (b *BoolLiteralExpr) expressionNode()      {}
func (b *BoolLiteralExpr) TokenLiteral() string { return b.Token.Literal }
func (b *BoolLiteralExpr) String() string       { return b.TokenLiteral() }

type StringLiteralExpr struct {
	Token token.Token
	Value string
}

func (s *StringLiteralExpr) expressionNode()      {}
func (s *StringLiteralExpr) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteralExpr) String() string       { return `"` + s.TokenLiteral() + `"` }

type PrefixExpr struct {
	Token    token.Token // prefix token
	Operator string
	Right    Expression
}

func (p *PrefixExpr) expressionNode()      {}
func (p *PrefixExpr) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpr struct {
	Token    token.Token // operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (p *InfixExpr) expressionNode()      {}
func (p *InfixExpr) TokenLiteral() string { return p.Token.Literal }
func (p *InfixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Left.String())
	out.WriteString(" " + p.Operator + " ")
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type AssignExpr struct {
	Token      token.Token // =
	Identifier *IdentifierExpr
	Expression Expression
}

func (a *AssignExpr) expressionNode()      {}
func (a *AssignExpr) TokenLiteral() string { return a.Token.Literal }
func (a *AssignExpr) String() string {
	var out bytes.Buffer
	out.WriteString(a.Identifier.String())
	out.WriteString(" ")
	out.WriteString(a.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(a.Expression.String())
	return out.String()
}

type IfExpr struct {
	Token     token.Token
	Condition Expression
	Then      Statement
	Else      Statement
}

func (i *IfExpr) expressionNode()      {}
func (i *IfExpr) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Then.String())
	if i.Else != nil {
		out.WriteString(" else ")
		out.WriteString(i.Else.String())
	}

	return out.String()
}

type FunctionExpr struct {
	Token      token.Token
	Parameters []*IdentifierExpr
	Body       *BlockStmt
}

func (f *FunctionExpr) expressionNode()      {}
func (f *FunctionExpr) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionExpr) String() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	last := len(f.Parameters) - 1
	for i, ident := range f.Parameters {
		out.WriteString(ident.Value)
		if i != last {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")

	out.WriteString(f.Body.String())

	return out.String()
}

type CallExpr struct {
	Token     token.Token // '('
	Function  Expression  // IdentifierExpr || FunctionExpr
	Arguments []Expression
}

func (c *CallExpr) expressionNode()      {}
func (c *CallExpr) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpr) String() string {
	var out bytes.Buffer

	out.WriteString(c.Function.String())
	out.WriteString("(")

	last := len(c.Arguments) - 1
	for i, arg := range c.Arguments {
		out.WriteString(arg.String())
		if i != last {
			out.WriteString(", ")
		}
	}

	out.WriteString(")")
	return out.String()
}
