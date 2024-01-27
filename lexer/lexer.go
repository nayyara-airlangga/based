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

func (l *Lexer) peekCh() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
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

func (l *Lexer) readString() string {
	pos := l.position + 1

	for {
		l.readCh()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[pos:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdent(checkFn func(ch byte) bool) string {
	pos := l.position

	for checkFn(l.ch) {
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

func newStringToken(str string) token.Token {
	return token.Token{Type: token.STRING, Literal: str}
}

func newEOFToken() token.Token {
	return token.Token{Type: token.EOF, Literal: ""}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaces()

	switch l.ch {
	case '=':
		if l.peekCh() == '=' {
			ch := l.ch
			l.readCh()
			lit := string(ch) + string(l.ch)
			tok = newIdentToken(token.EQ, lit)
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekCh() == '=' {
			ch := l.ch
			l.readCh()
			lit := string(ch) + string(l.ch)
			tok = newIdentToken(token.NEQ, lit)
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		if l.peekCh() == '=' {
			ch := l.ch
			l.readCh()
			lit := string(ch) + string(l.ch)
			tok = newIdentToken(token.LTE, lit)
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekCh() == '=' {
			ch := l.ch
			l.readCh()
			lit := string(ch) + string(l.ch)
			tok = newIdentToken(token.GTE, lit)
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok = newStringToken(l.readString())
	case 0:
		tok = newEOFToken()
	default:
		if isLetter(l.ch) {
			ident := l.readIdent(isLetter)
			tok = newIdentToken(token.LookupType(ident), ident)
			return tok
		} else if isDigit(l.ch) {
			num := l.readIdent(isDigit)
			tok = newIdentToken(token.INT, num)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readCh()
	return tok
}
