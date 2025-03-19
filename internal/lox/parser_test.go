package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_primary(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{desc: "TRUE", input: "true", want: literalExpr{true}},
		{desc: "FALSE", input: "false", want: literalExpr{false}},
		{desc: "NIL", input: "nil", want: literalExpr{nil}},
		{desc: "NUMBER_int", input: "45", want: literalExpr{45}},
		{desc: "NUMBER_float", input: "49.67", want: literalExpr{49.67}},
		{
			desc:  "PAREN",
			input: "(true)",
			want:  groupingExpr{literalExpr{true}},
		},
		{
			desc:  "PAREN_nested",
			input: "((nil))",
			want:  groupingExpr{groupingExpr{literalExpr{nil}}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.primary()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_unary(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "BANG",
			input: "!true",
			want:  unaryExpr{operator: newTokenNoLiteralType(BANG, 1, 0), right: literalExpr{true}},
		},
		{
			desc:  "MINUS",
			input: "-56.19",
			want:  unaryExpr{operator: newTokenNoLiteralType(MINUS, 1, 0), right: literalExpr{56.19}},
		},
		{
			desc:  "NESTED",
			input: "!-!true",
			want: unaryExpr{
				operator: newTokenNoLiteralType(BANG, 1, 0),
				right: unaryExpr{
					operator: newTokenNoLiteralType(MINUS, 1, 1),
					right: unaryExpr{
						operator: newTokenNoLiteralType(BANG, 1, 2),
						right:    literalExpr{true},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.unary()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_factor(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "SLASH",
			input: "12/9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(SLASH, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "STAR",
			input: "12*9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(STAR, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "SLASH_STAR_SLASH",
			input: "12/9*78/6",
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteralType(SLASH, 1, 2),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteralType(STAR, 1, 4),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteralType(SLASH, 1, 7),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.factor()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_term(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "MINUS",
			input: "12-9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(MINUS, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "PLUS",
			input: "12+9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(PLUS, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "MINUS_PLUS_MINUS",
			input: "12-9+78-6",
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteralType(MINUS, 1, 2),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteralType(PLUS, 1, 4),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteralType(MINUS, 1, 7),
				right:    literalExpr{6},
			},
		},
		{
			desc:  "MINUS_STAR_MINUS",
			input: "12-9*78+6",
			want: binaryExpr{
				left: binaryExpr{
					left:     literalExpr{12},
					operator: newTokenNoLiteralType(MINUS, 1, 2),
					right: binaryExpr{
						left:     literalExpr{9},
						operator: newTokenNoLiteralType(STAR, 1, 4),
						right:    literalExpr{78},
					},
				},
				operator: newTokenNoLiteralType(PLUS, 1, 7),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.term()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_comparison(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "GREATER",
			input: "12>9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(GREATER, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "GREATER_EQUAL",
			input: "12>=9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(GREATER_EQUAL, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "LESS",
			input: "12<9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(LESS, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "LESS_EQUAL",
			input: "12<=9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(LESS_EQUAL, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "GREATER__LESS__GREATER_EQUAL",
			input: "12>9<78>=6",
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteralType(GREATER, 1, 2),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteralType(LESS, 1, 4),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteralType(GREATER_EQUAL, 1, 7),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.comparison()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_equality(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "BANG_EQUAL",
			input: "12!=9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(BANG_EQUAL, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "EQUAL_EQUAL",
			input: "12==9",
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteralType(EQUAL_EQUAL, 1, 2), right: literalExpr{9}},
		},
		{
			desc:  "BANG_EQUAL__BANG_EQUAL__EQUAL_EQUAL",
			input: "12!=9!=78==6",
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteralType(BANG_EQUAL, 1, 2),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteralType(BANG_EQUAL, 1, 5),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteralType(EQUAL_EQUAL, 1, 9),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.equality()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_ternary(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "Simple_Ternary",
			input: "23==2.3?true:false",
			want: binaryExpr{
				left: binaryExpr{
					left:     binaryExpr{left: literalExpr{23}, operator: newTokenNoLiteralType(EQUAL_EQUAL, 1, 2), right: literalExpr{2.3}},
					operator: newTokenNoLiteralType(QUESTION, 1, 7),
					right:    literalExpr{true},
				},
				operator: newTokenNoLiteralType(COLON, 1, 12),
				right:    literalExpr{false},
			},
		},
		{
			desc:  "Nested_Ternary",
			input: "10>5?(true?false:true):(false?true:nil)",
			want: binaryExpr{
				left: binaryExpr{
					left:     binaryExpr{left: literalExpr{10}, operator: newTokenNoLiteralType(GREATER, 1, 2), right: literalExpr{5}},
					operator: newTokenNoLiteralType(QUESTION, 1, 4),
					right: groupingExpr{
						binaryExpr{
							left: binaryExpr{
								left:     literalExpr{true},
								operator: newTokenNoLiteralType(QUESTION, 1, 10),
								right:    literalExpr{false},
							},
							operator: newTokenNoLiteralType(COLON, 1, 16),
							right:    literalExpr{true},
						},
					},
				},
				operator: newTokenNoLiteralType(COLON, 1, 22),
				right: groupingExpr{
					binaryExpr{
						left: binaryExpr{
							left:     literalExpr{false},
							operator: newTokenNoLiteralType(QUESTION, 1, 29),
							right:    literalExpr{true},
						},
						operator: newTokenNoLiteralType(COLON, 1, 34),
						right:    literalExpr{nil},
					},
				},
			},
		},
		{
			desc:  "Complex_Condition_Ternary",
			input: "(5+3)<10?\"yes\":\"no\"",
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left: groupingExpr{
							binaryExpr{
								left:     literalExpr{5},
								operator: newTokenNoLiteralType(PLUS, 1, 2),
								right:    literalExpr{3},
							},
						},
						operator: newTokenNoLiteralType(LESS, 1, 5),
						right:    literalExpr{10},
					},
					operator: newTokenNoLiteralType(QUESTION, 1, 8),
					right:    literalExpr{"yes"},
				},
				operator: newTokenNoLiteralType(COLON, 1, 14),
				right:    literalExpr{"no"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.ternary()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_or(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "simple_or",
			input: "true or false",
			want:  logicalExpr{left: literalExpr{true}, operator: newTokenNoLiteralType(OR, 1, 5), right: literalExpr{false}},
		},
		{
			desc:  "chained_or",
			input: "true or false or nil",
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteralType(OR, 1, 5),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteralType(OR, 1, 14),
				right:    literalExpr{nil},
			},
		},
		{
			desc:  "or_with_expressions",
			input: "1+2 or 3*4",
			want: logicalExpr{
				left: binaryExpr{
					left:     literalExpr{1},
					operator: newTokenNoLiteralType(PLUS, 1, 1),
					right:    literalExpr{2},
				},
				operator: newTokenNoLiteralType(OR, 1, 4),
				right: binaryExpr{
					left:     literalExpr{3},
					operator: newTokenNoLiteralType(STAR, 1, 8),
					right:    literalExpr{4},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.or()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_and(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
	}{
		{
			desc:  "simple_and",
			input: "true and false",
			want:  logicalExpr{left: literalExpr{true}, operator: newTokenNoLiteralType(AND, 1, 5), right: literalExpr{false}},
		},
		{
			desc:  "chained_and",
			input: "true and false and nil",
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteralType(AND, 1, 5),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteralType(AND, 1, 15),
				right:    literalExpr{nil},
			},
		},
		{
			desc:  "and_with_expressions",
			input: "1+2 and 3*4",
			want: logicalExpr{
				left: binaryExpr{
					left:     literalExpr{1},
					operator: newTokenNoLiteralType(PLUS, 1, 1),
					right:    literalExpr{2},
				},
				operator: newTokenNoLiteralType(AND, 1, 4),
				right: binaryExpr{
					left:     literalExpr{3},
					operator: newTokenNoLiteralType(STAR, 1, 9),
					right:    literalExpr{4},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.and()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_assignment(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
		err   error
	}{
		{
			desc:  "simple_assignment",
			input: "x=42",
			want: assignExpr{
				name:  newToken(IDENTIFIER, "x", "x", 1, 0),
				value: literalExpr{42},
			},
		},
		{
			desc:  "chained_assignment",
			input: "x=y=42",
			want: assignExpr{
				name: newToken(IDENTIFIER, "x", "x", 1, 0),
				value: assignExpr{
					name:  newToken(IDENTIFIER, "y", "y", 1, 2),
					value: literalExpr{42},
				},
			},
		},
		{
			desc:  "assignment_with_expression",
			input: "x=10+5",
			want: assignExpr{
				name: newToken(IDENTIFIER, "x", "x", 1, 0),
				value: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(PLUS, 1, 4),
					right:    literalExpr{5},
				},
			},
		},
		{
			desc:  "invalid_assignment_target",
			input: "42=10",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(EQUAL, 1, 2), "Invalid assignment target."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.assignment()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_call(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
		err   error
	}{
		{
			desc:  "simple_call",
			input: "say(\"hello\")",
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "say", "say", 1, 0)},
				paren:  newTokenNoLiteralType(RIGHT_PAREN, 1, 11),
				arguments: []expr{
					literalExpr{"hello"},
				},
			},
		},
		{
			desc:  "no_arguments",
			input: "clock()",
			want: callExpr{
				callee:    variableExpr{newToken(IDENTIFIER, "clock", "clock", 1, 0)},
				paren:     newTokenNoLiteralType(RIGHT_PAREN, 1, 6),
				arguments: []expr{},
			},
		},
		{
			desc:  "multiple_arguments",
			input: "sum(1,2,3)",
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "sum", "sum", 1, 0)},
				paren:  newTokenNoLiteralType(RIGHT_PAREN, 1, 9),
				arguments: []expr{
					literalExpr{1},
					literalExpr{2},
					literalExpr{3},
				},
			},
		},
		{
			desc:  "nested_calls",
			input: "outer(inner(42))",
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "outer", "outer", 1, 0)},
				paren:  newTokenNoLiteralType(RIGHT_PAREN, 1, 15),
				arguments: []expr{
					callExpr{
						callee: variableExpr{newToken(IDENTIFIER, "inner", "inner", 1, 6)},
						paren:  newTokenNoLiteralType(RIGHT_PAREN, 1, 14),
						arguments: []expr{
							literalExpr{42},
						},
					},
				},
			},
		},
		{
			desc:  "multiple_consecutive_calls",
			input: "first()()()",
			want: callExpr{
				callee: callExpr{
					callee: callExpr{
						callee:    variableExpr{newToken(IDENTIFIER, "first", "first", 1, 0)},
						paren:     newTokenNoLiteralType(RIGHT_PAREN, 1, 6),
						arguments: []expr{},
					},
					paren:     newTokenNoLiteralType(RIGHT_PAREN, 1, 8),
					arguments: []expr{},
				},
				paren:     newTokenNoLiteralType(RIGHT_PAREN, 1, 10),
				arguments: []expr{},
			},
		},
		{
			desc:  "call_with_expressions",
			input: "calc(1+2,3*4)",
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "calc", "calc", 1, 0)},
				paren:  newTokenNoLiteralType(RIGHT_PAREN, 1, 12),
				arguments: []expr{
					binaryExpr{
						left:     literalExpr{1},
						operator: newTokenNoLiteralType(PLUS, 1, 6),
						right:    literalExpr{2},
					},
					binaryExpr{
						left:     literalExpr{3},
						operator: newTokenNoLiteralType(STAR, 1, 10),
						right:    literalExpr{4},
					},
				},
			},
		},
		{
			desc:  "missing_right_paren",
			input: "say(\"hello\"",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(EOF, 1, 11), "Expect ')' after arguments."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.call()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_expression(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  expr
		err   error
	}{
		{
			desc:  "expr_COMMA_expr_COMMA_expr",
			input: "12>9,78>=6,13.5!=51.3",
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteralType(GREATER, 1, 2),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteralType(COMMA, 1, 4),
					right: binaryExpr{
						left:     literalExpr{78},
						operator: newTokenNoLiteralType(GREATER_EQUAL, 1, 7),
						right:    literalExpr{6},
					},
				},
				operator: newTokenNoLiteralType(COMMA, 1, 10),
				right: binaryExpr{
					left:     literalExpr{13.5},
					operator: newTokenNoLiteralType(BANG_EQUAL, 1, 15),
					right:    literalExpr{51.3},
				},
			},
		},
		{
			desc:  "Missing_left_operand_in_binary",
			input: "/13.5!=51.3",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(SLASH, 1, 0), "Expect left operand."),
		},
		{
			desc:  "or_simple",
			input: "true or false",
			want: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteralType(OR, 1, 5),
				right:    literalExpr{false},
			},
		},
		{
			desc:  "and_simple",
			input: "true and false",
			want: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteralType(AND, 1, 5),
				right:    literalExpr{false},
			},
		},
		{
			desc:  "chained_or",
			input: "true or false or nil",
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteralType(OR, 1, 5),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteralType(OR, 1, 14),
				right:    literalExpr{nil},
			},
		},
		{
			desc:  "chained_and",
			input: "true and false and nil",
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteralType(AND, 1, 5),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteralType(AND, 1, 15),
				right:    literalExpr{nil},
			},
		},
		{
			desc:  "mixed_and_or",
			input: "true and false or nil",
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteralType(AND, 1, 5),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteralType(OR, 1, 15),
				right:    literalExpr{nil},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.expression()
			if err != nil {
				assert.Equal(t, tC.err, err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_printStmt(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "printStmt_Simple",
			input: "print 42;",
			want:  printStmt{expr: literalExpr{42}},
		},
		{
			desc:  "printStmt_with_binaryExpr",
			input: "print 42+8;",
			want: printStmt{
				expr: binaryExpr{
					left:     literalExpr{42},
					operator: newTokenNoLiteralType(PLUS, 1, 8),
					right:    literalExpr{8},
				},
			},
		},
		{
			desc:  "missing_semicolon",
			input: "print 42",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(EOF, 1, 8), "Expect ';' after expression."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.statement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_exprStmt(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "exprStmt",
			input: "42+8;",
			want: exprStmt{
				expr: binaryExpr{
					left:     literalExpr{42},
					operator: newTokenNoLiteralType(PLUS, 1, 2),
					right:    literalExpr{8},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.statement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_blockStmt(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "simple_block",
			input: "{print 42;}",
			want: blockStmt{
				statements: []stmt{
					printStmt{expr: literalExpr{42}},
				},
			},
		},
		{
			desc:  "empty_block",
			input: "{}",
			want:  blockStmt{statements: []stmt{}},
		},
		{
			desc:  "nested_blocks",
			input: "{print 1;{print 2;}}",
			want: blockStmt{
				statements: []stmt{
					printStmt{expr: literalExpr{1}},
					blockStmt{
						statements: []stmt{
							printStmt{expr: literalExpr{2}},
						},
					},
				},
			},
		},
		{
			desc:  "block_with_declarations",
			input: "{var x = 10;print x;}",
			want: blockStmt{
				statements: []stmt{
					varStmt{
						name:        newToken(IDENTIFIER, "x", "x", 1, 5),
						initializer: literalExpr{10},
					},
					printStmt{
						expr: variableExpr{newToken(IDENTIFIER, "x", "x", 1, 18)},
					},
				},
			},
		},
		{
			desc:  "missing_right_brace",
			input: "{print 42;",
			want:  nil,
			err:   NewParseError(newToken(EOF, "", nil, 1, 10), "Expect '}' after block."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.statement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_ifStmt(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "simple_if_then",
			input: "if 10>5 {print true;}",
			want: ifStmt{
				condition: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(GREATER, 1, 5),
					right:    literalExpr{5},
				},
				thenBranch: blockStmt{
					statements: []stmt{
						printStmt{expr: literalExpr{true}},
					},
				},
				elseBranch: nil,
			},
		},
		{
			desc:  "simple_if_then_else",
			input: "if 10>5 {print true;} else {print false;}",
			want: ifStmt{
				condition: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(GREATER, 1, 5),
					right:    literalExpr{5},
				},
				thenBranch: blockStmt{
					statements: []stmt{
						printStmt{expr: literalExpr{true}},
					},
				},
				elseBranch: blockStmt{
					statements: []stmt{
						printStmt{expr: literalExpr{false}},
					},
				},
			},
		},
		{
			desc:  "if_with_parentheses",
			input: "if (10>5) {print true;}",
			want: ifStmt{
				condition: groupingExpr{binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(GREATER, 1, 6),
					right:    literalExpr{5},
				}},
				thenBranch: blockStmt{statements: []stmt{printStmt{expr: literalExpr{true}}}},
				elseBranch: nil,
			},
		},
		{
			desc:  "if_missing_block",
			input: "if true print true;",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(PRINT, 1, 8), "Expect block."),
		},
		{
			desc:  "else_missing_block",
			input: "if true {print true;} else print true;",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(PRINT, 1, 27), "Expect block."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.statement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_whileStmt(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "simple_while",
			input: "while true {print 42;}",
			want: whileStmt{
				condition: literalExpr{true},
				body: blockStmt{
					statements: []stmt{
						printStmt{expr: literalExpr{42}},
					},
				},
			},
		},
		{
			desc:  "while_complex_condition",
			input: "while 10>5 {print true;}",
			want: whileStmt{
				condition: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(GREATER, 1, 8),
					right:    literalExpr{5},
				},
				body: blockStmt{
					statements: []stmt{
						printStmt{expr: literalExpr{true}},
					},
				},
			},
		},
		{
			desc:  "while_missing_block",
			input: "while true print true;",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(PRINT, 1, 11), "Expect block."),
		},
		{
			desc:  "while_with_parentheses",
			input: "while (10>5) {print true;}",
			want: whileStmt{
				condition: groupingExpr{binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(GREATER, 1, 9),
					right:    literalExpr{5},
				}},
				body: blockStmt{statements: []stmt{printStmt{expr: literalExpr{true}}}},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.statement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_forStatement(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "basic_for_loop",
			input: "for i = 0; i < 10; i = i + 1 {print i;}",
			want: blockStmt{
				statements: []stmt{
					exprStmt{expr: assignExpr{
						name:  newToken(IDENTIFIER, "i", "i", 1, 4),
						value: literalExpr{0},
					}},
					whileStmt{
						condition: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "i", "i", 1, 11)},
							operator: newTokenNoLiteralType(LESS, 1, 13),
							right:    literalExpr{10},
						},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 1, 36)}},
								exprStmt{expr: assignExpr{
									name: newToken(IDENTIFIER, "i", "i", 1, 19),
									value: binaryExpr{
										left:     variableExpr{newToken(IDENTIFIER, "i", "i", 1, 23)},
										operator: newTokenNoLiteralType(PLUS, 1, 25),
										right:    literalExpr{1},
									},
								}},
							},
						},
					},
				},
			},
		},
		{
			desc:  "for_loop_without_initializer",
			input: "for ; x < 5; x = x + 1 {print x;}",
			want: whileStmt{
				condition: binaryExpr{
					left:     variableExpr{newToken(IDENTIFIER, "x", "x", 1, 6)},
					operator: newTokenNoLiteralType(LESS, 1, 8),
					right:    literalExpr{5},
				},
				body: blockStmt{
					statements: []stmt{
						printStmt{expr: variableExpr{newToken(IDENTIFIER, "x", "x", 1, 30)}},
						exprStmt{expr: assignExpr{
							name: newToken(IDENTIFIER, "x", "x", 1, 13),
							value: binaryExpr{
								left:     variableExpr{newToken(IDENTIFIER, "x", "x", 1, 17)},
								operator: newTokenNoLiteralType(PLUS, 1, 19),
								right:    literalExpr{1},
							},
						}},
					},
				},
			},
		},
		{
			desc:  "for_loop_without_condition",
			input: "for i = 0; ; i = i + 1 {print i;}",
			want: blockStmt{
				statements: []stmt{
					exprStmt{expr: assignExpr{
						name:  newToken(IDENTIFIER, "i", "i", 1, 4),
						value: literalExpr{0},
					}},
					whileStmt{
						condition: literalExpr{true},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 1, 30)}},
								exprStmt{expr: assignExpr{
									name: newToken(IDENTIFIER, "i", "i", 1, 13),
									value: binaryExpr{
										left:     variableExpr{newToken(IDENTIFIER, "i", "i", 1, 17)},
										operator: newTokenNoLiteralType(PLUS, 1, 19),
										right:    literalExpr{1},
									},
								}},
							},
						},
					},
				},
			},
		},
		{
			desc:  "for_loop_without_increment",
			input: "for i = 0; i < 10; {print i;}",
			want: blockStmt{
				statements: []stmt{
					exprStmt{expr: assignExpr{
						name:  newToken(IDENTIFIER, "i", "i", 1, 4),
						value: literalExpr{0},
					}},
					whileStmt{
						condition: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "i", "i", 1, 11)},
							operator: newTokenNoLiteralType(LESS, 1, 13),
							right:    literalExpr{10},
						},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 1, 26)}},
							},
						},
					},
				},
			},
		},
		{
			desc:  "for_loop_with_var_declaration",
			input: "for var i = 0; i < 5; i = i + 1 {print i;}",
			want: blockStmt{
				statements: []stmt{
					varStmt{
						name:        newToken(IDENTIFIER, "i", "i", 1, 8),
						initializer: literalExpr{0},
					},
					whileStmt{
						condition: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "i", "i", 1, 15)},
							operator: newTokenNoLiteralType(LESS, 1, 17),
							right:    literalExpr{5},
						},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 1, 39)}},
								exprStmt{expr: assignExpr{
									name: newToken(IDENTIFIER, "i", "i", 1, 22),
									value: binaryExpr{
										left:     variableExpr{newToken(IDENTIFIER, "i", "i", 1, 26)},
										operator: newTokenNoLiteralType(PLUS, 1, 28),
										right:    literalExpr{1},
									},
								}},
							},
						},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.statement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_declaration(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "var_declaration_no_initializer",
			input: "var foo;",
			want: varStmt{
				name:        newToken(IDENTIFIER, "foo", "foo", 1, 4),
				initializer: nil,
			},
		},
		{
			desc:  "var_declaration_with_initializer",
			input: "var foo = 42;",
			want: varStmt{
				name:        newToken(IDENTIFIER, "foo", "foo", 1, 4),
				initializer: literalExpr{42},
			},
		},
		{
			desc:  "missing_semicolon",
			input: "var foo",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(EOF, 1, 7), "Expect ';' after variable declaration."),
		},
		{
			desc:  "missing_identifier",
			input: "var;",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(SEMICOLON, 1, 3), "Expect variable name."),
		},
		{
			desc:  "expr_statement_fallback",
			input: "42;",
			want:  exprStmt{expr: literalExpr{42}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.declaration()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_function(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "function_declaration",
			input: "fn greet(name) {print \"Hello, \" + name;}",
			want: functionStmt{
				name: newToken(IDENTIFIER, "greet", "greet", 1, 3),
				params: []token{
					newToken(IDENTIFIER, "name", "name", 1, 9),
				},
				body: []stmt{
					printStmt{
						expr: binaryExpr{
							left:     literalExpr{"Hello, "},
							operator: newTokenNoLiteralType(PLUS, 1, 32),
							right:    variableExpr{newToken(IDENTIFIER, "name", "name", 1, 34)},
						},
					},
				},
			},
		},
		{
			desc:  "function_no_params",
			input: "fn hello() {print \"Hello!\";}",
			want: functionStmt{
				name:   newToken(IDENTIFIER, "hello", "hello", 1, 3),
				params: []token{},
				body: []stmt{
					printStmt{expr: literalExpr{"Hello!"}},
				},
			},
		},
		{
			desc:  "function_multiple_params",
			input: "fn add(a,b) {print a + b;}",
			want: functionStmt{
				name: newToken(IDENTIFIER, "add", "add", 1, 3),
				params: []token{
					newToken(IDENTIFIER, "a", "a", 1, 7),
					newToken(IDENTIFIER, "b", "b", 1, 9),
				},
				body: []stmt{
					printStmt{
						expr: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "a", "a", 1, 19)},
							operator: newTokenNoLiteralType(PLUS, 1, 21),
							right:    variableExpr{newToken(IDENTIFIER, "b", "b", 1, 23)},
						},
					},
				},
			},
		},
		{
			desc:  "missing_function_name",
			input: "fn() {}",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(LEFT_PAREN, 1, 2), "Expect function name."),
		},
		{
			desc:  "missing_left_paren",
			input: "fn test param) {}",
			want:  nil,
			err:   NewParseError(newToken(IDENTIFIER, "param", "param", 1, 8), "Expect '(' after function name."),
		},
		{
			desc:  "missing_right_paren",
			input: "fn test(param {}",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(LEFT_BRACE, 1, 14), "Expect ')' after parameters."),
		},
		{
			desc:  "missing_body",
			input: "fn test()",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(EOF, 1, 9), "Expect block."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.function("function")
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_returnStatement(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  stmt
		err   error
	}{
		{
			desc:  "return_with_value",
			input: "return 42;",
			want: returnStmt{
				keyword: newTokenNoLiteralType(RETURN, 1, 0),
				value:   literalExpr{42},
			},
		},
		{
			desc:  "return_without_value",
			input: "return;",
			want: returnStmt{
				keyword: newTokenNoLiteralType(RETURN, 1, 0),
				value:   nil,
			},
		},
		{
			desc:  "return_with_expression",
			input: "return 10 + 5;",
			want: returnStmt{
				keyword: newTokenNoLiteralType(RETURN, 1, 0),
				value: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteralType(PLUS, 1, 10),
					right:    literalExpr{5},
				},
			},
		},
		{
			desc:  "missing_semicolon",
			input: "return 42",
			want:  nil,
			err:   NewParseError(newTokenNoLiteralType(EOF, 1, 9), "Expect ';' after return value."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.returnStatement()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_Parse(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  []stmt
		err   error
	}{
		{
			desc:  "single_print_statement",
			input: "print 42;",
			want:  []stmt{printStmt{expr: literalExpr{42}}},
		},
		{
			desc:  "multiple_statements",
			input: "var foo = 42;\nprint foo;",
			want: []stmt{
				varStmt{
					name:        newToken(IDENTIFIER, "foo", "foo", 1, 4),
					initializer: literalExpr{42},
				},
				printStmt{expr: variableExpr{newToken(IDENTIFIER, "foo", "foo", 2, 20)}},
			},
		},
		{
			desc:  "variable_with_ternary",
			input: "var result = 10 > 5 ? \"yes\" : \"no\";",
			want: []stmt{
				varStmt{
					name: newToken(IDENTIFIER, "result", "result", 1, 4),
					initializer: binaryExpr{
						left: binaryExpr{
							left: binaryExpr{
								left:     literalExpr{10},
								operator: newToken(GREATER, ">", ">", 1, 16),
								right:    literalExpr{5},
							},
							operator: newToken(QUESTION, "?", "?", 1, 20),
							right:    literalExpr{"yes"},
						},
						operator: newToken(COLON, ":", ":", 1, 28),
						right:    literalExpr{"no"},
					},
				},
			},
		},
		{
			desc:  "empty_input",
			input: "",
			want:  []stmt{},
		},
		{
			desc:  "parse_error",
			input: "print 42",
			want:  []stmt{},
		},
		{
			desc:  "synchronize_recovers_at_statement_boundary",
			input: "var 42\nprint \"hello\";",
			want: []stmt{
				printStmt{expr: literalExpr{"hello"}},
			},
			err: nil, // Synchronize should allow parsing to continue after error
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner([]byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			parser := NewParser(tokens)
			got, err := parser.Parse()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}
