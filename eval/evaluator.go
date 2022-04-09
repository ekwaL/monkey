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

const (
	ERR_UNKNOWN_OPERATOR = "unknown operator: "
	ERR_TYPE_MISMATCH    = "type mismatch: "
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.BlockStmt:
		return evalBlockStatement(node.Statements)
	case *ast.ExpressionStmt:
		return Eval(node.Expression)
	case *ast.ReturnStmt:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IntLiteralExpr:
		return &object.Integer{Value: node.Value}
	case *ast.BoolLiteralExpr:
		return boolToBooleanObject(node.Value)
	case *ast.PrefixExpr:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		left := Eval(node.Left)
		// TODO: definitely need a more elegant way to handle errors
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpr(left, node.Operator, right)
	case *ast.IfExpr:
		return evalIfExpr(node)
	default:
		println(fmt.Sprintf("eval is unimplemented for %T", node))
		return nil
	}
}

func evalProgram(statements []ast.Statement) (result object.Object) {
	for _, stmt := range statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return
}

func evalBlockStatement(statements []ast.Statement) (result object.Object) {
	for _, stmt := range statements {
		result = Eval(stmt)

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

func evalIfExpr(expr *ast.IfExpr) object.Object {
	condition := Eval(expr.Condition)

	if isTruthy(condition) {
		return Eval(expr.Then)
	} else if expr.Else != nil {
		return Eval(expr.Else)
	} else {
		return NULL
	}
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

func unknownPrefixOperatorError(operator string, right object.ObjectType) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_UNKNOWN_OPERATOR+"%s%s", operator, right),
	}
}

func unknownInfixOperatorError(
	left object.ObjectType,
	operator string,
	right object.ObjectType) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_UNKNOWN_OPERATOR+"%s %s %s", left, operator, right),
	}
}

func typeMismatchError(
	left object.ObjectType,
	operator string,
	right object.ObjectType) *object.Error {

	return &object.Error{
		Message: fmt.Sprintf(ERR_TYPE_MISMATCH+"%s %s %s", left, operator, right),
	}
}
