package main

import (
	"bufio"
	"fmt"
	"gokid/evaluator"
	"gokid/lexer"
	"gokid/parser"
	"gokid/repl"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please specify a .gokid file to run")
			fmt.Println("Usage: gokid run <file.gokid>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "repl", "interactive":
		startREPL()
	case "version", "--version", "-v":
		printVersion()
	case "help", "--help", "-h":
		printHelp()
	default:
		// If it ends with .gokid, try to run it
		if strings.HasSuffix(command, ".gokid") {
			runFile(command)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Printf("GoKid Language Interpreter v%s\n", VERSION)
	fmt.Println("Created by xspoilt-dev")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gokid run <file.gokid>    Execute a GoKid source file")
	fmt.Println("  gokid repl               Start interactive REPL")
	fmt.Println("  gokid <file.gokid>       Execute a GoKid source file (shorthand)")
	fmt.Println("  gokid version            Show version information")
	fmt.Println("  gokid help               Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gokid run hello.gokid")
	fmt.Println("  gokid hello.gokid")
	fmt.Println("  gokid repl")
}

func printVersion() {
	fmt.Printf("GoKid Language Interpreter v%s\n", VERSION)
	fmt.Println("Created by xspoilt-dev")
	fmt.Println("GitHub: https://github.com/xspoilt-dev/gokid")
}

func printHelp() {
	printUsage()
	fmt.Println()
	fmt.Println("GoKid Language Features:")
	fmt.Println("  • Variables: let, const, var")
	fmt.Println("  • Data types: integers, floats, strings, booleans, arrays, objects")
	fmt.Println("  • Operators: arithmetic, comparison, logical, assignment")
	fmt.Println("  • Control flow: if/else statements, while loops")
	fmt.Println("  • Functions: function expressions and calls")
	fmt.Println("  • Built-ins: print(), len(), type()")
	fmt.Println("  • Interactive REPL for development and testing")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/xspoilt-dev/gokid")
}

func runFile(filename string) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found\n", filename)
		os.Exit(1)
	}

	// Check file extension
	if !strings.HasSuffix(filename, ".gokid") {
		fmt.Printf("Warning: File '%s' doesn't have .gokid extension\n", filename)
		fmt.Print("Continue anyway? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Execution cancelled.")
			return
		}
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	// Get absolute path for better error reporting
	absPath, _ := filepath.Abs(filename)

	fmt.Printf("Executing: %s\n", absPath)
	fmt.Println(strings.Repeat("-", 50))

	// Execute the program
	executeProgram(string(content), filename)
}

func executeProgram(source string, filename string) {
	// Create lexer and parser
	l := lexer.NewLexer(source)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing errors
	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Printf("Parsing errors in %s:\n", filename)
		for i, err := range errors {
			fmt.Printf("  %d: %s\n", i+1, err)
		}
		os.Exit(1)
	}

	// Execute the program
	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)

	// Handle runtime errors
	if result != nil && result.Type() == "ERROR" {
		fmt.Printf("Runtime error: %s\n", result.Inspect())
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Program executed successfully.")
}

func startREPL() {
	fmt.Printf("GoKid Language REPL v%s\n", VERSION)
	fmt.Println("Created by xspoilt-dev")
	fmt.Println("Type 'exit' or press Ctrl+C to quit")
	fmt.Println(strings.Repeat("-", 40))

	repl.Start(os.Stdin, os.Stdout)
}
