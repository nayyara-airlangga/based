package parser

import (
	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/lexer"
	"github.com/nayyara-airlangga/basedlang/token"
)

type Parser struct {
	l *lexer.Lexer

	curTok  token.Token
	peekTok token.Token
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Set curTok and peekTok
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() *ast.Program {
	return nil
}
