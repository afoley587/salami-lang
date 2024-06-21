package ast

import (
	"github.com/afoley/salami-lang/tok"
)

type Node interface {
	Literal() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	}
	return ""
}

type VarStatement struct {
	Token tok.Tok
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()  {}
func (vs *VarStatement) Literal() string { return vs.Token.Literal }

type Identifier struct {
	Token tok.Tok // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Literal() string { return i.Token.Literal }

type IntegerLiteral struct {
	Token tok.Tok // The token.INT token
	Value int64   // The actual value of the integer
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) Literal() string { return il.Token.Literal }

type InfixExpression struct {
	Token    tok.Tok // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) Literal() string { return ie.Token.Literal }

type IfExpression struct {
	Token       tok.Tok // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) statementNode()  {}
func (ie *IfExpression) Literal() string { return ie.Token.Literal }

type BlockStatement struct {
	Token      tok.Tok // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()  {}
func (bs *BlockStatement) Literal() string { return bs.Token.Literal }

type BooleanLiteral struct {
	Token tok.Tok
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) Literal() string { return bl.Token.Literal }

type ExitStatement struct {
	Token tok.Tok // The 'exit' token
	Value Expression
}

func (es *ExitStatement) statementNode()  {}
func (es *ExitStatement) Literal() string { return es.Token.Literal }
