package parser

import (
	"fmt"
	"testing"

	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/lexer"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testIdentifier(t, stmt.Expression, "foobar") {
		return
	}
}

func TestIntLiteralExpression(t *testing.T) {
	input := `
100;
202
	`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 2, len(program.Statements))
	}

	expectedExprs := []struct {
		value int64
	}{{100}, {202}}

	for i, ee := range expectedExprs {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not ast.ExpressionStatement. got=%T", i, program.Statements[i])
		}

		if !testIntegerLiteral(t, stmt.Expression, ee.value) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{"true;", true}, {"false", false}}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testBooleanLiteral(t, stmt.Expression, tc.expected) {
			return
		}
	}
}

func TestStringExpression(t *testing.T) {
	input := `"hello world";`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	str, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if str.Value != "hello world" {
		t.Errorf("incorrect str.Value. expected=%q, got=%q", "hello world", str.Value)
	}
}

func TestArrayExpression(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elems) != 3 {
		t.Fatalf("incorrect length of array. expected=%d, got=%d", 3, len(array.Elems))
	}
	testIntegerLiteral(t, array.Elems[0], 1)
	testInfixExpression(t, array.Elems[1], 2, "*", 2)
	testInfixExpression(t, array.Elems[2], 3, "+", 3)
}

func TestIndexExpression(t *testing.T) {
	input := "arr[1 + 1]"

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	idxExpr, isIdx := stmt.Expression.(*ast.IndexExpression)
	if !isIdx {
		t.Fatalf("stmt.Expression is not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, idxExpr.Left, "arr") {
		return
	}
	testInfixExpression(t, idxExpr.Index, 1, "+", 1)
}

func TestIfExpression(t *testing.T) {
	input := `
	if (x < y) {
		x
	}
	`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt := program.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, isIfExpr := exprStmt.Expression.(*ast.IfExpression)
	if !isIfExpr {
		t.Fatalf("exprStmt.Expression is not *ast.IfExpression. got=%T", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpr.Body.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in Body. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt = ifExpr.Body.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("Statements[0] is not an *ast.ExpressionStatement. got=%T", ifExpr.Body.Statements[0])
	}
	if !testIdentifier(t, exprStmt.Expression, "x") {
		return
	}

	if ifExpr.Else != nil {
		t.Fatalf("Else was not nil. got=%+v", ifExpr.Else)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `
	if (x < y) {
		x
	} else {
		y
	}
	`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt := program.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, isIfExpr := exprStmt.Expression.(*ast.IfExpression)
	if !isIfExpr {
		t.Fatalf("exprStmt.Expression is not *ast.IfExpression. got=%T", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpr.Body.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in Body. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt = ifExpr.Body.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("Statements[0] is not an *ast.ExpressionStatement. got=%T", ifExpr.Body.Statements[0])
	}
	if !testIdentifier(t, exprStmt.Expression, "x") {
		return
	}

	if ifExpr.Else == nil {
		t.Fatalf("Else block is nil")
	}

	bl, isBlock := ifExpr.Else.(*ast.BlockStatement)
	if !isBlock {
		t.Fatalf("Else is not an *ast.BlockStatement. got=%T", ifExpr.Else)
	}
	if len(bl.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in bl.Statements. expected=%d, got=%d", 1, len(bl.Statements))
	}

	exprStmt, isExprStmt = bl.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("bl.Statements[0] is not an *ast.ExpressionStatement. got=%T", bl.Statements[0])
	}
	if !testIdentifier(t, exprStmt.Expression, "y") {
		return
	}
}

func TestIfElseIfExpression(t *testing.T) {
	input := `
	if (x < y) {
		x
	} else if (x == y) {
		x - y
	} else {
		y
	}
	`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt := program.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, isIfExpr := exprStmt.Expression.(*ast.IfExpression)
	if !isIfExpr {
		t.Fatalf("exprStmt.Expression is not *ast.IfExpression. got=%T", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpr.Body.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in Body. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt = ifExpr.Body.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("Statements[0] is not an *ast.ExpressionStatement. got=%T", ifExpr.Body.Statements[0])
	}
	if !testIdentifier(t, exprStmt.Expression, "x") {
		return
	}

	if ifExpr.Else == nil {
		t.Fatalf("Else block is nil")
	}

	elif, isElifBlock := ifExpr.Else.(*ast.IfExpression)
	if !isElifBlock {
		t.Fatalf("Else is not an *ast.IfExpression. got=%T", ifExpr.Else)
	}
	if !testInfixExpression(t, elif.Condition, "x", "==", "y") {
		return
	}
	if len(elif.Body.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in elif.Body. expected=%d, got=%d", 1, len(elif.Body.Statements))
	}

	exprStmt, isExprStmt = elif.Body.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("elif.Body.Statements[0] is not an *ast.ExpressionStatement. got=%T", elif.Body.Statements[0])
	}
	if !testInfixExpression(t, exprStmt.Expression, "x", "-", "y") {
		return
	}

	if elif.Else == nil {
		t.Fatalf("elif.Else block is nil")
	}

	bl, isBlock := elif.Else.(*ast.BlockStatement)
	if !isBlock {
		t.Fatalf("Else is not an *ast.BlockStatement. got=%T", ifExpr.Else)
	}
	if len(bl.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in bl.Statements. expected=%d, got=%d", 1, len(bl.Statements))
	}

	exprStmt, isExprStmt = bl.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("bl.Statements[0] is not an *ast.ExpressionStatement. got=%T", bl.Statements[0])
	}
	if !testIdentifier(t, exprStmt.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt := program.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, isFunction := exprStmt.Expression.(*ast.FunctionLiteral)
	if !isFunction {
		t.Fatalf("exprStmt.Expression is not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}
	if len(function.Params) != 2 {
		t.Fatalf("Unexpected number of parameters. expected=%d, got=%d", 2, len(function.Params))
	}

	testLiteralExpression(t, function.Params[0], "x")
	testLiteralExpression(t, function.Params[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Unexpected number of statements in Body.Statements. expected=%d, got=%d", 1, len(function.Body.Statements))
	}

	bodyStmt, isExprStmt := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParams(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		program := p.Parse()

		checkParserErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		function := exprStmt.Expression.(*ast.FunctionLiteral)

		if len(function.Params) != len(tc.expected) {
			t.Fatalf("Unexpected number of parameters. expected=%d, got=%d", len(tc.expected), len(function.Params))
		}

		for i, param := range tc.expected {
			testLiteralExpression(t, function.Params[i], param)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
	}

	exprStmt, isExprStmt := program.Statements[0].(*ast.ExpressionStatement)
	if !isExprStmt {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	call, isCall := exprStmt.Expression.(*ast.CallExpression)
	if !isCall {
		t.Fatalf("exprStmt.Expression is not *ast.CallExpression. got=%T", exprStmt.Expression)
	}
	if !testIdentifier(t, call.Function, "add") {
		return
	}
	if len(call.Args) != 3 {
		t.Fatalf("Unexpected number of arguments. expected=%d, got=%d", 3, len(call.Args))
	}

	testIntegerLiteral(t, call.Args[0], 1)
	testInfixExpression(t, call.Args[1], 2, "*", 3)
	testInfixExpression(t, call.Args[2], 4, "+", 5)
}

func TestCallArgs(t *testing.T) {
	tests := []struct {
		input string
		ident string
		args  []string
	}{
		{"add();", "add", []string{}},
		{"add(1);", "add", []string{"1"}},
		{"add(1, 2 * 3, 4 + 5);", "add", []string{"1", "(2 * 3)", "(4 + 5)"}},
		{`
		add(
		  1, 
		  2 * 3,
		  4 + 5,
		);`, "add", []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		program := p.Parse()

		checkParserErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		call, isCall := exprStmt.Expression.(*ast.CallExpression)
		if !isCall {
			t.Fatalf("exptStmt.Expression is not ast.CallExpression. got=%T", exprStmt.Expression)
		}
		if !testIdentifier(t, call.Function, tc.ident) {
			return
		}
		if len(call.Args) != len(tc.args) {
			t.Fatalf("Unexpected number of arguments. expected=%d, got=%d", len(tc.args), len(call.Args))
		}

		for i, arg := range tc.args {
			if call.Args[i].String() != arg {
				t.Errorf("Unexpected value for argument %d. expected=%q, got=%q", i, arg, call.Args[i].String())
			}
		}
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, isIdent := expr.(*ast.Identifier)
	if !isIdent {
		t.Errorf("expr is not *ast.Identifier. got=%T", expr)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value is not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLit, ok := il.(*ast.IntLiteral)
	if !ok {
		t.Errorf("il not *ast.IntLiteral. got=%T", il)
		return false
	}
	if intLit.Value != value {
		t.Errorf("intLit.Value not %d. got=%d", value, intLit.Value)
		return false
	}
	if intLit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("intLit.TokenLiteral() not %d. got=%s", value, intLit.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, expected bool) bool {
	bo, isBool := expr.(*ast.BooleanLiteral)
	if !isBool {
		t.Errorf("expr is not *ast.BooleanLiteral. got=%T", expr)
		return false
	}
	if bo.Value != expected {
		t.Errorf("bo.Value is not %t. got=%t", expected, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", expected) {
		t.Errorf("bo.TokenLiteral() is not %t. got=%s", expected, bo.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected any) bool {
	switch eVal := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(eVal))
	case int64:
		return testIntegerLiteral(t, expr, eVal)
	case string:
		return testIdentifier(t, expr, eVal)
	case bool:
		return testBooleanLiteral(t, expr, eVal)
	default:
		t.Errorf("Unhandled type for expr. got=%T", expr)
		return false
	}
}

func testInfixExpression(t *testing.T, expr ast.Expression, left any, op string, right any) bool {
	opExpr, isOpExpr := expr.(*ast.InfixExpression)
	if !isOpExpr {
		t.Errorf("expr is not ast.InfixExpression, got=%T(%s)", expr, expr)
		return false
	}
	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}
	if opExpr.Operator != op {
		t.Errorf("opExpr.Operator is not '%s'. got=%s", op, opExpr.Operator)
		return false
	}
	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		val      any
	}{
		{"!120;", "!", 120},
		{"-67", "-", 67},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))

		program := p.Parse()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tc.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tc.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tc.val) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  any
		op       string
		rightVal any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"foobar <= barfoo;", "foobar", "<=", "barfoo"},
		{"foobar >= barfoo;", "foobar", ">=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))

		program := p.Parse()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		if !testInfixExpression(t, stmt.Expression, tc.leftVal, tc.op, tc.rightVal) {
			return
		}
	}
}

func TestOperatorPrecedences(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"true;",
			"true",
		},
		{
			"false",
			"false",
		},

		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"6 <= 6 == 6 < 7",
			"((6 <= 6) == (6 < 7))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			`add(
				a, 
				b, 
				1, 
				2 * 3, 
				4 + 5, 
				add(6, 7 * 8),
			);`,
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		program := p.Parse()

		checkParserErrors(t, p)

		if actual := program.String(); actual != tc.expected {
			t.Errorf("expected=%q, got=%q", tc.expected, actual)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		ident string
		value any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		program := p.Parse()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tc.ident) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tc.value) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		val   any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tc := range tests {
		p := New(lexer.New(tc.input))
		program := p.Parse()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 1, len(program.Statements))
		}

		returnStmt, isReturn := program.Statements[0].(*ast.ReturnStatement)
		if !isReturn {
			t.Fatalf("program.Statements[0] is not *ast.ReturnStatement. got=%T", program.Statements[0])
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
		if !testLiteralExpression(t, returnStmt.ReturnValue, tc.val) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errs()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
