package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			os.Exit(65)
		}
	} else {
		runPrompt()
	}
}

func runFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	run(b)
	return nil
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
			log.Fatalf("input error %v: ", err)
		}

		run(input[:len(input)-1]) // Remove \n from input before running
	}
}

// TODO: this is the core processor
func run(source []byte) {
	fmt.Print(string(source))
}

func report(line int, where, msg string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, msg)
}
