package evaluator

import (
	"testing"

	"github.com/nayyara-airlangga/basedlang/lexer"
	"github.com/nayyara-airlangga/basedlang/object"
	"github.com/nayyara-airlangga/basedlang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func testEval(input string) object.Object {
	p := parser.New(lexer.New(input))
	program := p.Parse()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	intObj, isInt := obj.(*object.Integer)
	if !isInt {
		t.Errorf("obj is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if intObj.Value != expected {
		t.Errorf("obj.Value is incorrect. expected=%d, got=%d", expected, intObj.Value)
		return false
	}
	return true
}
