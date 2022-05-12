package object

import (
	"bytes"
	"monkey/ast"
	"monkey/token"
	"strconv"
	"strings"
)

type ObjectType string
type BuiltinFunction func(args ...Object) Object

const (
	NULL_OBJ         = "NULL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	CLASS_OBJ        = "CLASS"
	INSTANCE_OBJ     = "INSTANCE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	els := []string{}
	for _, e := range a.Elements {
		els = append(els, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(els, ", "))
	out.WriteString("]")

	return out.String()
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "Runtime error: " + e.Message }

type Function struct {
	Parameters []*ast.IdentifierExpr
	Body       *ast.BlockStmt
	Env        *Environment
	IsInit     bool
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
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

func (f *Function) Bind(inst *Instance) *Function {
	env := NewEnclosedEnvironment(f.Env)
	env.Set(token.THIS_KEYWORD, inst)
	return &Function{
		Parameters: f.Parameters,
		Body:       f.Body,
		Env:        env,
		IsInit:     f.IsInit,
	}
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Class struct {
	Name    *ast.IdentifierExpr
	Super   *Class
	Methods map[string]*Function
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string {
	return "<class " + c.Name.Value + ">"
}

func (c *Class) FindMethod(name string) *Function {
	fn, ok := c.Methods[name]
	if ok {
		return fn
	}

	if c.Super != nil {
		return c.Super.FindMethod(name)
	}

	return nil
}

type Instance struct {
	Class  *Class
	Fields map[string]Object
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string {
	return "<instance of " + i.Class.Name.Value + ">"
}
