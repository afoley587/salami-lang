package interpreter

import (
	"fmt"

	"github.com/afoley/salami-lang/ast"
)

type Environment struct {
	store map[string]interface{}
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]interface{})}
}

func (e *Environment) Get(name string) (interface{}, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

func (e *Environment) Set(name string, value interface{}) interface{} {
	e.store[name] = value
	return value
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fn *Function) Literal() string { return "gorlami" }

type ReturnValue struct {
	Value interface{}
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
	if i.Exited {
		return i.ExitCode
	}
	switch node := node.(type) {
	case *ast.Program:
		return i.evalProgram(node)
	case *ast.VarStatement:
		return i.evalVarStatement(node)
	case *ast.FunctionStatement:
		return i.evalFunctionStatement(node)
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
	case *ast.FunctionLiteral:
		return i.evalFunctionLiteral(node)
	case *ast.CallExpression:
		return i.evalCallExpression(node)
	case *ast.ReturnStatement:
		return i.evalReturnStatement(node)
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
		i.env.Set(stmt.Name.Value, val)
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
	case ">":
		return left > right
	case "<":
		return left < right

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

		if i.Exited {
			return result
		}
	}

	return result
}

func (i *Interpreter) evalFunctionLiteral(fl *ast.FunctionLiteral) interface{} {
	params := fl.Parameters
	body := fl.Body
	env := i.env

	fmt.Println("parsed fn literal")

	return &Function{Parameters: params, Body: body, Env: env}
}

func (i *Interpreter) evalCallExpression(ce *ast.CallExpression) interface{} {
	function := i.Interpret(ce.Function)

	if function == nil {
		return nil
	}

	fn, ok := function.(*Function)
	if !ok {
		return nil
	}

	args := []interface{}{}
	for _, arg := range ce.Arguments {
		args = append(args, i.Interpret(arg))
	}

	extendedEnv := extendFunctionEnv(fn, args)
	evaluated := i.evalBlockStatementWithEnv(fn.Body, extendedEnv)
	return evaluated
}

func (i *Interpreter) evalExitStatement(stmt *ast.ExitStatement) interface{} {
	val := i.Interpret(stmt.Value)
	if val != nil {
		i.ExitCode = val.(int64)
		i.Exited = true
	}
	return val
}

func (i *Interpreter) evalFunctionStatement(stmt *ast.FunctionStatement) interface{} {
	fn := &Function{
		Parameters: stmt.Parameters,
		Body:       stmt.Body,
		Env:        i.env,
	}

	i.env.Set(stmt.Name.Value, fn)
	return fn
}

func (i *Interpreter) evalReturnStatement(rs *ast.ReturnStatement) interface{} {
	value := i.Interpret(rs.ReturnValue)
	return &ReturnValue{Value: value}
}

func (i *Interpreter) evalBlockStatementWithEnv(block *ast.BlockStatement, env *Environment) interface{} {
	previousEnv := i.env
	i.env = env

	var result interface{}
	for _, stmt := range block.Statements {
		result = i.Interpret(stmt)
		if returnValue, ok := result.(*ReturnValue); ok {
			i.env = previousEnv
			return returnValue.Value
		}
		if i.Exited {
			break
		}
	}

	i.env = previousEnv
	return result
}

func extendFunctionEnv(fn *Function, args []interface{}) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, (args[paramIdx]).(int64))
	}

	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
