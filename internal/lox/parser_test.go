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

func Test_expression(t *testing.T) {
	testCases := []struct {
		desc  string
		input []token
		want  expr
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
			got, err := parser.expression()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}

func TestParse_missing_left_operand(t *testing.T) {
	testCases := []struct {
		desc  string
		input []token
		want  expr
		err   error
	}{
		{
			desc: "Missing_left_operand_in_binary",
			input: []token{
				newTokenNoLiteral(SLASH),
				newToken(NUMBER, "13.5", 13.5, 0),
				newTokenNoLiteral(BANG_EQUAL),
				newToken(NUMBER, "51.3", 51.3, 0),
			},
			want: nil,
			err:  NewLoxError(0, "'SLASH'", "expected left operand"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := NewParser(tC.input)
			got, err := parser.Parse()
			if err != nil {
				assert.Equal(t, tC.err, err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}
