package evaluator

import (
	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/object"
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
