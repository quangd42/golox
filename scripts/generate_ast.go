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
	"Call: callee expr, paren token, arguments []expr",
	"Get: object expr, name token",
	"Grouping: expr expr",
	"Literal: value any",
	"Logical: left expr, operator token, right expr",
	"Set: object expr, name token, value expr",
	"This: keyword token",
	"Unary: operator token, right expr",
	"Variable: name token",
	"Assign: name token, value expr",
	"Ternary: condition expr, thenExpr expr, elseExpr expr",
}

var StmtTypes = []string{
	"Expr: expr expr",
	"Function: name token, params []token, body []stmt",
	"If: condition expr, thenBranch stmt, elseBranch stmt",
	"Print: expr expr",
	"Return: keyword token, value expr",
	"Var: name token, initializer expr",
	"While: condition expr, body stmt, label token, increment stmt",
	"For: initializer stmt, whileBody whileStmt",
	"Break: keyword token, label token",
	"Continue: keyword token, label token",
	"Block: statements []stmt",
	"Class: name token, methods []functionStmt",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: gen_ast <output directory>")
		os.Exit(64)
	}
	// filepath.Abs() ?
	outDir := os.Args[1]
	defineAST(outDir, "Expr", "(any, error)", ExprTypes)
	defineAST(outDir, "Stmt", "error", StmtTypes)
}

func defineAST(outDir, baseName, returnStr string, types []string) {
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
	defineBaseInterface(w, lower(baseName), returnStr)
	// visitor interface
	defineVisitorInterface(w, baseName, returnStr, types)
	defineTypes(w, baseName, returnStr, types)
	fmt.Printf("output to %s\n", path)
}

func lower(s string) string {
	return strings.ToLower(s)
}

func defineBaseInterface(w io.Writer, baseName, returnStr string) {
	baseName = lower(baseName)
	fmt.Fprintf(w, "type %s interface {\n", baseName)
	fmt.Fprintf(w, "	accept(visitor %sVisitor) %s\n", baseName, returnStr)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "")
}

func defineVisitorInterface(w io.Writer, baseName, returnStr string, types []string) {
	fmt.Fprintf(w, "type %sVisitor interface {\n", lower(baseName))
	for _, t := range types {
		name, _, found := strings.Cut(t, ":")
		if !found {
			log.Fatalf("invalid ast format %s\n", t)
		}
		fmt.Fprintf(w, "	visit%s%s(e %s%s) %s\n", name, baseName, lower(name), baseName, returnStr)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "")
}

func defineTypes(w io.Writer, baseName, returnStr string, types []string) {
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

		fmt.Fprintf(w, "func (e %s%s) accept(v %sVisitor) %s {\n", lower(typeName), baseName, lower(baseName), returnStr)
		fmt.Fprintf(w, "	return v.visit%s%s(e)\n", typeName, baseName)
		fmt.Fprintln(w, "}")
		fmt.Fprintln(w, "")
	}
}
