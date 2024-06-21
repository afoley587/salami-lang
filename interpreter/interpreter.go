package interpreter

import (
	"github.com/afoley/salami-lang/ast"
)

type Environment struct {
	store map[string]int64
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]int64)}
}

func (e *Environment) Get(name string) (int64, bool) {
	value, ok := e.store[name]
	return value, ok
}

func (e *Environment) Set(name string, value int64) int64 {
	e.store[name] = value
	return value
}

type Interpreter struct {
	env      *Environment
	ExitCode int64
	Exited   bool
}

func New() *Interpreter {
	return &Interpreter{env: NewEnvironment()}
}

func (i *Interpreter) Interpret(node ast.Node) interface{} {
	switch node := node.(type) {
	case *ast.Program:
		return i.evalProgram(node)
	case *ast.VarStatement:
		return i.evalVarStatement(node)
	case *ast.Identifier:
		return i.evalIdentifier(node)
	case *ast.IntegerLiteral:
		return node.Value
	case *ast.BooleanLiteral:
		return node.Value
	case *ast.IfExpression:
		return i.evalIfExpression(node)
	case *ast.BlockStatement:
		return i.evalBlockStatement(node)
	case *ast.InfixExpression:
		return i.evalInfixExpression(node)
	case *ast.ExitStatement:
		return i.evalExitStatement(node)
	default:
		return nil
	}
}

func (i *Interpreter) evalProgram(program *ast.Program) interface{} {
	var result interface{}
	for _, stmt := range program.Statements {
		result = i.Interpret(stmt)
	}
	return result
}

func (i *Interpreter) evalVarStatement(stmt *ast.VarStatement) interface{} {
	val := i.Interpret(stmt.Value)
	if val != nil {
		i.env.Set(stmt.Name.Value, val.(int64))
	}
	return val
}

func (i *Interpreter) evalIdentifier(node *ast.Identifier) interface{} {
	if val, ok := i.env.Get(node.Value); ok {
		return val
	}
	return nil
}

func (i *Interpreter) evalInfixExpression(node *ast.InfixExpression) interface{} {
	left := i.Interpret(node.Left).(int64)
	right := i.Interpret(node.Right).(int64)

	switch node.Operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	default:
		return nil
	}
}

func (i *Interpreter) evalIfExpression(node *ast.IfExpression) interface{} {
	condition := i.Interpret(node.Condition).(bool)

	if condition {
		return i.Interpret(node.Consequence)
	} else if node.Alternative != nil {
		return i.Interpret(node.Alternative)
	} else {
		return nil
	}
}

func (i *Interpreter) evalBlockStatement(block *ast.BlockStatement) interface{} {
	var result interface{}

	for _, stmt := range block.Statements {
		result = i.Interpret(stmt)
	}

	return result
}

func (i *Interpreter) evalExitStatement(stmt *ast.ExitStatement) interface{} {
	val := i.Interpret(stmt.Value)
	if val != nil {
		i.ExitCode = val.(int64)
		i.Exited = true
	}
	return val
}
