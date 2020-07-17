package tree

import (
	"bytes"
	"concurrent-programming-christos-pitsikas/lexer"
)

// TreeNode : node interface
type TreeNode interface {
	// TokenVal() returns the value of the token
	TokenVal() string
	String() string
}

// Statement : statement type
type Statement interface {
	TreeNode
	statementNode()
}

// Expression : expression interface
type Expression interface {
	TreeNode
	expressionNode()
}

// Root : the root node of every tree
type Root struct {
	Statements []Statement
}

// The whole programm as a string
func (r *Root) String() string {
	var out bytes.Buffer

	for _, s := range r.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// TokenVal : return the value of the token
func (r *Root) TokenVal() string {
	if len(r.Statements) > 0 {
		return r.Statements[0].TokenVal()
	}
	return ""
}

// AssignStatement : Assignment
type AssignStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()   {}
func (as *AssignStatement) TokenVal() string { return as.Token.Val }
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	return out.String()
}

// Identifier - to hold the identifier of the binding
type Identifier struct {
	Token lexer.Token // IDENT token
	Value string
}

func (i *Identifier) expressionNode()  {}
func (i *Identifier) TokenVal() string { return i.Token.Val }
func (i *Identifier) String() string   { return i.Value }

// PrintStatement : Print statement
type PrintStatement struct {
	Token lexer.Token // PRINT token
	Value Expression
}

func (ps *PrintStatement) statementNode()   {}
func (ps *PrintStatement) TokenVal() string { return ps.Token.Val }
func (ps *PrintStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ps.TokenVal() + " ")
	if ps.Value != nil {
		out.WriteString(ps.Value.String())
	}

	// out.WriteString("\n") // Sunday
	return out.String()
}

// ExpressionStatement : Expressions
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()   {}
func (es *ExpressionStatement) TokenVal() string { return es.Token.Val }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral : integers
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()  {}
func (il *IntegerLiteral) TokenVal() string { return il.Token.Val }
func (il *IntegerLiteral) String() string   { return il.Token.Val }

// InfixExpression : operators like +, - etc
type InfixExpression struct {
	Token    lexer.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()  {}
func (ie *InfixExpression) TokenVal() string { return ie.Token.Val }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type PrefixExpression struct {
	Token    lexer.Token // prefix for (
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()  {}
func (pe *PrefixExpression) TokenVal() string { return pe.Token.Val }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	TrueBranch  *BlockStatement
	FalseBranch *BlockStatement
}

func (ie *IfExpression) expressionNode()  {}
func (ie *IfExpression) TokenVal() string { return ie.Token.Val }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())

	out.WriteString(" ")
	out.WriteString(ie.TrueBranch.String())

	if ie.FalseBranch != nil {
		out.WriteString("else ")
		out.WriteString(ie.FalseBranch.String())
	}

	return out.String()
}

type WhileExpression struct {
	Token     lexer.Token // The 'while' token
	Condition Expression
	Action    *BlockStatement
}

func (we *WhileExpression) expressionNode()  {}
func (we *WhileExpression) TokenVal() string { return we.Token.Val }
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(we.Condition.String())

	out.WriteString(" ")
	out.WriteString(we.Action.String())

	return out.String()
}

type BlockStatement struct {
	Token      lexer.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()   {}
func (bs *BlockStatement) TokenVal() string { return bs.Token.Val }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
