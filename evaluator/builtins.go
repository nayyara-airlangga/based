package evaluator

import "github.com/nayyara-airlangga/basedlang/object"

var (
	ErrInvalidLen = "invalid argument: %s (%s) not supported for len"
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
			default:
				return newError(ErrInvalidLen, arg.Inspect(), arg.Type())
			}
		},
	},
}
