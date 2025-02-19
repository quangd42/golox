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
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(filename string) {
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

	err = run(b)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(65)
	}
}

func runPrompt() {
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

			continue // If err != nil delim is missing from input, so keep scanning for more
		}

		err = run(input[:len(input)-1]) // Remove delim \n from input before running
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(65)
		}
	}
}

func run(source []byte) error {
	scanner := lox.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", tokens)
	return nil
}
