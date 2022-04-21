package eval

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStmt:
		return evalBlockStatement(node.Statements, env)
	case *ast.ExpressionStmt:
		return Eval(node.Expression, env)
	case *ast.ReturnStmt:
		val := Eval(node.Value, env)
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
	case *ast.IntLiteralExpr:
		return &object.Integer{Value: node.Value}
	case *ast.BoolLiteralExpr:
		return boolToBooleanObject(node.Value)
	case *ast.StringLiteralExpr:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpr:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
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
	case *ast.IfExpr:
		return evalIfExpr(node, env)
	case *ast.IdentifierExpr:
		return evalIdentifier(node, env)
	case *ast.FunctionExpr:
		return evalFunctionExpr(node, env)
	case *ast.CallExpr:
		return evalCallExpr(node, env)
	default:
		println(fmt.Sprintf("eval is unimplemented for %T", node))
		return nil
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
	return
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
		return typeMismatchError(left.Type(), operator, right.Type())
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

func evalIdentifier(node *ast.IdentifierExpr, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return identifierNotFound(node.Value)
	}
	return val
}

func evalFunctionExpr(node *ast.FunctionExpr, env *object.Environment) *object.Function {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
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
	function, ok := fn.(*object.Function)
	if !ok {
		return notAFunction(string(fn.Type()), fn.Inspect())
	}
	if len(function.Parameters) != len(args) {
		return wrongArgumentsCount(len(function.Parameters), len(args))
	}

	extendedEnv := extendFunctionEnv(function, args)
	result := Eval(function.Body, extendedEnv)
	if returnValue, ok := result.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return result
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
	return obj.Type() == object.ERROR_OBJ
}
