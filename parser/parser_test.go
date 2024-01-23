package parser

import (
	"testing"

	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	p := New(lexer.New(input))
	program := p.Parse()

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
