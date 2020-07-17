package evaluator

import (
	"concurrent-programming-christos-pitsikas/object"
	"concurrent-programming-christos-pitsikas/tree"
	"fmt"
)

// Eval : evaluation function
// Uncomment printing messages to check if we walk the tree correctly
func Eval(node tree.TreeNode, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *tree.Root:
		// fmt.Println("Evaluate Root")
		return evalProgram(node, env)
	case *tree.ExpressionStatement:
		// fmt.Println("Evaluate expression statement")
		return Eval(node.Expression, env)
	case *tree.IntegerLiteral:
		// fmt.Println("Evaluate integer")
		return &object.Integer{Value: node.Value}
	case *tree.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)
	case *tree.BlockStatement:
		return evalBlockStatement(node, env)
	case *tree.WhileExpression:
		// fmt.Println("Evaluate While expression")
		return evalWhileExpression(node, env)
	case *tree.IfExpression:
		// fmt.Println("Evaluate If expression")
		return evalIfExpression(node, env)
	case *tree.PrintStatement:
		// fmt.Println("Evaluate Print")
		val := Eval(node.Value, env)
		return &object.PrintValue{Value: val}
	case *tree.AssignStatement:
		// fmt.Println("Evaluate assignment")
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)
	case *tree.Identifier:
		// fmt.Println("Evaluate identifier")
		return evalIdentifier(node, env)
	}

	return nil
}

func evalProgram(program *tree.Root, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)
	}
	return result
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return nil
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case ">":
		if leftVal > rightVal {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case ">=":
		if leftVal > rightVal || leftVal == rightVal {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case "<":
		if leftVal < rightVal {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case "<=":
		if leftVal < rightVal || leftVal == rightVal {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case "==":
		if leftVal == rightVal {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case "!=":
		if leftVal != rightVal {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case "||":
		if leftVal == 1 || rightVal == 1 {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case "&&":
		if leftVal == 1 && rightVal == 1 {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	default:
		return nil
	}
}

func evalIfExpression(ie *tree.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	// printing messages to check if the correct branch is evaluated
	if checkCondition(condition) {
		fmt.Println("Evaluating true branch")
		return Eval(ie.TrueBranch, env)
	} else if ie.FalseBranch != nil {
		fmt.Println("Evaluating false branch")
		return Eval(ie.FalseBranch, env)
	} else {
		return nil
	}
}

func evalWhileExpression(we *tree.WhileExpression, env *object.Environment) object.Object {
	for {
		condition := Eval(we.Condition, env)

		if checkCondition(condition) {
			rt := Eval(we.Action, env)

			if rt != nil {
				return rt
			}
		} else {
			break
		}
	}
	return nil
}

func evalBlockStatement(block *tree.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			return result
		}
	}
	return result
}

func checkCondition(obj object.Object) bool {
	switch obj.Inspect() {
	case "1":
		return true
	case "0":
		return false
	default:
		return false
	}
}

func evalIdentifier(node *tree.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return nil
	}
	return val
}
