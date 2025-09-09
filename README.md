# GoKid Programming Language 🚀

<div align="center">

![GoKid Logo](https://img.shields.io/badge/GoKid-v1.0.0-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)

*A simple, yet powerful interpreted programming language built in Go*

**Created by [xspoilt-dev](https://github.com/xspoilt-dev)**

</div>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Language Syntax](#-language-syntax)
- [Built-in Functions](#-built-in-functions)
- [Examples](#-examples)
- [Architecture](#-architecture)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 Overview

**GoKid** is a tree-walking interpreted programming language implemented in Go. It features a clean, JavaScript-like syntax with support for variables, functions, control flow, and built-in data structures. GoKid is designed to be simple to learn while providing powerful programming constructs.

### Key Highlights
- 🏗️ **Tree-walking interpreter** with real-time execution
- 📝 **Clean, familiar syntax** inspired by JavaScript and C-family languages
- 🔧 **Interactive REPL** for development and testing
- 🎪 **Dynamic typing** with support for multiple data types
- ⚡ **Fast execution** with efficient parsing and evaluation
- 🧩 **Extensible architecture** for adding new features

---

## ✨ Features

### 🔤 Data Types
- **Numbers**: Integers (`42`) and Floats (`3.14`)
- **Strings**: `"Hello, World!"`
- **Booleans**: `true`, `false`
- **Arrays**: `[1, 2, 3, "mixed", true]`
- **Objects**: `{"name": "John", "age": 30}`
- **Functions**: First-class functions with closures
- **Null**: `null` value

### 🎮 Control Flow
- **Conditionals**: `if/else` statements (simple form)
- **Loops**: `while` loops with `break` and `continue`
- **Function calls** with parameters and return values

*Note: `else if` chaining and advanced control structures are planned for future releases*

### 🛠️ Operators
- **Arithmetic**: `+`, `-`, `*`, `/`, `**` (power)
- **Comparison**: `==`, `!=`, `<`, `>` 
- **Logical**: `&&`, `||`, `!`
- **Assignment**: `=`, `+=`, `-=`, `*=`, `/=`

*Note: `<=` and `>=` operators are planned for future releases*

### 🎛️ Language Constructs
- **Variable declarations**: `let`, `const`, `var`
- **Function expressions**: `let add = function(a, b) { return a + b; }`
- **Object property access**: `obj["property"]` (bracket notation)
- **Array indexing**: `arr[0]`
- **Comments**: `// Single line comments`

*Note: Dot notation for object properties (`obj.property`) is planned for future releases*

---

## 📦 Installation

### Prerequisites
- **Go 1.19+** installed on your system
- Basic familiarity with command line

### Build from Source

```bash
# Clone the repository
git clone https://github.com/xspoilt-dev/gokid.git
cd gokid

# Build the interpreter
go build -o gokid main.go

# Or run directly
go run main.go --help
```

### Quick Test
```bash
# Test the installation
go run main.go version
```

---

## 🚀 Quick Start

### 1. Create Your First GoKid Program

Create a file called `hello.gokid`:

```javascript
// hello.gokid
print("Hello, GoKid World!");

let name = "Developer";
let greeting = "Welcome to GoKid, " + name + "!";
print(greeting);

let numbers = [1, 2, 3, 4, 5];
print("Array length: " + len(numbers));
```

### 2. Run Your Program

```bash
# Method 1: Direct execution
go run main.go run hello.gokid

# Method 2: Shorthand
go run main.go hello.gokid

# Method 3: After building
./gokid run hello.gokid
```

### 3. Interactive Development

```bash
# Start the REPL
go run main.go repl

# Now you can type GoKid code interactively:
GoKid Language REPL v1.0.0
Created by xspoilt-dev
Type 'exit' or press Ctrl+C to quit
----------------------------------------
>> let x = 42;
>> print("The answer is: " + x);
The answer is: 42
>> let double = function(n) { return n * 2; };
>> double(21);
42
```

---

## 📚 Language Syntax

### Variables

```javascript
// Variable declarations
let mutableVar = 42;          // Mutable variable
const immutableVar = "hello"; // Immutable variable  
var legacyVar = true;         // Legacy style variable

// Assignment operations
mutableVar = 100;             // Basic assignment
mutableVar += 10;             // Addition assignment (110)
mutableVar *= 2;              // Multiplication assignment (220)
mutableVar /= 4;              // Division assignment (55)
```

### Functions

```javascript
// Function expressions
let add = function(a, b) {
    return a + b;
};

let factorial = function(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
};

// Function calls
let result = add(5, 3);        // 8
let fact5 = factorial(5);      // 120
```

### Control Flow

```javascript
// Conditionals
if (x > 0) {
    print("Positive");
} else if (x < 0) {
    print("Negative");
} else {
    print("Zero");
}

// While loops
let i = 0;
while (i < 5) {
    print("Count: " + i);
    i += 1;
}

// Loop control
let j = 0;
while (true) {
    if (j >= 10) {
        break;        // Exit loop
    }
    if (j % 2 == 0) {
        j += 1;
        continue;     // Skip to next iteration
    }
    print("Odd: " + j);
    j += 1;
}
```

### Data Structures

```javascript
// Arrays
let fruits = ["apple", "banana", "orange"];
fruits[0] = "grape";                    // Modify element
let firstFruit = fruits[0];             // Access element
print("First fruit: " + firstFruit);

// Objects
let person = {
    name: "Alice",
    age: 30,
    city: "New York"
};

// Property access
let name = person.name;           // Dot notation
let age = person["age"];          // Bracket notation
person.age = 31;                  // Modify property
```

### Advanced Examples

```javascript
// Closures
let createCounter = function() {
    let count = 0;
    return function() {
        count += 1;
        return count;
    };
};

let counter = createCounter();
print(counter()); // 1
print(counter()); // 2
print(counter()); // 3

// Higher-order functions
let map = function(arr, fn) {
    let result = [];
    let i = 0;
    while (i < len(arr)) {
        result[i] = fn(arr[i]);
        i += 1;
    }
    return result;
};

let numbers = [1, 2, 3, 4, 5];
let doubled = map(numbers, function(x) { return x * 2; });
print(doubled); // [2, 4, 6, 8, 10]
```

---

## 🔧 Built-in Functions

### `print(value)`
Outputs a value to the console.

```javascript
print("Hello World");
print(42);
print([1, 2, 3]);
```

### `len(collection)`
Returns the length of arrays, objects, or strings.

```javascript
len([1, 2, 3]);           // 3
len("hello");             // 5
len({a: 1, b: 2});        // 2
```

### `type(value)`
Returns the type of a value as a string.

```javascript
type(42);                 // "INTEGER"
type("hello");            // "STRING" 
type([1, 2, 3]);          // "ARRAY"
type({name: "John"});     // "HASH"
type(true);               // "BOOLEAN"
```

---

## 💡 Examples

### 1. Fibonacci Sequence

```javascript
// fibonacci.gokid
let fibonacci = function(n) {
    if (n <= 1) {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
};

let i = 0;
while (i <= 10) {
    print("fibonacci(" + i + ") = " + fibonacci(i));
    i += 1;
}
```

### 2. Array Processing

```javascript
// array_utils.gokid
let findMax = function(arr) {
    if (len(arr) == 0) {
        return null;
    }
    
    let max = arr[0];
    let i = 1;
    while (i < len(arr)) {
        if (arr[i] > max) {
            max = arr[i];
        }
        i += 1;
    }
    return max;
};

let numbers = [3, 7, 2, 9, 1, 5];
print("Array: " + numbers);
print("Maximum: " + findMax(numbers));
```

### 3. Simple Calculator

```javascript
// calculator.gokid
let calculator = function(op, a, b) {
    if (op == "+") {
        return a + b;
    } else if (op == "-") {
        return a - b;
    } else if (op == "*") {
        return a * b;
    } else if (op == "/") {
        if (b == 0) {
            print("Error: Division by zero");
            return null;
        }
        return a / b;
    } else {
        print("Error: Unknown operation");
        return null;
    }
};

print("5 + 3 = " + calculator("+", 5, 3));
print("10 - 4 = " + calculator("-", 10, 4));
print("6 * 7 = " + calculator("*", 6, 7));
print("15 / 3 = " + calculator("/", 15, 3));
```

---

## 🏗️ Architecture

GoKid follows a clean, modular architecture:

```
gokid/
├── main.go           # CLI interface and file execution
├── tokens/           # Token type definitions
│   └── tokens.go
├── lexer/           # Lexical analysis (tokenization)
│   └── lexer.go
├── parser/          # Syntax analysis (AST generation)
│   ├── parser.go
│   └── ast.go
├── evaluator/       # Semantic analysis and execution
│   ├── evaluator.go
│   ├── object.go
│   ├── environment.go
│   └── builtins.go
├── repl/            # Interactive Read-Eval-Print Loop
│   └── repl.go
└── tokenizer/       # High-level tokenization utilities
    └── tokenizer.go
```

### Processing Pipeline

1. **Lexical Analysis**: Source code → Tokens
2. **Syntax Analysis**: Tokens → Abstract Syntax Tree (AST)
3. **Semantic Analysis**: AST → Execution with Environment
4. **Evaluation**: Direct interpretation of AST nodes

### Key Components

- **Lexer**: Converts source text into tokens
- **Parser**: Builds AST using recursive descent parsing
- **Evaluator**: Tree-walking interpreter with environment chains
- **REPL**: Interactive development environment
- **Object System**: Runtime value representation

---

## 🤝 Contributing

We welcome contributions to GoKid! Here's how you can help:

### Getting Started
1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Commit your changes: `git commit -am "Add feature"`
7. Push to the branch: `git push origin feature-name`
8. Submit a pull request

### Areas for Contribution
- 🔧 **Language Features**: For loops, switch statements, try/catch
- 📚 **Standard Library**: More built-in functions and utilities  
- 🎨 **Tooling**: Syntax highlighting, IDE integration
- 📖 **Documentation**: Examples, tutorials, API docs
- 🧪 **Testing**: More comprehensive test coverage
- 🚀 **Performance**: Optimization and benchmarking

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gokid.git
cd gokid

# Run tests
go test ./...

# Run the test suite
go run main.go run examples/test.gokid

# Start development REPL
go run main.go repl
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🎉 Acknowledgments

- Inspired by "Writing An Interpreter In Go" by Thorsten Ball
- Thanks to the Go community for excellent documentation and tools
- Special thanks to all contributors and users of GoKid

---

## 📞 Contact

- **Author**: [xspoilt-dev](https://github.com/xspoilt-dev)
- **Repository**: [https://github.com/xspoilt-dev/gokid](https://github.com/xspoilt-dev/gokid)
- **Issues**: [Report a bug or request a feature](https://github.com/xspoilt-dev/gokid/issues)

---

<div align="center">

**⭐ Star this repository if you find GoKid useful! ⭐**

Made with ❤️ by [xspoilt-dev](https://github.com/xspoilt-dev)

</div>