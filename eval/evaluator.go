package eval

import (
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var Locals = map[ast.Expression]int{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStmt:
		return evalBlockStatement(node.Statements, object.NewEnclosedEnvironment(env))
	case *ast.ExpressionStmt:
		return Eval(node.Expression, env)
	case *ast.ReturnStmt:
		var val object.Object

		if node.Value == nil {
			val = NULL
		} else {
			val = Eval(node.Value, env)
		}

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.LetStmt:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val
	case *ast.ClassStmt:
		return evalClassStmt(node, env)
	case *ast.ThisExpr:
		return lookupVariable(token.THIS_KEYWORD, node, env)
	case *ast.SuperExpr:
		return evalSuperExpr(node, env)
	case *ast.IntLiteralExpr:
		return &object.Integer{Value: node.Value}
	case *ast.BoolLiteralExpr:
		return boolToBooleanObject(node.Value)
	case *ast.StringLiteralExpr:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteralExpr:
		arr := &object.Array{}
		for _, el := range node.Elements {
			val := Eval(el, env)
			if isError(val) {
				return val
			}
			arr.Elements = append(arr.Elements, val)
		}
		return arr
	case *ast.NullExpr:
		return NULL
	case *ast.PrefixExpr:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		if node.Token.Type == token.OR || node.Token.Type == token.AND {
			return evalLogicalExpr(node.Left, node.Operator, node.Right, env)
		}

		left := Eval(node.Left, env)
		// TODO: definitely need a more elegant way to handle errors
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpr(left, node.Operator, right)
	case *ast.IndexExpr:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		idx := Eval(node.Index, env)
		if isError(idx) {
			return idx
		}

		if left.Type() != object.ARRAY_OBJ || idx.Type() != object.INTEGER_OBJ {
			return indexOperatorError(left.Type(), idx.Type())
		}

		arr := left.(*object.Array)
		l := int64(len(arr.Elements))
		i := idx.(*object.Integer).Value
		if i >= 0 && i < l {
			return arr.Elements[i]
		} else if i < 0 && l + i >= 0 {
			return arr.Elements[l + i]
		} else {
			return outOfBoundsError(left.Type(), i)
		}
	case *ast.GetExpr:
		return evalGetExpr(node, env)
	case *ast.SetExpr:
		return evalSetExpr(node, env)
	case *ast.AssignExpr:
		val := Eval(node.Expression, env)
		if isError(val) {
			return val
		}

		if depth, ok := Locals[node.Identifier]; ok {
			env.AssignAt(depth, node.Identifier.Value, val)
		} else {
			return identifierNotFoundError(node.Identifier.Value)
		}

		return val
	case *ast.IfExpr:
		return evalIfExpr(node, env)
	case *ast.IdentifierExpr:
		return evalIdentifier(node, env)
	case *ast.FunctionExpr:
		return evalFunctionExpr(node, env, false)
	case *ast.CallExpr:
		return evalCallExpr(node, env)
	default:
		return NULL
	}
}

func evalProgram(statements []ast.Statement, env *object.Environment) (result object.Object) {
	for _, stmt := range statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	if result == nil {
		result = NULL
	}

	return
}

func evalBlockStatement(statements []ast.Statement, env *object.Environment) (result object.Object) {
	for _, stmt := range statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	if result == nil {
		result = NULL
	}

	return
}

func evalClassStmt(node *ast.ClassStmt, env *object.Environment) object.Object {
	var super *object.Class = nil
	if node.Superclass != nil {
		sc := Eval(node.Superclass, env)

		if isError(sc) {
			return sc
		}

		if superclass, ok := sc.(*object.Class); ok {
			super = superclass
		} else {
			return superclassMustBeClassError(node.Name.Value, sc.Type())
		}
	}

	// env.Set(node.Name.Value, nil) // declare

	if super != nil {
		env = object.NewEnclosedEnvironment(env)
		env.Set(token.SUPER_KEYWORD, super)
	}

	methods := make(map[string]*object.Function)
	for _, field := range node.Methods {
		if method, ok := field.Value.(*ast.FunctionExpr); ok {
			isInit := false
			if field.Name.Value == token.INITIALIZER_KEYWORD {
				isInit = true
			}
			fn := evalFunctionExpr(method, env, isInit)
			methods[field.Name.Value] = fn
		}
	}

	class := &object.Class{
		Name:    node.Name,
		Super:   super,
		Methods: methods,
	}

	if super != nil {
		env = env.Outer
	}

	env.Set(node.Name.Value, class)
	return class
}

func evalGetExpr(node *ast.GetExpr, env *object.Environment) object.Object {
	obj := Eval(node.Expression, env)
	if isError(obj) {
		return obj
	}

	if obj.Type() != object.INSTANCE_OBJ {
		return wrongGetTargetError(obj.Type(), node.Field.Value)
	}

	inst := obj.(*object.Instance)
	if field, ok := inst.Fields[node.Field.Value]; ok {
		return field
	} else if method := inst.Class.FindMethod(node.Field.Value); method != nil {
		return method.Bind(inst)
	} else {
		return undefinedPropertyError(node.Field.Value)
	}
}

func evalSetExpr(node *ast.SetExpr, env *object.Environment) object.Object {
	obj := Eval(node.Expression, env)
	if isError(obj) {
		return obj
	}

	inst, ok := obj.(*object.Instance)
	if !ok {
		return wrongSetTargetError(obj.Type(), node.Field.Value)
	}

	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	inst.Fields[node.Field.Value] = val

	return val
}

func evalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpr(right)
	case token.MINUS:
		return evalMinusOperatorExpr(right)
	default:
		return unknownPrefixOperatorError(operator, "")
	}
}

func evalBangOperatorExpr(right object.Object) object.Object {
	switch right {
	case NULL:
		return TRUE
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return unknownPrefixOperatorError("-", right.Type())
	}

	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalLogicalExpr(
	leftExpr ast.Expression,
	operator string,
	rightExpr ast.Expression,
	env *object.Environment,
) object.Object {
	left := Eval(leftExpr, env)

	if isError(left) {
		return left
	}

	if operator == token.AND {
		if !isTruthy(left) {
			return left
		}
	} else { // operator == token.OR
		if isTruthy(left) {
			return left
		}
	}

	right := Eval(rightExpr, env)
	return right
}

func evalInfixExpr(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpr(left, operator, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpr(left, operator, right)
	case operator == token.EQUAL_EQUAL:
		return boolToBooleanObject(left == right)
	case operator == token.NOT_EQUAL:
		return boolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return infixTypeMismatchError(left.Type(), operator, right.Type())
	default:
		return unknownInfixOperatorError(left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpr(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftValue + rightValue}
	case token.MINUS:
		return &object.Integer{Value: leftValue - rightValue}
	case token.STAR:
		return &object.Integer{Value: leftValue * rightValue}
	case token.SLASH:
		return &object.Integer{Value: leftValue / rightValue}
	case token.GREATER:
		return boolToBooleanObject(leftValue > rightValue)
	case token.GREATER_EQUAL:
		return boolToBooleanObject(leftValue >= rightValue)
	case token.LESS:
		return boolToBooleanObject(leftValue < rightValue)
	case token.LESS_EQUAL:
		return boolToBooleanObject(leftValue <= rightValue)
	case token.EQUAL_EQUAL:
		return boolToBooleanObject(leftValue == rightValue)
	case token.NOT_EQUAL:
		return boolToBooleanObject(leftValue != rightValue)
	default:
		return unknownInfixOperatorError(left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpr(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	switch operator {
	case token.PLUS:
		return &object.String{Value: leftValue + rightValue}
	case token.EQUAL_EQUAL:
		return boolToBooleanObject(leftValue == rightValue)
	case token.NOT_EQUAL:
		return boolToBooleanObject(leftValue != rightValue)
	default:
		return unknownInfixOperatorError(left.Type(), operator, right.Type())
	}
}

func evalIfExpr(expr *ast.IfExpr, env *object.Environment) object.Object {
	condition := Eval(expr.Condition, env)

	if isTruthy(condition) {
		return Eval(expr.Then, env)
	} else if expr.Else != nil {
		return Eval(expr.Else, env)
	} else {
		return NULL
	}
}

func evalSuperExpr(node *ast.SuperExpr, env *object.Environment) object.Object {
	depth, ok := Locals[node]
	if !ok {
		return internalResolveError(node.String())
	}

	super, ok := env.GetAt(depth, token.SUPER_KEYWORD)
	if !ok {
		return internalResolveError(node.String())
	}

	superClass, ok := super.(*object.Class)
	if !ok {
		return internalResolveError(node.String())
	}

	inst, ok := env.GetAt(depth-1, token.THIS_KEYWORD)
	if !ok {
		return internalResolveError(node.String())
	}

	instObj, ok := inst.(*object.Instance)
	if !ok {
		return internalResolveError(node.String())
	}

	fn := superClass.FindMethod(node.Method.Value)
	if fn == nil {
		return undefinedPropertyError(node.Method.Value)
	}

	return fn.Bind(instObj)
}

func evalIdentifier(node *ast.IdentifierExpr, env *object.Environment) object.Object {
	return lookupVariable(node.Value, node, env)
}

func lookupVariable(name string, node ast.Expression, env *object.Environment) object.Object {
	if depth, ok := Locals[node]; ok {
		if val, ok := env.GetAt(depth, name); ok {
			return val
		}
	} else if val, ok := env.GetGlobal(name); ok {
		return val
	}

	if builtin, ok := builtins[name]; ok {
		return builtin
	}

	return identifierNotFoundError(name)
}

func evalFunctionExpr(node *ast.FunctionExpr, env *object.Environment, isInit bool) *object.Function {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
		IsInit:     isInit,
	}
}

func evalCallExpr(node *ast.CallExpr, env *object.Environment) object.Object {
	fn := Eval(node.Function, env)
	if isError(fn) {
		return fn
	}
	args := evalExpressions(node.Arguments, env)
	if len(args) > 0 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(fn, args)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) (result []object.Object) {
	for _, expr := range expressions {
		arg := Eval(expr, env)
		if isError(arg) {
			return []object.Object{arg}
		}
		result = append(result, arg)
	}
	return
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return wrongArgumentsCountError(len(fn.Parameters), len(args))
		}

		extendedEnv := extendFunctionEnv(fn, args)
		result := evalBlockStatement(fn.Body.Statements, extendedEnv)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			if fn.IsInit {
				this, ok := fn.Env.GetAt(0, token.THIS_KEYWORD)
				if ok {
					return this
				}
				return internalResolveError(token.THIS_KEYWORD)
			}
			return returnValue.Value
		}

		if fn.IsInit {
			this, ok := fn.Env.GetAt(0, token.THIS_KEYWORD)
			if ok {
				return this
			}
			return internalResolveError(token.THIS_KEYWORD)
		}

		return result

	case *object.Builtin:
		return fn.Fn(args...)

	case *object.Class:
		init := fn.FindMethod(token.INITIALIZER_KEYWORD)

		expectArgs := 0
		if init != nil {
			expectArgs = len(init.Parameters)
		}

		if len(args) != expectArgs {
			return wrongArgumentsCountError(expectArgs, len(args))
		}

		inst := &object.Instance{
			Class:  fn,
			Fields: make(map[string]object.Object),
		}

		if init != nil {
			applyFunction(init.Bind(inst), args)
		}

		return inst
	}

	return notAFunctionError(string(fn.Type()), fn.Inspect())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, arg := range args {
		env.Set(fn.Parameters[i].Value, arg)
	}
	return env
}

func boolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}
