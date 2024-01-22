package lexer

import (
	"github.com/nayyara-airlangga/basedlang/token"
)

type Lexer struct {
	input        string
	position     int  // current position
	nextPosition int  // position after current
	ch           byte // current char being read
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readCh()
	return l
}

func (l *Lexer) readCh() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}

	l.position = l.nextPosition
	l.nextPosition++
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdent() string {
	pos := l.position

	for isLetter(l.ch) {
		l.readCh()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readCh()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newIdentToken(tokenType token.TokenType, ident string) token.Token {
	return token.Token{Type: tokenType, Literal: ident}
}

func newEOFToken() token.Token {
	return token.Token{Type: token.EOF, Literal: ""}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaces()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok = newEOFToken()
	default:
		if isLetter(l.ch) {
			ident := l.readIdent()
			tok = newIdentToken(token.LookupType(ident), ident)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readCh()
	return tok
}
