package parser

import (
	"fmt"
	"gokid/lexer"
	"gokid/tokens"
	"strconv"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN      // =, +=, -=, etc.
	TERNARY     // ? :
	OR          // ||
	AND         // &&
	EQUALS      // ==, !=
	LESSGREATER // > or <, >=, <=
	SUM         // +, -
	PRODUCT     // *, /, %
	POWER       // **
	PREFIX      // -X, !X
	CALL        // myFunction(X)
	INDEX       // array[index], obj.prop
)

// Precedence map
var precedences = map[tokens.TokenType]int{
	tokens.ASSIGN:          ASSIGN,
	tokens.PLUS_ASSIGN:     ASSIGN,
	tokens.MINUS_ASSIGN:    ASSIGN,
	tokens.MULTIPLY_ASSIGN: ASSIGN,
	tokens.DIVIDE_ASSIGN:   ASSIGN,
	tokens.QUESTION:        TERNARY,
	tokens.OR:              OR,
	tokens.AND:             AND,
	tokens.EQ:              EQUALS,
	tokens.NOT_EQ:          EQUALS,
	tokens.LT:              LESSGREATER,
	tokens.GT:              LESSGREATER,
	tokens.LTE:             LESSGREATER,
	tokens.GTE:             LESSGREATER,
	tokens.PLUS:            SUM,
	tokens.MINUS:           SUM,
	tokens.SLASH:           PRODUCT,
	tokens.ASTERISK:        PRODUCT,
	tokens.MODULO:          PRODUCT,
	tokens.POWER:           POWER,
	tokens.LPAREN:          CALL,
	tokens.LBRACKET:        INDEX,
	tokens.DOT:             INDEX,
}

// Parser function types
type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// Parser struct
type Parser struct {
	l *lexer.Lexer

	curToken  tokens.Token
	peekToken tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn

	errors []string
}

// New creates a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize parse function maps
	p.prefixParseFns = make(map[tokens.TokenType]prefixParseFn)
	p.infixParseFns = make(map[tokens.TokenType]infixParseFn)

	// Register prefix parse functions
	p.registerPrefix(tokens.IDENT, p.parseIdentifier)
	p.registerPrefix(tokens.INT, p.parseIntegerLiteral)
	p.registerPrefix(tokens.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(tokens.STRING, p.parseStringLiteral)
	p.registerPrefix(tokens.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(tokens.NULL, p.parseNullLiteral)
	p.registerPrefix(tokens.NOT, p.parsePrefixExpression)
	p.registerPrefix(tokens.MINUS, p.parsePrefixExpression)
	p.registerPrefix(tokens.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(tokens.IF, p.parseIfExpression)
	p.registerPrefix(tokens.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(tokens.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(tokens.LBRACE, p.parseObjectLiteral)

	// Register built-in functions as identifiers (they'll be handled by parseIdentifier)
	p.registerPrefix(tokens.PRINT, p.parseIdentifier)
	p.registerPrefix(tokens.LEN, p.parseIdentifier)
	p.registerPrefix(tokens.TYPE, p.parseIdentifier)

	// Register infix parse functions
	p.registerInfix(tokens.PLUS, p.parseInfixExpression)
	p.registerInfix(tokens.MINUS, p.parseInfixExpression)
	p.registerInfix(tokens.SLASH, p.parseInfixExpression)
	p.registerInfix(tokens.ASTERISK, p.parseInfixExpression)
	p.registerInfix(tokens.MODULO, p.parseInfixExpression)
	p.registerInfix(tokens.POWER, p.parseInfixExpression)
	p.registerInfix(tokens.EQ, p.parseInfixExpression)
	p.registerInfix(tokens.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(tokens.LT, p.parseInfixExpression)
	p.registerInfix(tokens.GT, p.parseInfixExpression)
	p.registerInfix(tokens.LTE, p.parseInfixExpression)
	p.registerInfix(tokens.GTE, p.parseInfixExpression)
	p.registerInfix(tokens.AND, p.parseInfixExpression)
	p.registerInfix(tokens.OR, p.parseInfixExpression)
	p.registerInfix(tokens.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(tokens.PLUS_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(tokens.MINUS_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(tokens.MULTIPLY_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(tokens.DIVIDE_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(tokens.LPAREN, p.parseCallExpression)
	p.registerInfix(tokens.LBRACKET, p.parseIndexExpression)
	p.registerInfix(tokens.DOT, p.parseDotExpression)
	p.registerInfix(tokens.QUESTION, p.parseTernaryExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Helper methods
func (p *Parser) registerPrefix(tokenType tokens.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType tokens.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t tokens.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Error handling
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t tokens.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// Main parsing method
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.curTokenIs(tokens.EOF) {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

// Statement parsing
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.CONST:
		return p.parseConstStatement()
	case tokens.VAR:
		return p.parseVarStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	case tokens.FUNCTION:
		return p.parseFunctionStatement()
	case tokens.WHILE:
		return p.parseWhileStatement()
	case tokens.FOR:
		return p.parseForStatement()
	case tokens.BREAK:
		return p.parseBreakStatement()
	case tokens.CONTINUE:
		return p.parseContinueStatement()
	case tokens.SWITCH:
		return p.parseSwitchStatement()
	case tokens.TRY:
		return p.parseTryStatement()
	case tokens.THROW:
		return p.parseThrowStatement()
	case tokens.IMPORT:
		return p.parseImportStatement()
	case tokens.EXPORT:
		return p.parseExportStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(tokens.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(tokens.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ConstStatement {
	stmt := &ConstStatement{Token: p.curToken}

	if !p.expectPeek(tokens.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(tokens.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() *VarStatement {
	stmt := &VarStatement{Token: p.curToken}

	if !p.expectPeek(tokens.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(tokens.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs(tokens.RBRACE) && !p.curTokenIs(tokens.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block
}

// Control flow statements
func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	// Initializer
	p.nextToken()
	if !p.curTokenIs(tokens.SEMICOLON) {
		stmt.Initializer = p.parseStatement()
	}

	if !p.expectPeek(tokens.SEMICOLON) {
		return nil
	}

	// Condition
	p.nextToken()
	if !p.curTokenIs(tokens.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(tokens.SEMICOLON) {
		return nil
	}

	// Increment
	p.nextToken()
	if !p.curTokenIs(tokens.RPAREN) {
		stmt.Increment = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBreakStatement() *BreakStatement {
	stmt := &BreakStatement{Token: p.curToken}

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ContinueStatement {
	stmt := &ContinueStatement{Token: p.curToken}

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Function statement (function declarations)
func (p *Parser) parseFunctionStatement() *ExpressionStatement {
	// Function declarations are actually function literals assigned to the global scope
	// We'll parse them as expression statements for now
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseFunctionLiteral()

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseSwitchStatement() *SwitchStatement {
	stmt := &SwitchStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	p.nextToken()

	for !p.curTokenIs(tokens.RBRACE) && !p.curTokenIs(tokens.EOF) {
		if p.curTokenIs(tokens.CASE) {
			caseStmt := p.parseCaseStatement()
			if caseStmt != nil {
				stmt.Cases = append(stmt.Cases, caseStmt)
			}
		} else if p.curTokenIs(tokens.DEFAULT) {
			stmt.Default = p.parseDefaultStatement()
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseCaseStatement() *CaseStatement {
	stmt := &CaseStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.COLON) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseDefaultStatement() *DefaultStatement {
	stmt := &DefaultStatement{Token: p.curToken}

	if !p.expectPeek(tokens.COLON) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// Exception handling
func (p *Parser) parseTryStatement() *TryStatement {
	stmt := &TryStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	if p.peekTokenIs(tokens.CATCH) {
		p.nextToken()
		stmt.Catch = p.parseCatchStatement()
	}

	if p.peekTokenIs(tokens.FINALLY) {
		p.nextToken()
		stmt.Finally = p.parseFinallyStatement()
	}

	return stmt
}

func (p *Parser) parseCatchStatement() *CatchStatement {
	stmt := &CatchStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.IDENT) {
		return nil
	}

	stmt.Parameter = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFinallyStatement() *FinallyStatement {
	stmt := &FinallyStatement{Token: p.curToken}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseThrowStatement() *ThrowStatement {
	stmt := &ThrowStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Import/Export
func (p *Parser) parseImportStatement() *ImportStatement {
	stmt := &ImportStatement{Token: p.curToken}

	if !p.expectPeek(tokens.STRING) {
		return nil
	}

	stmt.Path = &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(tokens.AS) {
		p.nextToken()
		if !p.expectPeek(tokens.IDENT) {
			return nil
		}
		stmt.Alias = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.peekTokenIs(tokens.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExportStatement() *ExportStatement {
	stmt := &ExportStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseStatement()

	return stmt
}

// Expression parsing
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(tokens.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix expressions
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(tokens.TRUE)}
}

func (p *Parser) parseNullLiteral() Expression {
	return &NullLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(tokens.ELSE) {
		p.nextToken()

		if !p.expectPeek(tokens.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(tokens.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(tokens.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekTokenIs(tokens.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(tokens.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(tokens.RBRACKET)

	return array
}

func (p *Parser) parseObjectLiteral() Expression {
	obj := &ObjectLiteral{Token: p.curToken}
	obj.Pairs = make(map[Expression]Expression)

	for !p.peekTokenIs(tokens.RBRACE) && !p.peekTokenIs(tokens.EOF) {
		p.nextToken()

		key := p.parseExpression(LOWEST)

		if !p.expectPeek(tokens.COLON) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)

		obj.Pairs[key] = value

		if !p.peekTokenIs(tokens.RBRACE) && !p.expectPeek(tokens.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(tokens.RBRACE) {
		return nil
	}

	return obj
}

// Infix expressions
func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	ident, ok := left.(*Identifier)
	if !ok {
		msg := fmt.Sprintf("expected identifier, got %T", left)
		p.errors = append(p.errors, msg)
		return nil
	}

	expression := &AssignmentExpression{
		Token:    p.curToken,
		Name:     ident,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Value = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(fn Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(tokens.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseDotExpression(left Expression) Expression {
	exp := &DotExpression{Token: p.curToken, Left: left}

	if !p.expectPeek(tokens.IDENT) {
		return nil
	}

	exp.Property = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return exp
}

func (p *Parser) parseTernaryExpression(condition Expression) Expression {
	exp := &TernaryExpression{Token: p.curToken, Condition: condition}

	p.nextToken()
	exp.Consequence = p.parseExpression(LOWEST)

	if !p.expectPeek(tokens.COLON) {
		return nil
	}

	p.nextToken()
	exp.Alternative = p.parseExpression(LOWEST)

	return exp
}

// Helper methods
func (p *Parser) parseExpressionList(end tokens.TokenType) []Expression {
	args := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}
