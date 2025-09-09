#!/bin/bash
# build.sh - Build script for GoKid Language Interpreter

echo "Building GoKid Language Interpreter..."

# Clean any existing builds
rm -f gokid gokid.exe

# Build for current platform
echo "Building for current platform..."
go build -o gokid main.go

# Build for Windows (if not on Windows)
if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "win32" ]]; then
    echo "Building for Windows..."
    GOOS=windows GOARCH=amd64 go build -o gokid.exe main.go
fi

# Build for Linux (if not on Linux)
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo "Building for Linux..."
    GOOS=linux GOARCH=amd64 go build -o gokid-linux main.go
fi

# Build for macOS (if not on macOS)
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "Building for macOS..."
    GOOS=darwin GOARCH=amd64 go build -o gokid-macos main.go
fi

echo "Build complete!"
echo ""
echo "Usage:"
echo "  ./gokid examples/demo.gokid    # Run a GoKid program"
echo "  ./gokid repl                   # Start interactive REPL"
echo "  ./gokid help                   # Show help"
