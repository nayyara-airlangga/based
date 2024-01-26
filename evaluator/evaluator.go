package evaluator

import (
	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(n ast.Node) object.Object {
	switch n := n.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.BlockStatement:
		return evalStatements(n.Statements)
	// Expressions
	case *ast.IntLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToObjBool(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)
		return evalInfixExpression(n.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(n)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	cond := Eval(ie.Condition)

	if isTruthy(cond) {
		return Eval(ie.Body)
	} else if ie.Else != nil {
		switch el := ie.Else.(type) {
		case *ast.BlockStatement, *ast.IfExpression:
			return Eval(el)
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
	// The following cases are only for boolean expressions
	case op == "==":
		return nativeBoolToObjBool(left == right)
	case op == "!=":
		return nativeBoolToObjBool(left != right)
	default:
		return NULL
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
		return NULL
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
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
	return NULL
}

func evalStatements(stmts []ast.Statement) (res object.Object) {
	for _, s := range stmts {
		res = Eval(s)
	}
	return res
}

func nativeBoolToObjBool(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}
