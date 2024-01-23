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

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() not %s. got=%s", "foobar", ident.TokenLiteral())
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

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		intVal   int64
	}{
		{"!120;", "!", 120},
		{"-67", "-", 67},
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
		if !testIntegerLiteral(t, exp.Right, tc.intVal) {
			return
		}
	}
}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("Parse() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 3, len(program.Statements))
	}

	expectedStmts := []struct {
		ident string
	}{{"x"}, {"y"}, {"foobar"}}

	for i, es := range expectedStmts {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, es.ident) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	p := New(lexer.New(input))
	program := p.Parse()

	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Unexpected number of statements. expected=%d, got=%d", 3, len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
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
