package parser

import (
	"fmt"
	"strconv"

	"github.com/nayyara-airlangga/basedlang/ast"
	"github.com/nayyara-airlangga/basedlang/lexer"
	"github.com/nayyara-airlangga/basedlang/token"
)

// Pratt parser function types
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)

type precedence int

const (
	_ precedence = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	l *lexer.Lexer

	curTok  token.Token
	peekTok token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Set curTok and peekTok
	p.nextToken()
	p.nextToken()

	// Register prefix functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntLiteral)
	p.registerPrefix(token.STRING, p.parseString)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	// Register infix functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) Errs() []string { return p.errors }

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curTok}

	// Expects an identifier after the let keyword
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curTok}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curTok}

	stmt.Expression = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(pr precedence) ast.Expression {
	prefixFn := p.prefixParseFns[p.curTok.Type]
	if prefixFn == nil {
		p.noPrefixParseFnErr(p.curTok.Type)
		return nil
	}

	leftExpr := prefixFn()

	for !p.peekTokenIs(token.SEMICOLON) && pr < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.peekTok.Type]
		if infixFn == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infixFn(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseIntLiteral() ast.Expression {
	lit := &ast.IntLiteral{Token: p.curTok}

	value, err := strconv.ParseInt(p.curTok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curTok, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curTok}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			expr.Else = p.parseIfExpression()
		} else if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			expr.Else = p.parseBlockStatement()
		} else {
			return nil
		}
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	b := &ast.BlockStatement{Token: p.curTok, Statements: []ast.Statement{}}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			b.Statements = append(b.Statements, stmt)
		}
		p.nextToken()
	}

	return b
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	f := &ast.FunctionLiteral{Token: p.curTok}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	f.Params = p.parseFunctionParameters()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	f.Body = p.parseBlockStatement()

	return f
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()

		if p.peekTokenIs(end) {
			p.nextToken()
			return list
		} else {
			p.nextToken()
			list = append(list, p.parseExpression(LOWEST))
		}
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	a := &ast.ArrayLiteral{Token: p.curTok}
	a.Elems = p.parseExpressionList(token.RBRACKET)
	return a
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		return params
	}

	p.nextToken()

	params = append(params, &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal})
	}

	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curTok, Function: function}
	call.Args = p.parseExpressionList(token.RPAREN)
	return call
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.curTok, Operator: p.curTok.Literal}

	p.nextToken()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{Token: p.curTok, Left: left, Operator: p.curTok.Literal}
	pre := p.curPrecedence()

	p.nextToken()

	expr.Right = p.parseExpression(pre)

	return expr
}

func (p *Parser) noPrefixParseFnErr(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function found for %s", t)
	p.errors = append(p.errors, msg)
}

func getPrecedence(t token.TokenType) precedence {
	switch t {
	case token.EQ, token.NEQ:
		return EQUALS
	case token.LT, token.GT, token.LTE, token.GTE:
		return LESSGREATER
	case token.PLUS, token.MINUS:
		return SUM
	case token.ASTERISK, token.SLASH:
		return PRODUCT
	case token.LPAREN:
		return CALL
	default:
		return LOWEST
	}
}

func (p *Parser) curPrecedence() precedence {
	return getPrecedence(p.curTok.Type)
}

func (p *Parser) peekPrecedence() precedence {
	return getPrecedence(p.peekTok.Type)
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curTok.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekTok.Type == t
}

func (p *Parser) peekErr(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekTok.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekErr(t)
		return false
	}
}
