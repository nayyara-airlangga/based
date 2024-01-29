package evaluator

import (
	"fmt"

	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/object"
)

const (
	ErrUnsupportedOperatorInfix  = "unsupported operator: %s %s %s"
	ErrUnsupportedOperatorPrefix = "unsupported operator: %s%s"
	ErrTypeMismatch              = "type mismatch: %s %s %s"
	ErrIdentifierNotFound        = "identifier not found: %s"
	ErrNotAFunction              = "not a function: %s"
	ErrWrongNumberOfArgs         = "wrong number of arguments. got=%d, want=%d"
)

func newError(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(n ast.Node, env *object.Environment) object.Object {
	switch n := n.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(n.Statements, env)
	case *ast.LetStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(n.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
		// Expressions
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.IntLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToObjBool(n.Value)
	case *ast.StringLiteral:
		return &object.String{Value: n.Value}
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.FunctionLiteral:
		return &object.Function{Params: n.Params, Body: n.Body, Env: env}
	case *ast.CallExpression:
		f := Eval(n.Function, env)
		if isError(f) {
			return f
		}
		args := evalExpressions(n.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(f, args)
	default:
		return NULL
	}

	return nil
}

func applyFunction(f object.Object, args []object.Object) object.Object {
	switch fn := f.(type) {
	case *object.Function:
		fun, isFunc := f.(*object.Function)
		if !isFunc {
			return newError(ErrNotAFunction, fun.Type())
		}
		if len(fn.Params) != len(args) {
			return newError(ErrWrongNumberOfArgs, len(args), len(fn.Params))
		}
		extEnv := extendFunctionEnv(fun, args)
		evaluated := Eval(fun.Body, extEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError(ErrNotAFunction, f.Type())
	}

}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewLocalEnvironment(fn.Env)
	for i, arg := range fn.Params {
		env.Set(arg.Value, args[i])
	}
	return env
}
func unwrapReturnValue(obj object.Object) object.Object {
	if rv, isRetVal := obj.(*object.ReturnValue); isRetVal {
		return rv.Value
	}
	return obj
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) (result []object.Object) {
	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	val, exists := env.Get(id.Value)
	if exists {
		return val
	}

	builtin, exists := builtins[id.Value]
	if exists {
		return builtin
	}

	return newError(ErrIdentifierNotFound, id.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)

	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Body, env)
	} else if ie.Else != nil {
		switch el := ie.Else.(type) {
		case *ast.BlockStatement, *ast.IfExpression:
			return Eval(el, env)
		default:
			return NULL
		}
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(op, left, right)
	// The following cases are only for boolean expressions
	case op == "==":
		return nativeBoolToObjBool(left == right)
	case op == "!=":
		return nativeBoolToObjBool(left != right)
	case left.Type() != right.Type():
		return newError(ErrTypeMismatch, left.Type(), op, right.Type())
	default:
		return newError(ErrUnsupportedOperatorInfix, left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	leftInt := left.(*object.Integer)
	rightInt := right.(*object.Integer)

	switch op {
	// Arithmetics
	case "*":
		return &object.Integer{Value: leftInt.Value * rightInt.Value}
	case "/":
		return &object.Integer{Value: leftInt.Value / rightInt.Value}
	case "+":
		return &object.Integer{Value: leftInt.Value + rightInt.Value}
	case "-":
		return &object.Integer{Value: leftInt.Value - rightInt.Value}
	// Relational
	case "<":
		return nativeBoolToObjBool(leftInt.Value < rightInt.Value)
	case "<=":
		return nativeBoolToObjBool(leftInt.Value <= rightInt.Value)
	case ">":
		return nativeBoolToObjBool(leftInt.Value > rightInt.Value)
	case ">=":
		return nativeBoolToObjBool(leftInt.Value >= rightInt.Value)
	case "==":
		return nativeBoolToObjBool(leftInt.Value == rightInt.Value)
	case "!=":
		return nativeBoolToObjBool(leftInt.Value != rightInt.Value)
	default:
		return newError(ErrUnsupportedOperatorInfix, left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right object.Object) object.Object {
	if op != "+" {
		return newError(ErrUnsupportedOperatorInfix, left.Type(), op, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError(ErrUnsupportedOperatorPrefix, op, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Boolean:
		return nativeBoolToObjBool(!right.Value)
	case *object.Integer:
		if right.Value == 0 {
			return TRUE
		}
		return FALSE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if intObj, isInt := right.(*object.Integer); isInt {
		intObj.Value = -intObj.Value
		return intObj
	}
	return newError(ErrUnsupportedOperatorPrefix, "-", right.Type())
}

func evalProgram(stmts []ast.Statement, env *object.Environment) (res object.Object) {
	for _, s := range stmts {
		res = Eval(s, env)

		if err, isErr := res.(*object.Error); isErr {
			return err
		}
		if rv, isRetVal := res.(*object.ReturnValue); isRetVal {
			return rv.Value
		}
	}
	return res
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) (res object.Object) {
	for _, s := range stmts {
		res = Eval(s, env)

		if err, isErr := res.(*object.Error); isErr {
			return err
		}

		if rv, isRetVal := res.(*object.ReturnValue); isRetVal {
			return rv
		}
	}
	return res
}

func nativeBoolToObjBool(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}
