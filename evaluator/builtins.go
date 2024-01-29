package evaluator

import "github.com/nayyara-airlangga/basedlang/object"

const (
	ErrInvalidLen                  = "invalid argument: %s (%s) not supported for len"
	ErrNotEnoughArgsAppend         = "invalid argument: not enough arguments for append, expected>=1, got=0"
	ErrFirstArgShouldBeArrayAppend = "invalid argument: first argument for append must be an array. got=%s (%s)"
)

var builtins map[string]*object.Builtin = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(ErrWrongNumberOfArgs, len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elems))}
			default:
				return newError(ErrInvalidLen, arg.Inspect(), arg.Type())
			}
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError(ErrNotEnoughArgsAppend)
			}

			arr, isArr := args[0].(*object.Array)
			if !isArr {
				return newError(ErrFirstArgShouldBeArrayAppend, args[0].Inspect(), args[0].Type())
			}
			if len(args) == 1 {
				return arr
			}

			newArr := &object.Array{Elems: append(arr.Elems, args[1:]...)}

			return newArr
		},
	},
}
