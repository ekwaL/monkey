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

func typeMismatchError(
	left object.ObjectType,
	operator string,
	right object.ObjectType) *object.Error {

	return &object.Error{
		Message: fmt.Sprintf(ERR_TYPE_MISMATCH+"%s %s %s", left, operator, right),
	}
}

func identifierNotFound(identifier string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_IDENTIFIER_NOT_FOUND+"'%s'", identifier),
	}
}

func notAFunction(objType string, identifier string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_NOT_A_FUNCTION+"%s '%s'", objType, identifier),
	}
}

func wrongArgumentsCount(expect int, got int) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_WRONG_ARGUMENTS_COUNT+"expect %d, got %d", expect, got),
	}
}
