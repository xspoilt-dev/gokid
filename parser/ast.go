package parser

import "gokid/tokens"

// Base interfaces
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program - root node
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// Identifier
type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// Literals
type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

type FloatLiteral struct {
	Token tokens.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

type StringLiteral struct {
	Token tokens.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

type BooleanLiteral struct {
	Token tokens.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

type NullLiteral struct {
	Token tokens.Token
}

func (nl *NullLiteral) expressionNode() {}
func (nl *NullLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

// Array Literal
type ArrayLiteral struct {
	Token    tokens.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// Object Literal
type ObjectLiteral struct {
	Token tokens.Token
	Pairs map[Expression]Expression
}

func (ol *ObjectLiteral) expressionNode() {}
func (ol *ObjectLiteral) TokenLiteral() string {
	return ol.Token.Literal
}

// Variable Declarations
type LetStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type ConstStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode() {}
func (cs *ConstStatement) TokenLiteral() string {
	return cs.Token.Literal
}

type VarStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

// Return Statement
type ReturnStatement struct {
	Token       tokens.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// Expression Statement
type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// Block Statement
type BlockStatement struct {
	Token      tokens.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// Function Literal
type FunctionLiteral struct {
	Token      tokens.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// Call Expression
type CallExpression struct {
	Token     tokens.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// Prefix Expression
type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// Infix Expression
type InfixExpression struct {
	Token    tokens.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// If Expression
type IfExpression struct {
	Token       tokens.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// While Statement
type WhileStatement struct {
	Token     tokens.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}

// For Statement
type ForStatement struct {
	Token       tokens.Token
	Initializer Statement
	Condition   Expression
	Increment   Expression
	Body        *BlockStatement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

// Break Statement
type BreakStatement struct {
	Token tokens.Token
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// Continue Statement
type ContinueStatement struct {
	Token tokens.Token
}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// Switch Statement
type SwitchStatement struct {
	Token   tokens.Token
	Value   Expression
	Cases   []*CaseStatement
	Default *DefaultStatement
}

func (ss *SwitchStatement) statementNode() {}
func (ss *SwitchStatement) TokenLiteral() string {
	return ss.Token.Literal
}

// Case Statement
type CaseStatement struct {
	Token tokens.Token
	Value Expression
	Body  *BlockStatement
}

func (cs *CaseStatement) statementNode() {}
func (cs *CaseStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// Default Statement
type DefaultStatement struct {
	Token tokens.Token
	Body  *BlockStatement
}

func (ds *DefaultStatement) statementNode() {}
func (ds *DefaultStatement) TokenLiteral() string {
	return ds.Token.Literal
}

// Try Statement
type TryStatement struct {
	Token   tokens.Token
	Body    *BlockStatement
	Catch   *CatchStatement
	Finally *FinallyStatement
}

func (ts *TryStatement) statementNode() {}
func (ts *TryStatement) TokenLiteral() string {
	return ts.Token.Literal
}

// Catch Statement
type CatchStatement struct {
	Token     tokens.Token
	Parameter *Identifier
	Body      *BlockStatement
}

func (cs *CatchStatement) statementNode() {}
func (cs *CatchStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// Finally Statement
type FinallyStatement struct {
	Token tokens.Token
	Body  *BlockStatement
}

func (fs *FinallyStatement) statementNode() {}
func (fs *FinallyStatement) TokenLiteral() string {
	return fs.Token.Literal
}

// Throw Statement
type ThrowStatement struct {
	Token tokens.Token
	Value Expression
}

func (ts *ThrowStatement) statementNode() {}
func (ts *ThrowStatement) TokenLiteral() string {
	return ts.Token.Literal
}

// Import Statement
type ImportStatement struct {
	Token tokens.Token
	Path  *StringLiteral
	Alias *Identifier
}

func (is *ImportStatement) statementNode() {}
func (is *ImportStatement) TokenLiteral() string {
	return is.Token.Literal
}

// Export Statement
type ExportStatement struct {
	Token tokens.Token
	Value Statement
}

func (es *ExportStatement) statementNode() {}
func (es *ExportStatement) TokenLiteral() string {
	return es.Token.Literal
}

// Assignment Expression
type AssignmentExpression struct {
	Token    tokens.Token
	Name     *Identifier
	Operator string
	Value    Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) TokenLiteral() string {
	return ae.Token.Literal
}

// Index Expression
type IndexExpression struct {
	Token tokens.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// Dot Expression (for object property access)
type DotExpression struct {
	Token    tokens.Token
	Left     Expression
	Property *Identifier
}

func (de *DotExpression) expressionNode() {}
func (de *DotExpression) TokenLiteral() string {
	return de.Token.Literal
}

// Ternary Expression
type TernaryExpression struct {
	Token       tokens.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode() {}
func (te *TernaryExpression) TokenLiteral() string {
	return te.Token.Literal
}
