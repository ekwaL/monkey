package eval

import (
	"fmt"
	"monkey/object"
	"os"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return wrongArgumentsCountError(1, len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return builtinTypeMismatchError("len", args...)
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			fmt.Fprint(os.Stdout, stringifyArgs(args))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			fmt.Fprintln(os.Stdout, stringifyArgs(args))
			return NULL
		},
	},
}

func stringifyArgs(args []object.Object) string {
	result := []string{}
	for _, arg := range args {
		result = append(result, arg.Inspect())
	}
	return strings.Join(result, ", ")
}
