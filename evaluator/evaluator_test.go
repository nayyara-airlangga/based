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
			`"Hello" + 5`,
			"type mismatch: STRING + INTEGER",
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
			`"Hello" - "World"`,
			"unsupported operator: STRING - STRING",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`999[1]`,
			"unsupported operator: index not supported on 999 (INTEGER)",
		},
		{
			`[1, 2, 3][true]`,
			"invalid argument: index true (BOOLEAN) is not an integer",
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

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "invalid argument: 1 (INTEGER) not supported for len"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)
		switch expected := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + -4, true];"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elems) != 4 {
		t.Fatalf("incorrect number of elements. expected=%d, got=%d", 4, len(result.Elems))
	}
	testIntegerObject(t, result.Elems[0], 1)
	testIntegerObject(t, result.Elems[1], 4)
	testIntegerObject(t, result.Elems[2], -1)
	testBooleanObject(t, result.Elems[3], true)
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, isFunc := evaluated.(*object.Function)
	if !isFunc {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Params) != 1 {
		t.Fatalf("incorrect number of params. expected=%d, got=%d (%+v)", 1, len(fn.Params), fn.Params)
	}
	if fn.Params[0].String() != "x" {
		t.Fatalf("incorrect parameter. expected='%s' got=%q", "x", fn.Params[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("incorrect function body. expected=%q got=%q", expectedBody, fn.Body.String())
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"[1, 2, !false][1 + 1];",
			true,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil},
		{
			"[1, 2, 3][-1]",
			3,
		},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.input)

		switch val := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(val))
		case bool:
			testBooleanObject(t, evaluated, val)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"let add = fn(x, y) { x + y; }; add(5);", "wrong number of arguments. got=1, want=2"},
		{"fn(x) { x * 3; }(5)", 15},
		{`
		let muller = fn(x) { fn(y) { x * y } };
		let fiveMul = muller(5)
		fiveMul(3)
		`, 15},
	}
	for _, tc := range tests {
		evaluated := testEval(tc.input)
		switch expected := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
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

func TestEvalString(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, isStr := evaluated.(*object.String)
	if !isStr {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("incorrect String value. expected=%q, got=%q", "Hello World!", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("incorrect String value. expected=%q, got=%q", "Hello World!", str.Value)
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
