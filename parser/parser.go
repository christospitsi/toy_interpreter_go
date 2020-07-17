package parser

import (
	"concurrent-programming-christos-pitsikas/lexer"
	"concurrent-programming-christos-pitsikas/tree"
	"strconv"
)

// for precedence
const (
	_ int = iota // succesive integers, ignore first value by assigning to blank
	LOWEST
	LOGICALAND // &&
	LOGICALOR  // ||
	EQUALS     // ==
	LESSMORE   // > or <
	SUM        // + -
	PRODUCT    // * / %
	PREFIX     // to handle (
)

var precedences = map[lexer.TokenType]int{
	// precedences as in C language
	lexer.AND:     LOGICALAND,
	lexer.OR:      LOGICALOR,
	lexer.EQUAL:   EQUALS,
	lexer.N_EQUAL: EQUALS,
	lexer.LESS:    LESSMORE,
	lexer.MORE:    LESSMORE,
	lexer.MORE_EQ: LESSMORE,
	lexer.LESS_EQ: LESSMORE,
	lexer.PLUS:    SUM,
	lexer.MINUS:   SUM,
	lexer.DIVIDE:  PRODUCT,
	lexer.MULTIP:  PRODUCT,
	lexer.MODULO:  PRODUCT,
}

// Parser : Parse struct
type Parser struct {
	lex       *lexer.Lexer
	thisToken lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

// Parsing functions,
type (
	prefixParseFn func() tree.Expression
	infixParseFn  func(tree.Expression) tree.Expression
)

// ParsConstructor : parser constructor
func ParsConstructor(lex *lexer.Lexer) *Parser {
	pars := &Parser{lex: lex}

	pars.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	pars.infixParseFns = make(map[lexer.TokenType]infixParseFn)

	pars.registerInfix(lexer.PLUS, pars.parseInfixExpression)
	pars.registerInfix(lexer.MINUS, pars.parseInfixExpression)
	pars.registerInfix(lexer.DIVIDE, pars.parseInfixExpression)
	pars.registerInfix(lexer.MULTIP, pars.parseInfixExpression)
	pars.registerInfix(lexer.MODULO, pars.parseInfixExpression)
	pars.registerInfix(lexer.EQUAL, pars.parseInfixExpression)
	pars.registerInfix(lexer.N_EQUAL, pars.parseInfixExpression)
	pars.registerInfix(lexer.LESS, pars.parseInfixExpression)
	pars.registerInfix(lexer.LESS_EQ, pars.parseInfixExpression)
	pars.registerInfix(lexer.MORE, pars.parseInfixExpression)
	pars.registerInfix(lexer.MORE_EQ, pars.parseInfixExpression)
	pars.registerInfix(lexer.OR, pars.parseInfixExpression)
	pars.registerInfix(lexer.AND, pars.parseInfixExpression)

	pars.registerPrefix(lexer.IDENT, pars.parseIdentifier)
	pars.registerPrefix(lexer.NUM, pars.parseIntegerLiteral)
	pars.registerPrefix(lexer.LPAR, pars.parseGroupedExpression)
	pars.registerPrefix(lexer.IF, pars.parseIfExpression)
	pars.registerPrefix(lexer.WHILE, pars.parseWhileExpression)

	// read two tokens, one for thisToken and one for peekToken
	pars.nextToken()
	pars.nextToken()

	return pars
}

func (pars *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	pars.prefixParseFns[tokenType] = fn
}

func (pars *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	pars.infixParseFns[tokenType] = fn
}

// move to next thisToken and next peekToken
func (pars *Parser) nextToken() {
	pars.thisToken = pars.peekToken
	pars.peekToken = pars.lex.NextToken()
}

// ParseProgram : Parse the program statement until EOF is encountered
func (pars *Parser) ParseProgram() *tree.Root {
	program := &tree.Root{}
	program.Statements = []tree.Statement{}

	for !pars.curTokenIs(lexer.EOF) {
		stmt := pars.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		pars.nextToken()
	}
	return program
}

// parse assign or print statement
func (pars *Parser) parseStatement() tree.Statement {
	switch pars.thisToken.Type {
	case lexer.IDENT:
		if pars.peekTokenIs(lexer.ASSIGN) {
			return pars.parseAssignStatement()
		}
		return pars.parseExpressionStatement()
	case lexer.PRINT:
		return pars.parsePrintStatement()
	default:
		return pars.parseExpressionStatement()
	}
}

func (pars *Parser) parseAssignStatement() *tree.AssignStatement {
	stmt := &tree.AssignStatement{Token: pars.thisToken}

	stmt.Name = &tree.Identifier{Token: pars.thisToken, Value: pars.thisToken.Val}

	// check if next token is =
	if !pars.expectPeek(lexer.ASSIGN) {
		return nil
	}

	pars.nextToken()
	stmt.Value = pars.parseExpression(LOWEST)

	if pars.peekTokenIs(lexer.NEWLINE) {
		pars.nextToken()
	}

	return stmt
}

func (pars *Parser) parsePrintStatement() *tree.PrintStatement {
	stmt := &tree.PrintStatement{Token: pars.thisToken}
	pars.nextToken()

	stmt.Value = pars.parseExpression(LOWEST)

	// parse until carriage return
	for !pars.curTokenIs(lexer.NEWLINE) {
		pars.nextToken()
	}

	return stmt
}

func (pars *Parser) parseExpressionStatement() *tree.ExpressionStatement {
	stmt := &tree.ExpressionStatement{Token: pars.thisToken}

	stmt.Expression = pars.parseExpression(LOWEST)

	if pars.peekTokenIs(lexer.NEWLINE) {
		pars.nextToken()
	}
	return stmt
}

func (pars *Parser) parseExpression(precedence int) tree.Expression {
	prefix := pars.prefixParseFns[pars.thisToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for !pars.peekTokenIs(lexer.NEWLINE) && precedence < pars.peekPrecedence() {
		infix := pars.infixParseFns[pars.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		pars.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (pars *Parser) parseIdentifier() tree.Expression {
	return &tree.Identifier{Token: pars.thisToken, Value: pars.thisToken.Val}
}

func (pars *Parser) parseIntegerLiteral() tree.Expression {
	lit := &tree.IntegerLiteral{Token: pars.thisToken}

	value, err := strconv.ParseInt(pars.thisToken.Val, 0, 64)
	_ = err

	lit.Value = value

	return lit
}

// checks current token type
func (pars *Parser) curTokenIs(t lexer.TokenType) bool {
	return pars.thisToken.Type == t
}

// checks next token type
func (pars *Parser) peekTokenIs(t lexer.TokenType) bool {
	return pars.peekToken.Type == t
}

// checks if next token is what expected
func (pars *Parser) expectPeek(t lexer.TokenType) bool {
	if pars.peekTokenIs(t) {
		pars.nextToken()
		return true
	}
	return false
}

// check next precedence
func (pars *Parser) peekPrecedence() int {
	if pars, ok := precedences[pars.peekToken.Type]; ok {
		return pars
	}
	return LOWEST
}

// current precedence
func (pars *Parser) curPrecedence() int {
	if pars, ok := precedences[pars.thisToken.Type]; ok {
		return pars
	}
	return LOWEST
}

// parse inflix expressions
func (pars *Parser) parseInfixExpression(left tree.Expression) tree.Expression {
	expression := &tree.InfixExpression{
		Token:    pars.thisToken,
		Operator: pars.thisToken.Val,
		Left:     left,
	}

	precedence := pars.curPrecedence()
	pars.nextToken()
	expression.Right = pars.parseExpression(precedence)

	return expression
}

// parse expressions with ( prefix
func (pars *Parser) parsePrefixExpression() tree.Expression {
	expression := &tree.PrefixExpression{
		Token:    pars.thisToken,
		Operator: pars.thisToken.Val,
	}

	pars.nextToken()

	expression.Right = pars.parseExpression(PREFIX)
	return expression
}

// parse grouped expresisons in parentheses
func (pars *Parser) parseGroupedExpression() tree.Expression {
	pars.nextToken()

	exp := pars.parseExpression(LOWEST)
	if !pars.expectPeek(lexer.RPAR) {
		return nil
	}
	return exp
}

func (pars *Parser) parseIfExpression() tree.Expression {
	expression := &tree.IfExpression{Token: pars.thisToken}

	if !pars.expectPeek(lexer.LPAR) {
		return nil
	}

	pars.nextToken()
	expression.Condition = pars.parseExpression(LOWEST)

	if !pars.expectPeek(lexer.RPAR) {
		return nil
	}

	expression.TrueBranch = pars.parseBlockStatement()

	if pars.peekTokenIs(lexer.ELSE) {
		pars.nextToken()

		expression.FalseBranch = pars.parseBlockStatement()
	}
	return expression
}

func (pars *Parser) parseWhileExpression() tree.Expression {
	expression := &tree.WhileExpression{Token: pars.thisToken}

	if !pars.expectPeek(lexer.LPAR) {
		return nil
	}

	pars.nextToken()
	expression.Condition = pars.parseExpression(LOWEST)

	if !pars.expectPeek(lexer.RPAR) {
		return nil
	}

	expression.Action = pars.parseBlockStatement()

	return expression
}

// parse block of statements in {}
func (pars *Parser) parseBlockStatement() *tree.BlockStatement {
	block := &tree.BlockStatement{Token: pars.thisToken}
	block.Statements = []tree.Statement{}

	pars.nextToken()

	for !pars.curTokenIs(lexer.RBRAC) && !pars.curTokenIs(lexer.EOF) {
		stmt := pars.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		pars.nextToken()
	}
	return block
}
