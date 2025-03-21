package lox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Runtime struct {
	er ErrorReporter
	i  *Interpreter
	r  *Resolver
}

func NewRuntime(er ErrorReporter, i *Interpreter, r *Resolver) *Runtime {
	return &Runtime{er: er, i: i, r: r}
}

func (rt *Runtime) RunFile(filename string) {
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

	rt.run(b)
	if rt.er.HadError() {
		os.Exit(65)
	}
	if rt.er.HadRuntimeError() {
		os.Exit(70)
	}
}

func (rt *Runtime) RunPrompt() {
	var stdout io.Writer = os.Stdout
	fmt.Fprint(stdout, "Golox 0.02\n")
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
		rt.run(input[:len(input)-1])
		// If there is any error just continue, error will be reported somewhere else
		rt.er.ResetError()
	}
}

func (rt *Runtime) run(source []byte) {
	// Static errors are printed but operation is continued,
	// so only fatal errors are returned.
	scanner := NewScanner(rt.er, source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		log.Fatal("Scanner Error: ", err.Error())
	}
	if rt.er.HadError() {
		return
	}
	parser := NewParser(rt.er, tokens)
	stmts, err := parser.Parse()
	if err != nil {
		log.Fatal("Parser Error: ", err.Error())
	}
	if rt.er.HadError() {
		return
	}
	err = rt.r.Resolve(stmts)
	if err != nil {
		log.Fatal("Resolver Error: ", err.Error())
	}
	if rt.er.HadError() {
		return
	}

	err = rt.i.Interpret(stmts)
	if err != nil {
		log.Fatal("Interpreter Error: ", err.Error())
	}
	if rt.er.HadRuntimeError() {
		return
	}
}
