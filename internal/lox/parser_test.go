package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_primary(t *testing.T) {
	testCases := []struct {
		desc  string
		input []token
		want  expr
	}{
		{desc: "TRUE", input: []token{newTokenNoLiteral(TRUE)}, want: literalExpr{true}},
		{desc: "FALSE", input: []token{newTokenNoLiteral(FALSE)}, want: literalExpr{false}},
		{desc: "NIL", input: []token{newTokenNoLiteral(NIL)}, want: literalExpr{nil}},
		{desc: "NUMBER_int", input: []token{newToken(NUMBER, "45", 45, 0)}, want: literalExpr{45}},
		{desc: "NUMBER_float", input: []token{newToken(NUMBER, "49.67", 49.67, 0)}, want: literalExpr{49.67}},
		{
			desc: "PAREN", input: []token{
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: groupingExpr{literalExpr{true}},
		},
		{
			desc: "PAREN_nested", input: []token{
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(NIL),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: groupingExpr{groupingExpr{literalExpr{nil}}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
			got, err := parser.primary()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}
