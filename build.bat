@echo off
REM build.bat - Build script for GoKid Language Interpreter (Windows)

echo Building GoKid Language Interpreter...

REM Clean any existing builds
if exist gokid.exe del gokid.exe
if exist gokid-linux del gokid-linux
if exist gokid-macos del gokid-macos

REM Build for Windows
echo Building for Windows...
go build -o gokid.exe main.go

REM Build for Linux
echo Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o gokid-linux main.go

REM Build for macOS
echo Building for macOS...
set GOOS=darwin
set GOARCH=amd64
go build -o gokid-macos main.go

REM Reset environment
set GOOS=
set GOARCH=

echo Build complete!
echo.
echo Usage:
echo   gokid.exe examples\demo.gokid    # Run a GoKid program
echo   gokid.exe repl                   # Start interactive REPL
echo   gokid.exe help                   # Show help
