package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ExprTypes = []string{
	"Binary: left expr, operator token, right expr",
	"Grouping: expr expr",
	"Literal: value any",
	"Unary: operator token, right expr",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: gen_ast <output directory>")
		os.Exit(64)
	}
	// filepath.Abs() ?
	outDir := os.Args[1]
	defineAST(outDir, "Expr", ExprTypes)
}

func defineAST(outDir, baseName string, types []string) {
	path := outDir + "/" + lower(baseName) + ".go"
	path = filepath.Clean(path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to create %s\n", path)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	w.WriteString("// Generated with AST generator.\n\n")
	// file header
	w.WriteString("package lox\n\n")
	// main interface
	defineBaseInterface(w, "Expr")
	// visitor interface
	defineVisitorInterface(w, baseName, types)
	defineTypes(w, baseName, types)
	fmt.Printf("output to %s\n", path)
}

func lower(s string) string {
	return strings.ToLower(s)
}

func title(s string) string {
	return strings.ToTitle(s)
}

func defineBaseInterface(w io.Writer, baseName string) {
	baseName = lower(baseName)
	fmt.Fprintf(w, "type %s interface {\n", baseName)
	fmt.Fprintf(w, "	accept(visitor %sVisitor) (any, error)\n", baseName)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "")
}

func defineVisitorInterface(w io.Writer, baseName string, types []string) {
	fmt.Fprintf(w, "type %sVisitor interface {\n", lower(baseName))
	for _, t := range types {
		name, _, found := strings.Cut(t, ":")
		if !found {
			log.Fatalf("invalid ast format %s\n", t)
		}
		fmt.Fprintf(w, "	visit%s%s(e %s%s) (any, error)\n", name, baseName, lower(name), baseName)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "")
}

func defineTypes(w io.Writer, baseName string, types []string) {
	for _, t := range types {
		typeName, fieldsStr, found := strings.Cut(t, ":")
		if !found {
			log.Fatalf("invalid ast format %s\n", t)
		}
		fields := strings.Split(fieldsStr, ", ")
		fmt.Fprintf(w, "type %s%s struct {\n", lower(typeName), baseName)
		for _, f := range fields {
			fmt.Fprintf(w, "	%s\n", strings.TrimSpace(f))
		}
		fmt.Fprintln(w, "}")
		fmt.Fprintln(w, "")

		fmt.Fprintf(w, "func (e %s%s) accept(v %sVisitor) (any, error) {\n", lower(typeName), baseName, lower(baseName))
		fmt.Fprintf(w, "	return v.visit%s%s(e)\n", typeName, baseName)
		fmt.Fprintln(w, "}")
		fmt.Fprintln(w, "")
	}
}
