package lexer

import (
	"testing"

	"github.com/nayyara-airlangga/basedlang/token"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"

	expectedTokens := []struct {
		tokenType token.TokenType
		literal   string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, et := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != et.tokenType {
			t.Fatalf("expectedTokens[%d] - wrong token type. expected=%q, got=%q", i, et.tokenType, tok.Type)
		}

		if tok.Literal != et.literal {
			t.Fatalf("expectedTokens[%d] - literal wrong. expected=%q, got=%q", i, et.literal, tok.Literal)
		}
	}
}
