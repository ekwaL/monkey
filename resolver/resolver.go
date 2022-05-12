package resolver

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"monkey/utils"
)

type FnType int
type ClassType int

const (
	_ = iota
	NONE
	FUNCTION
	METHOD
	INITIALIZER
	CLASS
	SUBCLASS
)

const (
	ERR_READ_ON_OWN_INIT   = "Can not read variable '%s' in it's own initializer."
	ERR_ALREADY_DECLARED   = "Variable '%s' is already declared in current scope."
	ERR_CLASS_INHERIT_SELF = "Class '%s' can not inherit from itself."
	// ERR_TOP_LEVEL_RETURN         = "Can not return from top-level code."
	ERR_INITIALIZER_VAL_RETURN   = "Can not return value from initializer."
	ERR_THIS_OUTSIDE_OF_CLASS    = "Can not use 'this' outside of class."
	ERR_SUPER_OUTSIDE_OF_CLASS   = "Can not use 'super' outside of class."
	ERR_SUPER_WITHOUT_SUPERCLASS = "Can not use 'super' in a class with no superclass."
)

type resolver struct {
	scopes utils.Stack[map[string]bool]
	locals map[ast.Expression]int
	errors []string

	currFn    FnType
	currClass ClassType
}

func New() *resolver {
	r := &resolver{
		scopes:    utils.NewStack[map[string]bool](),
		locals:    make(map[ast.Expression]int),
		errors:    []string{},
		currFn:    NONE,
		currClass: NONE,
	}

	r.beginScope()

	return r
}

func (r *resolver) Locals() map[ast.Expression]int {
	return r.locals
}

func (r *resolver) Errors() []string {
	return r.errors
}

func (r *resolver) Resolve(node ast.Node) {
	switch node := node.(type) {
	case *ast.Program:
		r.resolveStatements(node.Statements)
	case *ast.BlockStmt:
		r.beginScope()
		r.resolveStatements(node.Statements)
		r.endScope()
	case *ast.ExpressionStmt:
		r.Resolve(node.Expression)
	case *ast.ReturnStmt:
		// if r.currFn == NONE {
		// 	r.error(ERR_TOP_LEVEL_RETURN)
		// }
		if node.Value != nil {
			if r.currFn == INITIALIZER {
				r.error(ERR_INITIALIZER_VAL_RETURN)
			}
			r.Resolve(node.Value)
		}
	case *ast.LetStmt:
		switch node.Value.(type) {
		case *ast.FunctionExpr:
			r.declare(node.Name)
			r.define(node.Name)
			r.Resolve(node.Value)
		default:
			r.declare(node.Name)
			r.Resolve(node.Value)
			r.define(node.Name)
		}
	case *ast.ClassStmt:
		enclosingClass := r.currClass
		r.currClass = CLASS
		r.declare(node.Name)
		r.define(node.Name)

		if node.Superclass != nil {
			r.currClass = SUBCLASS

			if node.Superclass.Value == node.Name.Value {
				r.error(fmt.Sprintf(ERR_CLASS_INHERIT_SELF, node.Name.Value))
			} else {
				r.Resolve(node.Superclass)
			}

			r.beginScope()
			r.defineName(token.SUPER_KEYWORD)
		}

		r.beginScope()
		r.defineName(token.THIS_KEYWORD)

		for _, field := range node.Methods {
			// parser only allows let expressions with functionExpr values for now
			method, _ := field.Value.(*ast.FunctionExpr)
			var methodType FnType = METHOD

			if field.Name.Value == token.INITIALIZER_KEYWORD {
				methodType = INITIALIZER
			}

			r.resolveFn(method, methodType)
		}

		r.endScope()

		if node.Superclass != nil {
			r.endScope()
		}

		r.currClass = enclosingClass

	case *ast.ThisExpr:
		if r.currClass == NONE {
			r.error(ERR_THIS_OUTSIDE_OF_CLASS)
		}
		r.resolveLocal(node, token.THIS_KEYWORD)
	case *ast.SuperExpr:
		if r.currClass == NONE {
			r.error(ERR_SUPER_OUTSIDE_OF_CLASS)
		} else if r.currClass == CLASS {
			r.error(ERR_SUPER_WITHOUT_SUPERCLASS)
		}
		r.resolveLocal(node, token.SUPER_KEYWORD)
	case *ast.GetExpr:
		r.Resolve(node.Expression)
	case *ast.SetExpr:
		r.Resolve(node.Value)
		r.Resolve(node.Expression)
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
		r.resolveFn(node, FUNCTION)
	case *ast.CallExpr:
		r.Resolve(node.Function)
		for _, a := range node.Arguments {
			r.Resolve(a)
		}
	case *ast.IndexExpr:
		r.Resolve(node.Left)
		r.Resolve(node.Index)
	case *ast.ArrayLiteralExpr:
		for _, el := range node.Elements {
			r.Resolve(el)
		}
	case *ast.IntLiteralExpr:
	case *ast.BoolLiteralExpr:
	case *ast.StringLiteralExpr:
	case *ast.NullExpr:
	}
}

func (r *resolver) resolveFn(fn *ast.FunctionExpr, t FnType) {
	enclosingFn := r.currFn
	r.currFn = t

	r.beginScope()

	for _, p := range fn.Parameters {
		r.declare(p)
		r.define(p)
	}
	r.resolveStatements(fn.Body.Statements)

	r.endScope()
	r.currFn = enclosingFn
}

func (r *resolver) resolveVariable(name *ast.IdentifierExpr) {
	r.resolveLocal(name, name.Value)
}

func (r *resolver) resolveLocal(expr ast.Expression, name string) {
	scopes := r.scopes.List()

	foundUndefined := false
	for i := len(scopes) - 1; i >= 0; i-- {
		if defined, ok := scopes[i][name]; ok {
			if defined {
				r.locals[expr] = len(scopes) - 1 - i
				return
			} else {
				// variable 'x' appears on the right side of initializer of variable 'x'
				foundUndefined = true
			}
		}
	}

	// no other 'x' variables found in outer scopes
	if foundUndefined {
		r.error(fmt.Sprintf(ERR_READ_ON_OWN_INIT, name))
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
	r.defineName(name.Value)
}

func (r *resolver) defineName(name string) {
	currScope, ok := r.scopes.Peek()
	if !ok {
		return
	}

	currScope[name] = true
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
