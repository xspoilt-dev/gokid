package evaluator

import (
	"fmt"
	"gokid/parser"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// Eval evaluates AST nodes and returns objects
func Eval(node parser.Node, env *Environment) Object {
	switch node := node.(type) {

	// Statements
	case *parser.Program:
		return evalProgram(node.Statements, env)

	case *parser.ExpressionStatement:
		return Eval(node.Expression, env)

	case *parser.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *parser.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *parser.VarStatement:
		var val Object = NULL
		if node.Value != nil {
			val = Eval(node.Value, env)
			if isError(val) {
				return val
			}
		}
		env.Set(node.Name.Value, val)
		return val

	case *parser.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *parser.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions
	case *parser.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *parser.FloatLiteral:
		return &Float{Value: node.Value}

	case *parser.BooleanLiteral:
		return nativeBoolToPyMonkeyBool(node.Value)

	case *parser.StringLiteral:
		return &String{Value: node.Value}

	case *parser.NullLiteral:
		return NULL

	case *parser.Identifier:
		return evalIdentifier(node, env)

	case *parser.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *parser.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *parser.IfExpression:
		return evalIfExpression(node, env)

	case *parser.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}

	case *parser.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *parser.ObjectLiteral:
		return evalObjectLiteral(node, env)

	case *parser.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *parser.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}

	case *parser.WhileStatement:
		return evalWhileStatement(node, env)

	case *parser.ForStatement:
		return evalForStatement(node, env)

	case *parser.BreakStatement:
		return &Break{}

	case *parser.ContinueStatement:
		return &Continue{}

	case *parser.AssignmentExpression:
		return evalAssignmentExpression(node, env)

	default:
		return newError("unknown node type: %T", node)
	}
}

// Helper functions
func evalProgram(stmts []parser.Statement, env *Environment) Object {
	var result Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *parser.BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ || rt == BREAK_OBJ || rt == CONTINUE_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToPyMonkeyBool(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	switch right := right.(type) {
	case *Integer:
		return &Integer{Value: -right.Value}
	case *Float:
		return &Float{Value: -right.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == FLOAT_OBJ || right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToPyMonkeyBool(left == right)
	case operator == "!=":
		return nativeBoolToPyMonkeyBool(left != right)
	case operator == "&&":
		return evalLogicalAndExpression(left, right)
	case operator == "||":
		return evalLogicalOrExpression(left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToPyMonkeyBool(leftVal < rightVal)
	case ">":
		return nativeBoolToPyMonkeyBool(leftVal > rightVal)
	case "==":
		return nativeBoolToPyMonkeyBool(leftVal == rightVal)
	case "!=":
		return nativeBoolToPyMonkeyBool(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalFloatInfixExpression(operator string, left, right Object) Object {
	var leftVal, rightVal float64

	if left.Type() == FLOAT_OBJ {
		leftVal = left.(*Float).Value
	} else {
		leftVal = float64(left.(*Integer).Value)
	}

	if right.Type() == FLOAT_OBJ {
		rightVal = right.(*Float).Value
	} else {
		rightVal = float64(right.(*Integer).Value)
	}

	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToPyMonkeyBool(leftVal < rightVal)
	case ">":
		return nativeBoolToPyMonkeyBool(leftVal > rightVal)
	case "==":
		return nativeBoolToPyMonkeyBool(leftVal == rightVal)
	case "!=":
		return nativeBoolToPyMonkeyBool(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToPyMonkeyBool(leftVal == rightVal)
	case "!=":
		return nativeBoolToPyMonkeyBool(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalBooleanInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Boolean).Value
	rightVal := right.(*Boolean).Value

	switch operator {
	case "==":
		return nativeBoolToPyMonkeyBool(leftVal == rightVal)
	case "!=":
		return nativeBoolToPyMonkeyBool(leftVal != rightVal)
	case "&&":
		return nativeBoolToPyMonkeyBool(leftVal && rightVal)
	case "||":
		return nativeBoolToPyMonkeyBool(leftVal || rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalLogicalAndExpression(left, right Object) Object {
	if !isTruthy(left) {
		return FALSE
	}
	return nativeBoolToPyMonkeyBool(isTruthy(right))
}

func evalLogicalOrExpression(left, right Object) Object {
	if isTruthy(left) {
		return TRUE
	}
	return nativeBoolToPyMonkeyBool(isTruthy(right))
}

func evalIfExpression(ie *parser.IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIdentifier(node *parser.Identifier, env *Environment) Object {
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalExpressions(exps []parser.Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index Object) Object {
	hashObject := hash.(*Hash)

	key, ok := index.(Hashable)
	if !ok {
		return newError("unusable as hash key: %T", index)
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalObjectLiteral(node *parser.ObjectLiteral, env *Environment) Object {
	pairs := make(map[HashKey]HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(Hashable)
		if !ok {
			return newError("unusable as hash key: %T", key)
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key: key, Value: value}
	}

	return &Hash{Pairs: pairs}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// Assignment expression evaluation
func evalAssignmentExpression(ae *parser.AssignmentExpression, env *Environment) Object {
	val := Eval(ae.Value, env)
	if isError(val) {
		return val
	}

	// Handle different assignment operators
	switch ae.Operator {
	case "=":
		env.Set(ae.Name.Value, val)
		return val
	case "+=":
		current, exists := env.Get(ae.Name.Value)
		if !exists {
			return newError("identifier not found: %s", ae.Name.Value)
		}
		result := evalInfixExpression("+", current, val)
		if isError(result) {
			return result
		}
		env.Set(ae.Name.Value, result)
		return result
	case "-=":
		current, exists := env.Get(ae.Name.Value)
		if !exists {
			return newError("identifier not found: %s", ae.Name.Value)
		}
		result := evalInfixExpression("-", current, val)
		if isError(result) {
			return result
		}
		env.Set(ae.Name.Value, result)
		return result
	case "*=":
		current, exists := env.Get(ae.Name.Value)
		if !exists {
			return newError("identifier not found: %s", ae.Name.Value)
		}
		result := evalInfixExpression("*", current, val)
		if isError(result) {
			return result
		}
		env.Set(ae.Name.Value, result)
		return result
	case "/=":
		current, exists := env.Get(ae.Name.Value)
		if !exists {
			return newError("identifier not found: %s", ae.Name.Value)
		}
		result := evalInfixExpression("/", current, val)
		if isError(result) {
			return result
		}
		env.Set(ae.Name.Value, result)
		return result
	default:
		return newError("unknown assignment operator: %s", ae.Operator)
	}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// Function application
func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// Loop evaluations
func evalWhileStatement(ws *parser.WhileStatement, env *Environment) Object {
	var result Object = NULL

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ws.Body, env)
		if result != nil {
			switch result.Type() {
			case RETURN_OBJ, ERROR_OBJ:
				return result
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				continue
			}
		}
	}

	return result
}

func evalForStatement(fs *parser.ForStatement, env *Environment) Object {
	// Create new environment for for loop scope
	forEnv := NewEnclosedEnvironment(env)

	// Initialize
	if fs.Initializer != nil {
		result := Eval(fs.Initializer, forEnv)
		if isError(result) {
			return result
		}
	}

	var result Object = NULL

	for {
		// Check condition
		if fs.Condition != nil {
			condition := Eval(fs.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// Execute body
		result = Eval(fs.Body, forEnv)
		if result != nil {
			switch result.Type() {
			case RETURN_OBJ, ERROR_OBJ:
				return result
			case BREAK_OBJ:
				return NULL
			case CONTINUE_OBJ:
				// Continue to increment
			default:
				// Continue normally
			}
		}

		// Increment
		if fs.Increment != nil {
			incrementResult := Eval(fs.Increment, forEnv)
			if isError(incrementResult) {
				return incrementResult
			}
		}
	}

	return result
}
