package lox

import (
	"testing"
)

func TestAstPrinter(t *testing.T) {
	expr := binaryExpr{
		left: unaryExpr{
			operator: newToken(MINUS, "-", nil, 1),
			right:    literalExpr{value: 123},
		},
		operator: newToken(STAR, "*", nil, 1),
		right: groupingExpr{
			expr: literalExpr{value: 45.67},
		},
	}
	printer := astPrinter{}

	t.Run("ast printer", func(t *testing.T) {
		got, err := printer.String(expr)
		if err != nil {
			t.Error(err)
		}
		expected := "(* (- 123) (group 45.67))"
		if got != expected {
			t.Errorf("expected %s, got %s\n", expected, got)
		}
	})
}
