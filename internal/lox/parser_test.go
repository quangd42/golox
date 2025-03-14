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

func Test_unary(t *testing.T) {
	testCases := []struct {
		desc  string
		input []token
		want  expr
	}{
		{
			desc:  "BANG",
			input: []token{newTokenNoLiteral(BANG), newTokenNoLiteral(TRUE)},
			want:  unaryExpr{operator: newTokenNoLiteral(BANG), right: literalExpr{true}},
		},
		{
			desc:  "MINUS",
			input: []token{newTokenNoLiteral(MINUS), newToken(NUMBER, "56.19", 56.19, 0)},
			want:  unaryExpr{operator: newTokenNoLiteral(MINUS), right: literalExpr{56.19}},
		},
		{
			desc:  "NESTED",
			input: []token{newTokenNoLiteral(BANG), newTokenNoLiteral(MINUS), newTokenNoLiteral(BANG), newTokenNoLiteral(TRUE)},
			want: unaryExpr{
				operator: newTokenNoLiteral(BANG),
				right: unaryExpr{
					operator: newTokenNoLiteral(MINUS),
					right: unaryExpr{
						operator: newTokenNoLiteral(BANG),
						right:    literalExpr{true},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc:  "SLASH",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(SLASH), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(SLASH), right: literalExpr{9}},
		},
		{
			desc:  "STAR",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(STAR), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(STAR), right: literalExpr{9}},
		},
		{
			desc: "SLASH_STAR_SLASH",
			input: []token{
				newToken(NUMBER, "12", 12, 0),
				newTokenNoLiteral(SLASH),
				newToken(NUMBER, "9", 9, 0),
				newTokenNoLiteral(STAR),
				newToken(NUMBER, "78", 78, 0),
				newTokenNoLiteral(SLASH),
				newToken(NUMBER, "6", 6, 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteral(SLASH),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteral(STAR),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteral(SLASH),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc:  "MINUS",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(MINUS), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(MINUS), right: literalExpr{9}},
		},
		{
			desc:  "PLUS",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(PLUS), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(PLUS), right: literalExpr{9}},
		},
		{
			desc: "MINUS_PLUS_MINUS",
			input: []token{
				newToken(NUMBER, "12", 12, 0),
				newTokenNoLiteral(MINUS),
				newToken(NUMBER, "9", 9, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "78", 78, 0),
				newTokenNoLiteral(MINUS),
				newToken(NUMBER, "6", 6, 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteral(MINUS),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteral(MINUS),
				right:    literalExpr{6},
			},
		},
		{
			desc: "MINUS_STAR_MINUS",
			input: []token{
				newToken(NUMBER, "12", 12, 0),
				newTokenNoLiteral(MINUS),
				newToken(NUMBER, "9", 9, 0),
				newTokenNoLiteral(STAR),
				newToken(NUMBER, "78", 78, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "6", 6, 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left:     literalExpr{12},
					operator: newTokenNoLiteral(MINUS),
					right: binaryExpr{
						left:     literalExpr{9},
						operator: newTokenNoLiteral(STAR),
						right:    literalExpr{78},
					},
				},
				operator: newTokenNoLiteral(PLUS),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc:  "GREATER",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(GREATER), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(GREATER), right: literalExpr{9}},
		},
		{
			desc:  "GREATER_EQUAL",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(GREATER_EQUAL), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(GREATER_EQUAL), right: literalExpr{9}},
		},
		{
			desc:  "LESS",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(LESS), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(LESS), right: literalExpr{9}},
		},
		{
			desc:  "LESS_EQUAL",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(LESS_EQUAL), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(LESS_EQUAL), right: literalExpr{9}},
		},
		{
			desc: "GREATER__LESS__GREATER_EQUAL",
			input: []token{
				newToken(NUMBER, "12", 12, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "9", 9, 0),
				newTokenNoLiteral(LESS),
				newToken(NUMBER, "78", 78, 0),
				newTokenNoLiteral(GREATER_EQUAL),
				newToken(NUMBER, "6", 6, 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteral(GREATER),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteral(LESS),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteral(GREATER_EQUAL),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc:  "BANG_EQUAL",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(BANG_EQUAL), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(BANG_EQUAL), right: literalExpr{9}},
		},
		{
			desc:  "EQUAL_EQUAL",
			input: []token{newToken(NUMBER, "12", 12, 0), newTokenNoLiteral(EQUAL_EQUAL), newToken(NUMBER, "9", 9, 0)},
			want:  binaryExpr{left: literalExpr{12}, operator: newTokenNoLiteral(EQUAL_EQUAL), right: literalExpr{9}},
		},
		{
			desc: "BANG_EQUAL__BANG_EQUAL__EQUAL_EQUAL",
			input: []token{
				newToken(NUMBER, "12", 12, 0),
				newTokenNoLiteral(BANG_EQUAL),
				newToken(NUMBER, "9", 9, 0),
				newTokenNoLiteral(BANG_EQUAL),
				newToken(NUMBER, "78", 78, 0),
				newTokenNoLiteral(EQUAL_EQUAL),
				newToken(NUMBER, "6", 6, 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteral(BANG_EQUAL),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteral(BANG_EQUAL),
					right:    literalExpr{78},
				},
				operator: newTokenNoLiteral(EQUAL_EQUAL),
				right:    literalExpr{6},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc: "Simple_Ternary",
			input: []token{
				newToken(NUMBER, "23", 23, 0),
				newTokenNoLiteral(EQUAL_EQUAL),
				newToken(NUMBER, "2.3", 2.3, 0),
				newTokenNoLiteral(QUESTION),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(COLON),
				newTokenNoLiteral(FALSE),
			},
			want: binaryExpr{
				left: binaryExpr{
					left:     binaryExpr{left: literalExpr{23}, operator: newTokenNoLiteral(EQUAL_EQUAL), right: literalExpr{2.3}},
					operator: newTokenNoLiteral(QUESTION),
					right:    literalExpr{true},
				},
				operator: newTokenNoLiteral(COLON),
				right:    literalExpr{false},
			},
		},
		{
			desc: "Nested_Ternary",
			input: []token{
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(QUESTION),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(QUESTION),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(COLON),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(COLON),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(QUESTION),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(COLON),
				newTokenNoLiteral(NIL),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: binaryExpr{
				left: binaryExpr{
					left:     binaryExpr{left: literalExpr{10}, operator: newTokenNoLiteral(GREATER), right: literalExpr{5}},
					operator: newTokenNoLiteral(QUESTION),
					right: groupingExpr{
						binaryExpr{
							left: binaryExpr{
								left:     literalExpr{true},
								operator: newTokenNoLiteral(QUESTION),
								right:    literalExpr{false},
							},
							operator: newTokenNoLiteral(COLON),
							right:    literalExpr{true},
						},
					},
				},
				operator: newTokenNoLiteral(COLON),
				right: groupingExpr{
					binaryExpr{
						left: binaryExpr{
							left:     literalExpr{false},
							operator: newTokenNoLiteral(QUESTION),
							right:    literalExpr{true},
						},
						operator: newTokenNoLiteral(COLON),
						right:    literalExpr{nil},
					},
				},
			},
		},
		{
			desc: "Complex_Condition_Ternary",
			input: []token{
				newTokenNoLiteral(LEFT_PAREN),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "3", 3, 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LESS),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(QUESTION),
				newToken(STRING, "yes", "yes", 0),
				newTokenNoLiteral(COLON),
				newToken(STRING, "no", "no", 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left: groupingExpr{
							binaryExpr{
								left:     literalExpr{5},
								operator: newTokenNoLiteral(PLUS),
								right:    literalExpr{3},
							},
						},
						operator: newTokenNoLiteral(LESS),
						right:    literalExpr{10},
					},
					operator: newTokenNoLiteral(QUESTION),
					right:    literalExpr{"yes"},
				},
				operator: newTokenNoLiteral(COLON),
				right:    literalExpr{"no"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc:  "simple_or",
			input: []token{newTokenNoLiteral(TRUE), newTokenNoLiteral(OR), newTokenNoLiteral(FALSE)},
			want:  logicalExpr{left: literalExpr{true}, operator: newTokenNoLiteral(OR), right: literalExpr{false}},
		},
		{
			desc: "chained_or",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(OR),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(OR),
				newTokenNoLiteral(NIL),
			},
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteral(OR),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{nil},
			},
		},
		{
			desc: "or_with_expressions",
			input: []token{
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "2", 2, 0),
				newTokenNoLiteral(OR),
				newToken(NUMBER, "3", 3, 0),
				newTokenNoLiteral(STAR),
				newToken(NUMBER, "4", 4, 0),
			},
			want: logicalExpr{
				left: binaryExpr{
					left:     literalExpr{1},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{2},
				},
				operator: newTokenNoLiteral(OR),
				right: binaryExpr{
					left:     literalExpr{3},
					operator: newTokenNoLiteral(STAR),
					right:    literalExpr{4},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
	}{
		{
			desc:  "simple_and",
			input: []token{newTokenNoLiteral(TRUE), newTokenNoLiteral(AND), newTokenNoLiteral(FALSE)},
			want:  logicalExpr{left: literalExpr{true}, operator: newTokenNoLiteral(AND), right: literalExpr{false}},
		},
		{
			desc: "chained_and",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(AND),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(AND),
				newTokenNoLiteral(NIL),
			},
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteral(AND),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{nil},
			},
		},
		{
			desc: "and_with_expressions",
			input: []token{
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "2", 2, 0),
				newTokenNoLiteral(AND),
				newToken(NUMBER, "3", 3, 0),
				newTokenNoLiteral(STAR),
				newToken(NUMBER, "4", 4, 0),
			},
			want: logicalExpr{
				left: binaryExpr{
					left:     literalExpr{1},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{2},
				},
				operator: newTokenNoLiteral(AND),
				right: binaryExpr{
					left:     literalExpr{3},
					operator: newTokenNoLiteral(STAR),
					right:    literalExpr{4},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
		err   error
	}{
		{
			desc: "simple_assignment",
			input: []token{
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "42", 42, 0),
			},
			want: assignExpr{
				name:  newToken(IDENTIFIER, "x", "x", 0),
				value: literalExpr{42},
			},
		},
		{
			desc: "chained_assignment",
			input: []token{
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(EQUAL),
				newToken(IDENTIFIER, "y", "y", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "42", 42, 0),
			},
			want: assignExpr{
				name: newToken(IDENTIFIER, "x", "x", 0),
				value: assignExpr{
					name:  newToken(IDENTIFIER, "y", "y", 0),
					value: literalExpr{42},
				},
			},
		},
		{
			desc: "assignment_with_expression",
			input: []token{
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "5", 5, 0),
			},
			want: assignExpr{
				name: newToken(IDENTIFIER, "x", "x", 0),
				value: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{5},
				},
			},
		},
		{
			desc: "invalid_assignment_target",
			input: []token{
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "10", 10, 0),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(EQUAL), "Invalid assignment target."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
		err   error
	}{
		{
			desc: "simple_call",
			input: []token{
				newToken(IDENTIFIER, "print", "print", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(STRING, "hello", "hello", 0),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "print", "print", 0)},
				paren:  newTokenNoLiteral(RIGHT_PAREN),
				arguments: []expr{
					literalExpr{"hello"},
				},
			},
		},
		{
			desc: "no_arguments",
			input: []token{
				newToken(IDENTIFIER, "clock", "clock", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: callExpr{
				callee:    variableExpr{newToken(IDENTIFIER, "clock", "clock", 0)},
				paren:     newTokenNoLiteral(RIGHT_PAREN),
				arguments: []expr{},
			},
		},
		{
			desc: "multiple_arguments",
			input: []token{
				newToken(IDENTIFIER, "sum", "sum", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(COMMA),
				newToken(NUMBER, "2", 2, 0),
				newTokenNoLiteral(COMMA),
				newToken(NUMBER, "3", 3, 0),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "sum", "sum", 0)},
				paren:  newTokenNoLiteral(RIGHT_PAREN),
				arguments: []expr{
					literalExpr{1},
					literalExpr{2},
					literalExpr{3},
				},
			},
		},
		{
			desc: "nested_calls",
			input: []token{
				newToken(IDENTIFIER, "outer", "outer", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(IDENTIFIER, "inner", "inner", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "outer", "outer", 0)},
				paren:  newTokenNoLiteral(RIGHT_PAREN),
				arguments: []expr{
					callExpr{
						callee: variableExpr{newToken(IDENTIFIER, "inner", "inner", 0)},
						paren:  newTokenNoLiteral(RIGHT_PAREN),
						arguments: []expr{
							literalExpr{42},
						},
					},
				},
			},
		},
		{
			desc: "multiple_consecutive_calls",
			input: []token{
				newToken(IDENTIFIER, "first", "first", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: callExpr{
				callee: callExpr{
					callee: callExpr{
						callee:    variableExpr{newToken(IDENTIFIER, "first", "first", 0)},
						paren:     newTokenNoLiteral(RIGHT_PAREN),
						arguments: []expr{},
					},
					paren:     newTokenNoLiteral(RIGHT_PAREN),
					arguments: []expr{},
				},
				paren:     newTokenNoLiteral(RIGHT_PAREN),
				arguments: []expr{},
			},
		},
		{
			desc: "call_with_expressions",
			input: []token{
				newToken(IDENTIFIER, "calc", "calc", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "2", 2, 0),
				newTokenNoLiteral(COMMA),
				newToken(NUMBER, "3", 3, 0),
				newTokenNoLiteral(STAR),
				newToken(NUMBER, "4", 4, 0),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: callExpr{
				callee: variableExpr{newToken(IDENTIFIER, "calc", "calc", 0)},
				paren:  newTokenNoLiteral(RIGHT_PAREN),
				arguments: []expr{
					binaryExpr{
						left:     literalExpr{1},
						operator: newTokenNoLiteral(PLUS),
						right:    literalExpr{2},
					},
					binaryExpr{
						left:     literalExpr{3},
						operator: newTokenNoLiteral(STAR),
						right:    literalExpr{4},
					},
				},
			},
		},
		{
			desc: "missing_right_paren",
			input: []token{
				newToken(IDENTIFIER, "print", "print", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(STRING, "hello", "hello", 0),
			},
			want: nil,
			err:  NewParseError(newToken(STRING, "hello", "hello", 0), "Expect ')' after arguments."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  expr
		err   error
	}{
		{
			desc: "expr_COMMA_expr_COMMA_expr",
			input: []token{
				newToken(NUMBER, "12", 12, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "9", 9, 0),
				newTokenNoLiteral(COMMA),
				newToken(NUMBER, "78", 78, 0),
				newTokenNoLiteral(GREATER_EQUAL),
				newToken(NUMBER, "6", 6, 0),
				newTokenNoLiteral(COMMA),
				newToken(NUMBER, "13.5", 13.5, 0),
				newTokenNoLiteral(BANG_EQUAL),
				newToken(NUMBER, "51.3", 51.3, 0),
			},
			want: binaryExpr{
				left: binaryExpr{
					left: binaryExpr{
						left:     literalExpr{12},
						operator: newTokenNoLiteral(GREATER),
						right:    literalExpr{9},
					},
					operator: newTokenNoLiteral(COMMA),
					right: binaryExpr{
						left:     literalExpr{78},
						operator: newTokenNoLiteral(GREATER_EQUAL),
						right:    literalExpr{6},
					},
				},
				operator: newTokenNoLiteral(COMMA),
				right: binaryExpr{
					left:     literalExpr{13.5},
					operator: newTokenNoLiteral(BANG_EQUAL),
					right:    literalExpr{51.3},
				},
			},
		},
		{
			desc: "Missing_left_operand_in_binary",
			input: []token{
				newTokenNoLiteral(SLASH),
				newToken(NUMBER, "13.5", 13.5, 0),
				newTokenNoLiteral(BANG_EQUAL),
				newToken(NUMBER, "51.3", 51.3, 0),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(SLASH), "Expect left operand."),
		},
		{
			desc: "or_simple",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(OR),
				newTokenNoLiteral(FALSE),
			},
			want: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{false},
			},
		},
		{
			desc: "and_simple",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(AND),
				newTokenNoLiteral(FALSE),
			},
			want: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{false},
			},
		},
		{
			desc: "chained_or",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(OR),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(OR),
				newTokenNoLiteral(NIL),
			},
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteral(OR),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{nil},
			},
		},
		{
			desc: "chained_and",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(AND),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(AND),
				newTokenNoLiteral(NIL),
			},
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteral(AND),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{nil},
			},
		},
		{
			desc: "mixed_and_or",
			input: []token{
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(AND),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(OR),
				newTokenNoLiteral(NIL),
			},
			want: logicalExpr{
				left: logicalExpr{
					left:     literalExpr{true},
					operator: newTokenNoLiteral(AND),
					right:    literalExpr{false},
				},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{nil},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "printStmt_Simple",
			input: []token{
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: printStmt{expr: literalExpr{42}},
		},
		{
			desc: "printStmt_with_binaryExpr",
			input: []token{
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "8", 8, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: printStmt{
				expr: binaryExpr{
					left:     literalExpr{42},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{8},
				},
			},
		},
		{
			desc: "missing_semicolon",
			input: []token{
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(EOF),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(EOF), "Expect ';' after expression."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "exprStmt",
			input: []token{
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "8", 8, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: exprStmt{
				expr: binaryExpr{
					left:     literalExpr{42},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{8},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "simple_block",
			input: []token{
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{
				statements: []stmt{
					printStmt{expr: literalExpr{42}},
				},
			},
		},
		{
			desc: "empty_block",
			input: []token{
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{statements: []stmt{}},
		},
		{
			desc: "nested_blocks",
			input: []token{
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "2", 2, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
				newTokenNoLiteral(RIGHT_BRACE),
			},
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
			desc: "block_with_declarations",
			input: []token{
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{
				statements: []stmt{
					varStmt{
						name:        newToken(IDENTIFIER, "x", "x", 0),
						initializer: literalExpr{10},
					},
					printStmt{
						expr: variableExpr{newToken(IDENTIFIER, "x", "x", 0)},
					},
				},
			},
		},
		{
			desc: "missing_right_brace",
			input: []token{
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(EOF),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(EOF), "Expect '}' after block."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "simple_if_then",
			input: []token{
				newTokenNoLiteral(IF),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: ifStmt{
				condition: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(GREATER),
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
			desc: "simple_if_then_else",
			input: []token{
				newTokenNoLiteral(IF),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
				newTokenNoLiteral(ELSE),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(FALSE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: ifStmt{
				condition: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(GREATER),
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
			desc: "if_with_parentheses",
			input: []token{
				newTokenNoLiteral(IF),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: ifStmt{
				condition: groupingExpr{binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(GREATER),
					right:    literalExpr{5},
				}},
				thenBranch: blockStmt{statements: []stmt{printStmt{expr: literalExpr{true}}}},
				elseBranch: nil,
			},
		},
		{
			desc: "if_missing_block",
			input: []token{
				newTokenNoLiteral(IF),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(PRINT), "Expect block."),
		},
		{
			desc: "else_missing_block",
			input: []token{
				newTokenNoLiteral(IF),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
				newTokenNoLiteral(ELSE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(PRINT), "Expect block."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "simple_while",
			input: []token{
				newTokenNoLiteral(WHILE),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
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
			desc: "while_complex_condition",
			input: []token{
				newTokenNoLiteral(WHILE),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: whileStmt{
				condition: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(GREATER),
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
			desc: "while_missing_block",
			input: []token{
				newTokenNoLiteral(WHILE),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(PRINT), "Expect block."),
		},
		{
			desc: "while_with_parentheses",
			input: []token{
				newTokenNoLiteral(WHILE),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: whileStmt{
				condition: groupingExpr{binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(GREATER),
					right:    literalExpr{5},
				}},
				body: blockStmt{statements: []stmt{printStmt{expr: literalExpr{true}}}},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "basic_for_loop",
			input: []token{
				newTokenNoLiteral(FOR),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "0", 0, 0),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(LESS),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{
				statements: []stmt{
					exprStmt{expr: assignExpr{
						name:  newToken(IDENTIFIER, "i", "i", 0),
						value: literalExpr{0},
					}},
					whileStmt{
						condition: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "i", "i", 0)},
							operator: newTokenNoLiteral(LESS),
							right:    literalExpr{10},
						},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 0)}},
								exprStmt{expr: assignExpr{
									name: newToken(IDENTIFIER, "i", "i", 0),
									value: binaryExpr{
										left:     variableExpr{newToken(IDENTIFIER, "i", "i", 0)},
										operator: newTokenNoLiteral(PLUS),
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
			desc: "for_loop_without_initializer",
			input: []token{
				newTokenNoLiteral(FOR),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(LESS),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(EQUAL),
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "x", "x", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: whileStmt{
				condition: binaryExpr{
					left:     variableExpr{newToken(IDENTIFIER, "x", "x", 0)},
					operator: newTokenNoLiteral(LESS),
					right:    literalExpr{5},
				},
				body: blockStmt{
					statements: []stmt{
						printStmt{expr: variableExpr{newToken(IDENTIFIER, "x", "x", 0)}},
						exprStmt{expr: assignExpr{
							name: newToken(IDENTIFIER, "x", "x", 0),
							value: binaryExpr{
								left:     variableExpr{newToken(IDENTIFIER, "x", "x", 0)},
								operator: newTokenNoLiteral(PLUS),
								right:    literalExpr{1},
							},
						}},
					},
				},
			},
		},
		{
			desc: "for_loop_without_condition",
			input: []token{
				newTokenNoLiteral(FOR),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "0", 0, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{
				statements: []stmt{
					exprStmt{expr: assignExpr{
						name:  newToken(IDENTIFIER, "i", "i", 0),
						value: literalExpr{0},
					}},
					whileStmt{
						condition: literalExpr{true},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 0)}},
								exprStmt{expr: assignExpr{
									name: newToken(IDENTIFIER, "i", "i", 0),
									value: binaryExpr{
										left:     variableExpr{newToken(IDENTIFIER, "i", "i", 0)},
										operator: newTokenNoLiteral(PLUS),
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
			desc: "for_loop_without_increment",
			input: []token{
				newTokenNoLiteral(FOR),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "0", 0, 0),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(LESS),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{
				statements: []stmt{
					exprStmt{expr: assignExpr{
						name:  newToken(IDENTIFIER, "i", "i", 0),
						value: literalExpr{0},
					}},
					whileStmt{
						condition: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "i", "i", 0)},
							operator: newTokenNoLiteral(LESS),
							right:    literalExpr{10},
						},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 0)}},
							},
						},
					},
				},
			},
		},
		{
			desc: "empty_for_loop",
			input: []token{
				newTokenNoLiteral(FOR),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newTokenNoLiteral(TRUE),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: whileStmt{
				condition: literalExpr{true},
				body: blockStmt{
					statements: []stmt{
						printStmt{expr: literalExpr{true}},
					},
				},
			},
		},
		{
			desc: "for_loop_with_var_declaration",
			input: []token{
				newTokenNoLiteral(FOR),
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "0", 0, 0),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(LESS),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(SEMICOLON),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(EQUAL),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "1", 1, 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "i", "i", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: blockStmt{
				statements: []stmt{
					varStmt{
						name:        newToken(IDENTIFIER, "i", "i", 0),
						initializer: literalExpr{0},
					},
					whileStmt{
						condition: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "i", "i", 0)},
							operator: newTokenNoLiteral(LESS),
							right:    literalExpr{5},
						},
						body: blockStmt{
							statements: []stmt{
								printStmt{expr: variableExpr{newToken(IDENTIFIER, "i", "i", 0)}},
								exprStmt{expr: assignExpr{
									name: newToken(IDENTIFIER, "i", "i", 0),
									value: binaryExpr{
										left:     variableExpr{newToken(IDENTIFIER, "i", "i", 0)},
										operator: newTokenNoLiteral(PLUS),
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
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "var_declaration_no_initializer",
			input: []token{
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "foo", "foo", 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: varStmt{
				name:        newToken(IDENTIFIER, "foo", "foo", 0),
				initializer: nil,
			},
		},
		{
			desc: "var_declaration_with_initializer",
			input: []token{
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "foo", "foo", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: varStmt{
				name:        newToken(IDENTIFIER, "foo", "foo", 0),
				initializer: literalExpr{42},
			},
		},
		{
			desc: "missing_semicolon",
			input: []token{
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "foo", "foo", 0),
				newTokenNoLiteral(EOF),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(EOF), "Expect ';' after variable declaration."),
		},
		{
			desc: "missing_identifier",
			input: []token{
				newTokenNoLiteral(VAR),
				newTokenNoLiteral(SEMICOLON),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(SEMICOLON), "Expect variable name."),
		},
		{
			desc: "expr_statement_fallback",
			input: []token{
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: exprStmt{expr: literalExpr{42}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "function_declaration",
			input: []token{
				newTokenNoLiteral(FN),
				newToken(IDENTIFIER, "greet", "greet", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(IDENTIFIER, "name", "name", 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(STRING, "Hello, ", "Hello, ", 0),
				newTokenNoLiteral(PLUS),
				newToken(IDENTIFIER, "name", "name", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: functionStmt{
				name: newToken(IDENTIFIER, "greet", "greet", 0),
				params: []token{
					newToken(IDENTIFIER, "name", "name", 0),
				},
				body: []stmt{
					printStmt{
						expr: binaryExpr{
							left:     literalExpr{"Hello, "},
							operator: newTokenNoLiteral(PLUS),
							right:    variableExpr{newToken(IDENTIFIER, "name", "name", 0)},
						},
					},
				},
			},
		},
		{
			desc: "function_no_params",
			input: []token{
				newTokenNoLiteral(FN),
				newToken(IDENTIFIER, "hello", "hello", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(STRING, "Hello!", "Hello!", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: functionStmt{
				name:   newToken(IDENTIFIER, "hello", "hello", 0),
				params: []token{},
				body: []stmt{
					printStmt{expr: literalExpr{"Hello!"}},
				},
			},
		},
		{
			desc: "function_multiple_params",
			input: []token{
				newTokenNoLiteral(FN),
				newToken(IDENTIFIER, "add", "add", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(IDENTIFIER, "a", "a", 0),
				newTokenNoLiteral(COMMA),
				newToken(IDENTIFIER, "b", "b", 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "a", "a", 0),
				newTokenNoLiteral(PLUS),
				newToken(IDENTIFIER, "b", "b", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: functionStmt{
				name: newToken(IDENTIFIER, "add", "add", 0),
				params: []token{
					newToken(IDENTIFIER, "a", "a", 0),
					newToken(IDENTIFIER, "b", "b", 0),
				},
				body: []stmt{
					printStmt{
						expr: binaryExpr{
							left:     variableExpr{newToken(IDENTIFIER, "a", "a", 0)},
							operator: newTokenNoLiteral(PLUS),
							right:    variableExpr{newToken(IDENTIFIER, "b", "b", 0)},
						},
					},
				},
			},
		},
		{
			desc: "missing_function_name",
			input: []token{
				newTokenNoLiteral(FN),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(LEFT_PAREN), "Expect function name."),
		},
		{
			desc: "missing_left_paren",
			input: []token{
				newTokenNoLiteral(FN),
				newToken(IDENTIFIER, "test", "test", 0),
				newToken(IDENTIFIER, "param", "param", 0),
				newTokenNoLiteral(RIGHT_PAREN),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: nil,
			err:  NewParseError(newToken(IDENTIFIER, "param", "param", 0), "Expect '(' after function name."),
		},
		{
			desc: "missing_right_paren",
			input: []token{
				newTokenNoLiteral(FN),
				newToken(IDENTIFIER, "test", "test", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newToken(IDENTIFIER, "param", "param", 0),
				newTokenNoLiteral(LEFT_BRACE),
				newTokenNoLiteral(RIGHT_BRACE),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(LEFT_BRACE), "Expect ')' after parameters."),
		},
		{
			desc: "missing_body",
			input: []token{
				newTokenNoLiteral(FN),
				newToken(IDENTIFIER, "test", "test", 0),
				newTokenNoLiteral(LEFT_PAREN),
				newTokenNoLiteral(RIGHT_PAREN),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(RIGHT_PAREN), "Expect block."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  stmt
		err   error
	}{
		{
			desc: "return_with_value",
			input: []token{
				newTokenNoLiteral(RETURN),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: returnStmt{
				keyword: newTokenNoLiteral(RETURN),
				value:   literalExpr{42},
			},
		},
		{
			desc: "return_without_value",
			input: []token{
				newTokenNoLiteral(RETURN),
				newTokenNoLiteral(SEMICOLON),
			},
			want: returnStmt{
				keyword: newTokenNoLiteral(RETURN),
				value:   nil,
			},
		},
		{
			desc: "return_with_expression",
			input: []token{
				newTokenNoLiteral(RETURN),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(PLUS),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(SEMICOLON),
			},
			want: returnStmt{
				keyword: newTokenNoLiteral(RETURN),
				value: binaryExpr{
					left:     literalExpr{10},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{5},
				},
			},
		},
		{
			desc: "missing_semicolon",
			input: []token{
				newTokenNoLiteral(RETURN),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(EOF),
			},
			want: nil,
			err:  NewParseError(newTokenNoLiteral(EOF), "Expect ';' after return value."),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
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
		input []token
		want  []stmt
		err   error
	}{
		{
			desc: "single_print_statement",
			input: []token{
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(EOF),
			},
			want: []stmt{printStmt{expr: literalExpr{42}}},
		},
		{
			desc: "multiple_statements",
			input: []token{
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "foo", "foo", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(PRINT),
				newToken(IDENTIFIER, "foo", "foo", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(EOF),
			},
			want: []stmt{
				varStmt{
					name:        newToken(IDENTIFIER, "foo", "foo", 0),
					initializer: literalExpr{42},
				},
				printStmt{expr: variableExpr{newToken(IDENTIFIER, "foo", "foo", 0)}},
			},
		},
		{
			desc: "variable_with_ternary",
			input: []token{
				newTokenNoLiteral(VAR),
				newToken(IDENTIFIER, "result", "result", 0),
				newTokenNoLiteral(EQUAL),
				newToken(NUMBER, "10", 10, 0),
				newTokenNoLiteral(GREATER),
				newToken(NUMBER, "5", 5, 0),
				newTokenNoLiteral(QUESTION),
				newToken(STRING, "yes", "yes", 0),
				newTokenNoLiteral(COLON),
				newToken(STRING, "no", "no", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(EOF),
			},
			want: []stmt{
				varStmt{
					name: newToken(IDENTIFIER, "result", "result", 0),
					initializer: binaryExpr{
						left: binaryExpr{
							left: binaryExpr{
								left:     literalExpr{10},
								operator: newTokenNoLiteral(GREATER),
								right:    literalExpr{5},
							},
							operator: newTokenNoLiteral(QUESTION),
							right:    literalExpr{"yes"},
						},
						operator: newTokenNoLiteral(COLON),
						right:    literalExpr{"no"},
					},
				},
			},
		},
		{
			desc: "empty_input",
			input: []token{
				newTokenNoLiteral(EOF),
			},
			want: []stmt{},
		},
		{
			desc: "parse_error",
			input: []token{
				newTokenNoLiteral(PRINT),
				newToken(NUMBER, "42", 42, 0),
				newTokenNoLiteral(EOF),
			},
			want: []stmt{},
		},
		{
			desc: "synchronize_recovers_at_statement_boundary",
			input: []token{
				newTokenNoLiteral(VAR),
				newToken(NUMBER, "42", 42, 0), // Missing semicolon
				newTokenNoLiteral(PRINT),      // Next statement boundary
				newToken(STRING, "hello", "hello", 0),
				newTokenNoLiteral(SEMICOLON),
				newTokenNoLiteral(EOF),
			},
			want: []stmt{
				printStmt{expr: literalExpr{"hello"}},
			},
			err: nil, // Synchronize should allow parsing to continue after error
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
			got, err := parser.Parse()
			if err != nil {
				assert.Equal(t, tC.err, err)
				return
			}
			assert.Equal(t, tC.want, got)
		})
	}
}
