package main

import (
	"os"

	"github.com/quangd42/golox/internal/lox"
)

func main() {
	if len(os.Args) > 2 {
		println("Usage: golox [script]")
		os.Exit(64)
	}

	er := lox.NewLoxErrorReporter()
	i := lox.NewInterpreter(er)
	r := lox.NewResolver(er, i)
	runtime := lox.NewRuntime(er, i, r)

	if len(os.Args) == 2 {
		runtime.RunFile(os.Args[1])
	} else {
		runtime.RunPrompt()
	}
}
