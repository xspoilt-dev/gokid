package tokenizer

import (
	"gokid/lexer"
	"gokid/tokens"
)

type Tokenizer struct {
	lexer *lexer.Lexer
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		lexer: lexer.NewLexer(input),
	}
}

func (t *Tokenizer) GetTokens() []tokens.Token {
	var allTokens []tokens.Token
	for tok := t.lexer.NextToken(); tok.Type != tokens.EOF; tok = t.lexer.NextToken() {
		allTokens = append(allTokens, tok)
	}
	return allTokens
}
