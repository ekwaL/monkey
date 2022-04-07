package eval

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStmt:
		return Eval(node.Expression)
	case *ast.IntLiteralExpr:
		return &object.Integer{Value: node.Value}
	case *ast.BoolLiteralExpr:
		return &object.Boolean{Value: node.Value}
	default:
		println("unimplemented")
	}

	return nil
}

func evalStatements(statements []ast.Statement) (result object.Object) {
	for _, stmt := range statements {
		result = Eval(stmt)
	}
	return
}
