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
	// Expressions
	case *ast.IntLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return nativeBoolToObjBool(n.Value)
	default:
		return nil
	}
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
