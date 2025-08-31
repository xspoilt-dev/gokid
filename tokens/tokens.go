package tokens

type TokenType string

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	POWER    = "**"

	// Assignment operators
	PLUS_ASSIGN     = "+="
	MINUS_ASSIGN    = "-="
	MULTIPLY_ASSIGN = "*="
	DIVIDE_ASSIGN   = "/="

	// Comparison operators
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LTE    = "<="
	GTE    = ">="

	// Logical operators
	AND = "&&"
	OR  = "||"
	NOT = "!"

	// Delimiters
	SEMICOLON = ";"
	COLON     = ":"
	COMMA     = ","
	DOT       = "."
	QUESTION  = "?"

	// Brackets and Braces
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Special characters
	AT    = "@"
	HASH  = "#"
	ARROW = "=>"

	// Keywords - Variables and Functions
	LET      = "LET"
	CONST    = "CONST"
	VAR      = "VAR"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"

	// Keywords - Control Flow
	IF       = "IF"
	ELSE     = "ELSE"
	WHILE    = "WHILE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"

	// Keywords - Data Types
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	NULL        = "NULL"
	STRING_TYPE = "STRING_TYPE"
	INT_TYPE    = "INT_TYPE"
	FLOAT_TYPE  = "FLOAT_TYPE"
	BOOL_TYPE   = "BOOL_TYPE"
	ARRAY_TYPE  = "ARRAY_TYPE"
	OBJECT_TYPE = "OBJECT_TYPE"

	// Keywords - Object Oriented
	CLASS   = "CLASS"
	THIS    = "THIS"
	NEW     = "NEW"
	EXTENDS = "EXTENDS"
	SUPER   = "SUPER"

	// Keywords - Exception Handling
	TRY     = "TRY"
	CATCH   = "CATCH"
	THROW   = "THROW"
	FINALLY = "FINALLY"

	// Keywords - Modules
	IMPORT = "IMPORT"
	EXPORT = "EXPORT"
	FROM   = "FROM"
	AS     = "AS"

	// Keywords - Scope
	GLOBAL = "GLOBAL"
	LOCAL  = "LOCAL"

	// Built-in Functions
	PRINT = "PRINT"
	LEN   = "LEN"
	TYPE  = "TYPE"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	// Variables and Functions
	"let":      LET,
	"const":    CONST,
	"var":      VAR,
	"fn":       FUNCTION,
	"function": FUNCTION,
	"return":   RETURN,

	// Control Flow
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,

	// Boolean and Null
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,

	// Data Types
	"string": STRING_TYPE,
	"int":    INT_TYPE,
	"float":  FLOAT_TYPE,
	"bool":   BOOL_TYPE,
	"array":  ARRAY_TYPE,
	"object": OBJECT_TYPE,

	// Object Oriented
	"class":   CLASS,
	"this":    THIS,
	"new":     NEW,
	"extends": EXTENDS,
	"super":   SUPER,

	// Exception Handling
	"try":     TRY,
	"catch":   CATCH,
	"throw":   THROW,
	"finally": FINALLY,

	// Modules
	"import": IMPORT,
	"export": EXPORT,
	"from":   FROM,
	"as":     AS,

	// Scope
	"global": GLOBAL,
	"local":  LOCAL,

	// Built-in Functions
	"print": PRINT,
	"len":   LEN,
	"type":  TYPE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
