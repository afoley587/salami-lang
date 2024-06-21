package parser

import (
	"fmt"
	"strconv"

	"github.com/afoley/salami-lang/ast"
	"github.com/afoley/salami-lang/lexer"
	"github.com/afoley/salami-lang/tok"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tok.Tok
	peekToken    tok.Tok
	errors       []string

	prefixParseFns map[tok.TokenType]prefixParseFn
	infixParseFns  map[tok.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          lexer,
		prefixParseFns: make(map[tok.TokenType]prefixParseFn),
		infixParseFns:  make(map[tok.TokenType]infixParseFn),
	}
	p.nextToken()
	p.nextToken()

	// Register prefix parse functions
	p.registerPrefix(tok.INT, p.parseIntegerLiteral)
	p.registerPrefix(tok.IDENT, p.parseIdentifier)
	p.registerPrefix(tok.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tok.FALSE, p.parseBooleanLiteral)

	// Register infix parse functions
	p.registerInfix(tok.PLUS, p.parseInfixExpression)
	p.registerInfix(tok.MINUS, p.parseInfixExpression)
	p.registerInfix(tok.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tok.SLASH, p.parseInfixExpression)
	p.registerInfix(tok.GT, p.parseInfixExpression)
	p.registerInfix(tok.LT, p.parseInfixExpression)

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) peekTokenIs(t tok.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t tok.TokenType) bool {

	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t tok.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != tok.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case tok.VAR:
		return p.parseVarStatement()
	case tok.IF:
		return p.parseIfExpression().(ast.Statement)
	case tok.EXIT:
		return p.parseExitStatement()
	default:
		return nil
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.currentToken}

	if !p.expectPeek(tok.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(tok.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tok.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

const (
	_ int = iota
	LOWEST
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
	COMPARE // > or <
)

var precedences = map[tok.TokenType]int{
	tok.PLUS:     SUM,
	tok.MINUS:    SUM,
	tok.ASTERISK: PRODUCT,
	tok.SLASH:    PRODUCT,
	tok.GT:       COMPARE,
	tok.LT:       COMPARE,
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(tok.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(tok.LPAREN) {
		fmt.Println(string(p.currentToken.Literal))
		fmt.Println("Expected '(' after 'if'")
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if expression.Condition == nil {
		fmt.Println("Failed to parse condition")
		return nil
	}

	if !p.expectPeek(tok.RPAREN) {
		fmt.Println("Expected ')' after condition")
		return nil
	}

	if !p.expectPeek(tok.LBRACE) {
		fmt.Println("Expected '{' after ')'")
		return nil
	}

	expression.Consequence = p.parseBlockStatement()
	if expression.Consequence == nil {
		return nil
	}

	if p.peekTokenIs(tok.ELSE) {
		p.nextToken()

		if !p.expectPeek(tok.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
		if expression.Alternative == nil {
			return nil
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !(p.currentToken.Type == tok.RBRACE) && !(p.currentToken.Type == tok.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExitStatement() *ast.ExitStatement {
	stmt := &ast.ExitStatement{Token: p.currentToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tok.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == tok.TRUE}
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType tok.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tok.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
