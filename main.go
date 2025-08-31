package main

import (
	"fmt"
	"gokid/tokenizer"
)

func main() {
	input := `
let x = 5 + 10;
let y = x - 2;
fn add(a, b) {
    a + b;
}
`

	t := tokenizer.NewTokenizer(input)
	tokens := t.GetTokens()

	fmt.Println("Tokens generated:")
	for _, tok := range tokens {
		fmt.Printf("%-10s : %s\n", tok.Type, tok.Literal)
	}
}
