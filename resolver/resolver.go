package resolver

import (
	"fmt"
	"monkey/ast"
	"monkey/utils"
)

const ERR_READ_ON_OWN_INIT = "Can not read variable '%s' in it's own initializer."
const ERR_ALREADY_DECLARED = "Variable '%s' is already declared in current scope."

type resolver struct {
	scopes utils.Stack[map[string]bool]
	locals map[*ast.IdentifierExpr]int
	errors []string
}

func New() *resolver {
	return &resolver{
		scopes: utils.NewStack[map[string]bool](),
		locals: make(map[*ast.IdentifierExpr]int),
		errors: []string{},
	}
}

func (r *resolver) Locals() map[*ast.IdentifierExpr]int {
	return r.locals
}

func (r *resolver) Errors() []string {
	return r.errors
}

func (r *resolver) Resolve(node ast.Node) {
	switch node := node.(type) {
	case *ast.Program:
		r.beginScope()
		r.resolveStatements(node.Statements)
		r.endScope()
	case *ast.BlockStmt:
		r.beginScope()
		r.resolveStatements(node.Statements)
		r.endScope()
	case *ast.ExpressionStmt:
		r.Resolve(node.Expression)
	case *ast.ReturnStmt:
		r.Resolve(node.Value)
	case *ast.LetStmt:
		r.declare(node.Name)
		r.Resolve(node.Value)
		r.define(node.Name)
	case *ast.PrefixExpr:
		r.Resolve(node.Right)
	case *ast.InfixExpr:
		r.Resolve(node.Left)
		r.Resolve(node.Right)
	case *ast.AssignExpr:
		r.Resolve(node.Expression)
		r.resolveVariable(node.Identifier)
	case *ast.IfExpr:
		r.Resolve(node.Condition)
		r.Resolve(node.Then)
		r.Resolve(node.Else)
	case *ast.IdentifierExpr:
		r.resolveVariable(node)
	case *ast.FunctionExpr:
		r.beginScope()
		for _, p := range node.Parameters {
			r.define(p)
		}
		r.resolveStatements(node.Body.Statements)
		r.endScope()
	case *ast.CallExpr:
		for _, a := range node.Arguments {
			r.Resolve(a)
		}
		r.Resolve(node.Function)
	case *ast.IntLiteralExpr:
	case *ast.BoolLiteralExpr:
	case *ast.StringLiteralExpr:
	}
}

func (r *resolver) resolveVariable(name *ast.IdentifierExpr) {
	scopes := r.scopes.List()

	foundUndefined := false
	for i := len(scopes) - 1; i >= 0; i-- {
		if defined, ok := scopes[i][name.Value]; ok {
			if defined {
				r.locals[name] = len(scopes) - 1 - i
				return
			} else {
				// variable 'x' appears on the right side of initializer of variable 'x'
				foundUndefined = true
			}
		}
	}

	// no other 'x' variables found in outer scopes
	if foundUndefined {
		r.error(fmt.Sprintf(ERR_READ_ON_OWN_INIT, name.Value))
	}
	// assume variable is global/builtin and do nothing ( and crash at runtime :D )
}

func (r *resolver) resolveStatements(statement []ast.Statement) {
	for _, stmt := range statement {
		r.Resolve(stmt)
	}
}

func (r *resolver) declare(name *ast.IdentifierExpr) {
	currScope, ok := r.scopes.Peek()
	if !ok {
		return
	}

	if _, alreadyDeclared := currScope[name.Value]; alreadyDeclared {
		r.error(fmt.Sprintf(ERR_ALREADY_DECLARED, name.Value))
	} else {
		currScope[name.Value] = false
	}
}

func (r *resolver) define(name *ast.IdentifierExpr) {
	currScope, ok := r.scopes.Peek()
	if !ok {
		return
	}

	currScope[name.Value] = true
}

func (r *resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *resolver) endScope() {
	r.scopes.Pop()
}

func (r *resolver) error(err string) {
	r.errors = append(r.errors, err)
}
