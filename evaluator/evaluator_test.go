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
	}{{"5", 5}, {"101", 101}}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{"true", true}, {"false", false}}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolObj, isBool := obj.(*object.Boolean)
	if !isBool {
		t.Errorf("obj is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if boolObj.Value != expected {
		t.Errorf("obj.Value is incorrect. expected=%t, got=%t", expected, boolObj.Value)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	p := parser.New(lexer.New(input))
	program := p.Parse()

	return Eval(program)
}
