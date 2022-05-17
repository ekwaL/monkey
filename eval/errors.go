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
	ERR_WRONG_GET_TARGET      = "only instances have properties: "
	ERR_WRONG_SET_TARGET      = "only instances have fields: "
	ERR_UNDEFINED_PROP        = "undefined property: "
	ERR_SUPERCLASS_NOT_CLASS  = "superclass must be a class: "
	ERR_INTERNAL              = "internal error: "
	ERR_OUT_OF_BOUNDS         = "out of bounds: "
	ERR_NOT_HASHABLE_KEY      = "unusable as hash key: "
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

func indexOperatorError(left object.ObjectType, right object.ObjectType) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_UNKNOWN_OPERATOR+"%s[%s]", left, right),
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
	return &object.Error{Message: ERR_TYPE_MISMATCH + name + "(" + argsStr + ")"}
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

func wrongGetTargetError(target object.ObjectType, prop string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_WRONG_GET_TARGET+"%s.%s", string(target), prop),
	}
}

func wrongSetTargetError(target object.ObjectType, prop string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_WRONG_SET_TARGET+"%s.%s", string(target), prop),
	}
}

func undefinedPropertyError(prop string) *object.Error {
	return &object.Error{Message: fmt.Sprintf(ERR_UNDEFINED_PROP+"'%s'", prop)}
}

func superclassMustBeClassError(class string, super object.ObjectType) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_SUPERCLASS_NOT_CLASS+"'%s < %s'", class, super),
	}
}

func outOfBoundsError(left object.ObjectType, idx int64) *object.Error {
	return &object.Error{Message: fmt.Sprintf(ERR_OUT_OF_BOUNDS+"%s[%d]", left, idx)}
}

func notHashableKeyError(key object.ObjectType) *object.Error {
	return &object.Error{Message: fmt.Sprintf(ERR_NOT_HASHABLE_KEY+"%s", key)}
}

func internalResolveError(keyword string) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(ERR_INTERNAL+"looks like expression '%s' is not resolved correctly", keyword),
	}
}
