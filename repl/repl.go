package repl

import (
	"bufio"
	"fmt"
	"gokid/evaluator"
	"gokid/lexer"
	"gokid/parser"
	"io"
)

const PROMPT = ">> "

const GOKID_FACE = `
    ____       _  ___     _ 
   / ___| ___ | |/ (_) __| |
  | |  _ / _ \| ' /| |/ _` + "`" + ` |
  | |_| | (_) | . \| | (_| |
   \____|\___/|_|\_\_|\__,_|

Welcome to the GoKid Programming Language!
Feel free to type in commands.
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := evaluator.NewEnvironment()

	fmt.Fprint(out, GOKID_FACE)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
