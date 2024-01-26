package evaluator

import (
	"testing"

	"github.com/nayyara-airlangga/basedlang/lexer"
	"github.com/nayyara-airlangga/basedlang/object"
	"github.com/nayyara-airlangga/basedlang/parser"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unsupported operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unsupported operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unsupported operator: BOOLEAN + BOOLEAN",
		}, {
			"if (10 > 1) { true + false; }",
			"unsupported operator: BOOLEAN + BOOLEAN",
		},
		{
			`
				if (10 > 1) {
				if (10 > 1) {
				return true + false;
				}
				return 1;
				}
				`,
			"unsupported operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		err, isErr := evaluated.(*object.Error)
		if !isErr {
			t.Errorf("no error returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if err.Message != tc.expected {
			t.Errorf("wrong error message. expected=%s, got=%s", tc.expected, err.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"101", 101},
		{"-5", -5},
		{"-101", -101},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{`
		if (1 < 1) { 
		    10 
		} else if (1 <= 1) { 
			30 
		} else { 
			20 
		}`, 30},
		{`
		if (1 < 1) { 
			10 
		} else if (1 <= 0) { 
			30
		} else { 
			20 
		}`, 20},
		{`
		if (1 < 1) { 
			10 
		} else if (1 <= 0) { 
			30
		}`, nil},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		expected, isInt := tc.expected.(int)

		if isInt {
			testIntegerObject(t, evaluated, int64(expected))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; return 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
			    return 10;
			}
			return 1;
		}`, 10},
		{`
		if (10 > 1) {
			if (10 > 11) {
			    return 10;
			}
			
			return 1;
		}`, 1},
		{`
		if (10 > 1) {
			if (10 > 11) {
			    return 10;
			}
		}`, nil},
		{`
		if (10 > 1) {
			if (10 > 11) {
			    return 10;
			}
		}
		return 9;
		`, 9},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		expected, isInt := tc.expected.(int)

		if isInt {
			testIntegerObject(t, evaluated, int64(expected))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!0", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("obj is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
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
	env := object.NewEnvironment()

	return Eval(program, env)
}
