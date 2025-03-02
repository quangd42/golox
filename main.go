package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/quangd42/golox/internal/lox"
)

func main() {
	if len(os.Args) > 2 {
		println("Usage: golox [script]")
		os.Exit(64)
	}

	interpreter := lox.NewInterpreter()
	if len(os.Args) == 2 {
		runFile(interpreter, os.Args[1])
	} else {
		runPrompt(interpreter)
	}
}

func runFile(i *lox.Interpreter, filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("can't open file '%s': %v\n", filename, err)
		os.Exit(2)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("can't read file '%s': %v\n", filename, err)
		os.Exit(2)
	}

	err = run(i, b)
	if err != nil {
		var rtErr *lox.RuntimeError
		fmt.Printf("%v\n", err)
		if errors.As(err, &rtErr) {
			os.Exit(70)
		}
		os.Exit(65)
	}
}

func runPrompt(i *lox.Interpreter) {
	var stdout io.Writer = os.Stdout
	fmt.Fprint(stdout, "Golox 0.01\nType \"help\" or something\n")
	for {
		fmt.Fprint(stdout, ">> ")

		// Wait for user input
		stdin := bufio.NewReader(os.Stdin)
		input, err := stdin.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				os.Exit(0)
			}

			// If err != nil delim is missing from input, so keep scanning for more
			continue
		}

		// Remove delim \n from input before running
		// If there is any error just continue, error will be reported somewhere else
		run(i, input[:len(input)-1])
	}
}

func run(i *lox.Interpreter, source []byte) error {
	scanner := lox.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}
	parser := lox.NewParser(tokens)
	stmts, err := parser.Parse()
	if err != nil {
		return err
	}

	return i.Interpret(stmts)
}
