package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_interpretLiteral(t *testing.T) {
	testCases := []struct {
		desc  string
		input literalExpr
		want  any
		err   error
	}{
		{
			desc:  "STRING",
			input: literalExpr{"a string"},
			want:  "a string",
			err:   nil,
		},
		{
			desc:  "NUMBER_float64",
			input: literalExpr{158.2},
			want:  158.2,
			err:   nil,
		},
		{
			desc:  "NUMBER_int",
			input: literalExpr{2389},
			want:  2389,
			err:   nil,
		},
		{
			desc:  "NIL",
			input: literalExpr{nil},
			want:  nil,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			got, err := interpreter.visitLiteralExpr(tC.input)
			assert.Equal(t, tC.want, got)
			assert.Equal(t, tC.err, err)
		})
	}
}

func Test_interpretUnary(t *testing.T) {
	minus := newTokenNoLiteral(MINUS)
	bang := newTokenNoLiteral(BANG)
	testCases := []struct {
		desc  string
		input unaryExpr
		want  any
		err   error
	}{
		{
			desc:  "MINUS__NUMBER__Float",
			input: unaryExpr{operator: minus, right: literalExpr{189.228}},
			want:  -189.228,
			err:   nil,
		},
		{
			desc:  "MINUS__NUMBER__Int",
			input: unaryExpr{operator: minus, right: literalExpr{189}},
			want:  float64(-189),
			err:   nil,
		},
		{
			desc:  "MINUS__NUMBER__NaN",
			input: unaryExpr{operator: minus, right: literalExpr{"NaN"}},
			want:  nil,
			err:   NewRuntimeError(minus, "Operand must be a number."),
		},
		{
			desc:  "BANG__TRUE",
			input: unaryExpr{operator: bang, right: literalExpr{true}},
			want:  false,
			err:   nil,
		},
		{
			desc:  "BANG__FALSE",
			input: unaryExpr{operator: bang, right: literalExpr{false}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "BANG__NIL",
			input: unaryExpr{operator: bang, right: literalExpr{nil}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "BANG__LITERAL",
			input: unaryExpr{operator: bang, right: literalExpr{"some string"}},
			want:  false,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			got, err := interpreter.visitUnaryExpr(tC.input)
			assert.Equal(t, tC.want, got)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}
		})
	}
}

func Test_interpretBinary(t *testing.T) {
	plus := newTokenNoLiteral(PLUS)
	minus := newTokenNoLiteral(MINUS)
	star := newTokenNoLiteral(STAR)
	slash := newTokenNoLiteral(SLASH)
	greater := newTokenNoLiteral(GREATER)
	greaterEqual := newTokenNoLiteral(GREATER_EQUAL)
	less := newTokenNoLiteral(LESS)
	lessEqual := newTokenNoLiteral(LESS_EQUAL)
	equalEqual := newTokenNoLiteral(EQUAL_EQUAL)
	bangEqual := newTokenNoLiteral(BANG_EQUAL)

	testCases := []struct {
		desc  string
		input binaryExpr
		want  any
		err   error
	}{
		{
			desc:  "PLUS_float_float",
			input: binaryExpr{left: literalExpr{5.0}, operator: plus, right: literalExpr{3.0}},
			want:  8.0,
			err:   nil,
		},
		{
			desc:  "PLUS_int_int",
			input: binaryExpr{left: literalExpr{5}, operator: plus, right: literalExpr{3}},
			want:  8.0,
			err:   nil,
		},
		{
			desc:  "PLUS_float_int",
			input: binaryExpr{left: literalExpr{5.0}, operator: plus, right: literalExpr{3}},
			want:  8.0,
			err:   nil,
		},
		{
			desc:  "PLUS_strings",
			input: binaryExpr{left: literalExpr{"hello"}, operator: plus, right: literalExpr{" world"}},
			want:  "hello world",
			err:   nil,
		},
		{
			desc:  "PLUS_invalid",
			input: binaryExpr{left: literalExpr{true}, operator: plus, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(plus, "Operands must be either numbers or strings."),
		},
		{
			desc:  "MINUS",
			input: binaryExpr{left: literalExpr{5.0}, operator: minus, right: literalExpr{3.0}},
			want:  2.0,
			err:   nil,
		},
		{
			desc:  "MINUS_invalid",
			input: binaryExpr{left: literalExpr{"string"}, operator: minus, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(minus, "Operands must be numbers."),
		},
		{
			desc:  "MULTIPLY",
			input: binaryExpr{left: literalExpr{5.0}, operator: star, right: literalExpr{3.0}},
			want:  15.0,
			err:   nil,
		},
		{
			desc:  "MULTIPLY_invalid",
			input: binaryExpr{left: literalExpr{true}, operator: star, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(star, "Operands must be numbers."),
		},
		{
			desc:  "DIVIDE",
			input: binaryExpr{left: literalExpr{15.0}, operator: slash, right: literalExpr{3.0}},
			want:  5.0,
			err:   nil,
		},
		{
			desc:  "DIVIDE_invalid",
			input: binaryExpr{left: literalExpr{"string"}, operator: slash, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(slash, "Operands must be numbers."),
		},
		{
			desc:  "GREATER",
			input: binaryExpr{left: literalExpr{5.0}, operator: greater, right: literalExpr{3.0}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "GREATER_invalid",
			input: binaryExpr{left: literalExpr{true}, operator: greater, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(greater, "Operands must be numbers."),
		},
		{
			desc:  "GREATER_EQUAL",
			input: binaryExpr{left: literalExpr{5.0}, operator: greaterEqual, right: literalExpr{5.0}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "GREATER_EQUAL_invalid",
			input: binaryExpr{left: literalExpr{"string"}, operator: greaterEqual, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(greaterEqual, "Operands must be numbers."),
		},
		{
			desc:  "LESS",
			input: binaryExpr{left: literalExpr{3.0}, operator: less, right: literalExpr{5.0}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "LESS_invalid",
			input: binaryExpr{left: literalExpr{true}, operator: less, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(less, "Operands must be numbers."),
		},
		{
			desc:  "LESS_EQUAL",
			input: binaryExpr{left: literalExpr{5.0}, operator: lessEqual, right: literalExpr{5.0}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "LESS_EQUAL_invalid",
			input: binaryExpr{left: literalExpr{"string"}, operator: lessEqual, right: literalExpr{5.0}},
			want:  nil,
			err:   NewRuntimeError(lessEqual, "Operands must be numbers."),
		},
		{
			desc:  "EQUAL_EQUAL",
			input: binaryExpr{left: literalExpr{5.0}, operator: equalEqual, right: literalExpr{5.0}},
			want:  true,
			err:   nil,
		},
		{
			desc:  "BANG_EQUAL",
			input: binaryExpr{left: literalExpr{5.0}, operator: bangEqual, right: literalExpr{3.0}},
			want:  true,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			got, err := interpreter.visitBinaryExpr(tC.input)
			assert.Equal(t, tC.want, got)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}
		})
	}
}
