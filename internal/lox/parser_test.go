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
			err:  NewParseError(newTokenNoLiteral(SLASH), "expect left operand"),
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

func Test_statement(t *testing.T) {
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
