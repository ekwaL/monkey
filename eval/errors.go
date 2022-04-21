package eval

import (
	"fmt"
	"monkey/object"
)

const (
	ERR_UNKNOWN_OPERATOR      = "unknown operator: "
	ERR_TYPE_MISMATCH         = "type mismatch: "
	ERR_IDENTIFIER_NOT_FOUND  = "identifier not found: "
	ERR_NOT_A_FUNCTION        = "not a function: "
	ERR_WRONG_ARGUMENTS_COUNT = "wrong arguments count: "
)

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

func infixTypeMismatchError(
	left object.ObjectType,
	operator string,
	right object.ObjectType) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_TYPE_MISMATCH+"%s %s %s", left, operator, right),
	}
}

func builtinTypeMismatchError(
	name string,
	args ...object.Object) *object.Error {
	argsStr := ""
	switch len(args) {
	case 0:
		argsStr = ""
	case 1:
		argsStr = string(args[0].Type())
	default:
		argsStr = string(args[0].Type())
		for _, arg := range args[1:] {
			argsStr += ", " + string(arg.Type())
		}
	}
	return &object.Error{
		Message: ERR_TYPE_MISMATCH + name + "(" + argsStr + ")",
	}
}

func identifierNotFoundError(identifier string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_IDENTIFIER_NOT_FOUND+"'%s'", identifier),
	}
}

func notAFunctionError(objType string, identifier string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_NOT_A_FUNCTION+"%s '%s'", objType, identifier),
	}
}

func wrongArgumentsCountError(expect int, got int) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_WRONG_ARGUMENTS_COUNT+"expect %d, got %d", expect, got),
	}
}
