package lox

import (
	"errors"
	"strconv"
	"strings"
)

type astPrinter struct{}

func NewASTPrinter() *astPrinter {
	return &astPrinter{}
}

func (p astPrinter) String(e expr) (any, error) {
	return e.accept(p)
}

func (p astPrinter) visitLiteralExpr(e literalExpr) (any, error) {
	if e.value == nil {
		return "nil", nil
	}
	switch t := e.value.(type) {
	case string:
		return t, nil
	case int:
		return strconv.Itoa(t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', 2, 64), nil
	case bool:
		if t {
			return "true", nil
		}
		return "false", nil
	default:
		return "", errors.New("cannot convert value to string")
	}
}

func (p astPrinter) visitUnaryExpr(e unaryExpr) (any, error) {
	return p.paren(e.operator.lexeme, e.right)
}

func (p astPrinter) visitBinaryExpr(e binaryExpr) (any, error) {
	return p.paren(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitGroupingExpr(e groupingExpr) (any, error) {
	return p.paren("group", e.expr)
}

func (p astPrinter) paren(name string, exprs ...expr) (any, error) {
	var w strings.Builder
	w.WriteString("(" + name)
	for _, e := range exprs {
		e, err := e.accept(p)
		if err != nil {
			return nil, err
		}
		eStr, ok := e.(string)
		if !ok {
			return nil, errors.New("expr does not have string representation")
		}
		w.WriteString(" ")
		w.WriteString(eStr)
	}
	w.WriteString(")")
	return w.String(), nil
}
