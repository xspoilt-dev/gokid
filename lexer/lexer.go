package lexer

import (
	"gokid/tokens"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() tokens.Token {
	var tok tokens.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.EQ, Literal: literal}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.ARROW, Literal: literal}
		} else {
			tok = newToken(tokens.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.PLUS_ASSIGN, Literal: literal}
		} else {
			tok = newToken(tokens.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.MINUS_ASSIGN, Literal: literal}
		} else {
			tok = newToken(tokens.MINUS, l.ch)
		}
	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.POWER, Literal: literal}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.MULTIPLY_ASSIGN, Literal: literal}
		} else {
			tok = newToken(tokens.ASTERISK, l.ch)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.DIVIDE_ASSIGN, Literal: literal}
		} else {
			tok = newToken(tokens.SLASH, l.ch)
		}
	case '%':
		tok = newToken(tokens.MODULO, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(tokens.NOT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.LTE, Literal: literal}
		} else {
			tok = newToken(tokens.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.GTE, Literal: literal}
		} else {
			tok = newToken(tokens.GT, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.AND, Literal: literal}
		} else {
			tok = newToken(tokens.ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = tokens.Token{Type: tokens.OR, Literal: literal}
		} else {
			tok = newToken(tokens.ILLEGAL, l.ch)
		}
	case '(':
		tok = newToken(tokens.LPAREN, l.ch)
	case ')':
		tok = newToken(tokens.RPAREN, l.ch)
	case '{':
		tok = newToken(tokens.LBRACE, l.ch)
	case '}':
		tok = newToken(tokens.RBRACE, l.ch)
	case '[':
		tok = newToken(tokens.LBRACKET, l.ch)
	case ']':
		tok = newToken(tokens.RBRACKET, l.ch)
	case ',':
		tok = newToken(tokens.COMMA, l.ch)
	case ';':
		tok = newToken(tokens.SEMICOLON, l.ch)
	case ':':
		tok = newToken(tokens.COLON, l.ch)
	case '.':
		tok = newToken(tokens.DOT, l.ch)
	case '?':
		tok = newToken(tokens.QUESTION, l.ch)
	case '@':
		tok = newToken(tokens.AT, l.ch)
	case '#':
		tok = newToken(tokens.HASH, l.ch)
	case '"':
		tok.Type = tokens.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = tokens.EOF
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := tokens.LookupIdent(literal)
			tok = tokens.Token{Type: tokType, Literal: literal}
			return tok
		} else if isDigit(l.ch) {
			literal, tokenType := l.readNumber()
			tok = tokens.Token{Type: tokenType, Literal: literal}
			return tok
		} else {
			tok = newToken(tokens.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType tokens.TokenType, ch byte) tokens.Token {
	return tokens.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() (string, tokens.TokenType) {
	pos := l.position
	var tokenType tokens.TokenType = tokens.INT

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = tokens.FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[pos:l.position], tokenType
}

func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.position]
}
