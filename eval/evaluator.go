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

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.BlockStmt:
		return evalStatements(node.Statements)
	case *ast.ExpressionStmt:
		return Eval(node.Expression)
	case *ast.IntLiteralExpr:
		return &object.Integer{Value: node.Value}
	case *ast.BoolLiteralExpr:
		return boolToBooleanObject(node.Value)
	case *ast.PrefixExpr:
		right := Eval(node.Right)
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpr(left, node.Operator, right)
	case *ast.IfExpr:
		return evalIfExpr(node)
	default:
		println(fmt.Sprintf("eval is unimplemented for %T", node))
		return nil
	}
}

func evalStatements(statements []ast.Statement) (result object.Object) {
	for _, stmt := range statements {
		result = Eval(stmt)
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
		return NULL // TODO: runtime error
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
		return NULL // TODO: runtime error
	}

	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfixExpr(left object.Object, operator string, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpr(left, operator, right)
	}

	switch operator {
	case token.EQUAL_EQUAL:
		return boolToBooleanObject(left == right)

	case token.NOT_EQUAL:
		return boolToBooleanObject(left != right)
	default:
		return NULL // TODO: runtime error
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
		return NULL
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
