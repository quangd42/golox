package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_interpretLiteralExpr(t *testing.T) {
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

func Test_interpretUnaryExpr(t *testing.T) {
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

func Test_interpretBinaryExpr(t *testing.T) {
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

func Test_interpretVariableExpr(t *testing.T) {
	testCases := []struct {
		desc    string
		input   variableExpr
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc:    "variable_exists",
			input:   variableExpr{name: newToken(IDENTIFIER, "x", nil, 1)},
			initEnv: map[string]any{"x": 42.0},
			want:    42.0,
			err:     nil,
		},
		{
			desc:    "variable_undefined",
			input:   variableExpr{name: newToken(IDENTIFIER, "y", nil, 1)},
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(IDENTIFIER, "y", nil, 1), "Undefined variable 'y'."),
		},
		{
			desc:    "variable_nil",
			input:   variableExpr{name: newToken(IDENTIFIER, "z", nil, 1)},
			initEnv: map[string]any{"z": nil},
			want:    nil,
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			got, err := interpreter.visitVariableExpr(tC.input)
			assert.Equal(t, tC.want, got)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}
		})
	}
}

func Test_interpretAssignExpr(t *testing.T) {
	testCases := []struct {
		desc    string
		input   assignExpr
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc: "valid_assignment",
			input: assignExpr{
				name:  newToken(IDENTIFIER, "x", nil, 1),
				value: literalExpr{100.0},
			},
			initEnv: map[string]any{"x": 42.0},
			want:    100.0,
			err:     nil,
		},
		{
			desc: "undefined_variable",
			input: assignExpr{
				name:  newToken(IDENTIFIER, "y", nil, 1),
				value: literalExpr{200.0},
			},
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(IDENTIFIER, "y", nil, 1), "Undefined variable 'y'."),
		},
		{
			desc: "assign_string",
			input: assignExpr{
				name:  newToken(IDENTIFIER, "z", nil, 1),
				value: literalExpr{"hello"},
			},
			initEnv: map[string]any{"z": "world"},
			want:    "hello",
			err:     nil,
		},
		{
			desc: "assign_string_to_int",
			input: assignExpr{
				name:  newToken(IDENTIFIER, "q", nil, 1),
				value: literalExpr{"string"},
			},
			initEnv: map[string]any{"q": 42},
			want:    "string",
			err:     nil,
		},
		{
			desc: "assign_nil",
			input: assignExpr{
				name:  newToken(IDENTIFIER, "w", nil, 1),
				value: literalExpr{nil},
			},
			initEnv: map[string]any{"w": 42.0},
			want:    nil,
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			got, err := interpreter.visitAssignExpr(tC.input)
			assert.Equal(t, tC.want, got)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}
		})
	}
}

func Test_interpretLogicalExpr(t *testing.T) {
	testCases := []struct {
		desc    string
		input   logicalExpr
		want    any
		wantErr error
	}{
		{
			desc: "OR_leftTrue",
			input: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{false},
			},
			want:    true,
			wantErr: nil,
		},
		{
			desc: "OR_leftFalse",
			input: logicalExpr{
				left:     literalExpr{false},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{true},
			},
			want:    true,
			wantErr: nil,
		},
		{
			desc: "OR_bothFalse",
			input: logicalExpr{
				left:     literalExpr{false},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{false},
			},
			want:    false,
			wantErr: nil,
		},
		{
			desc: "AND_bothTrue",
			input: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{true},
			},
			want:    true,
			wantErr: nil,
		},
		{
			desc: "AND_leftFalse",
			input: logicalExpr{
				left:     literalExpr{false},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{true},
			},
			want:    false,
			wantErr: nil,
		},
		{
			desc: "AND_rightFalse",
			input: logicalExpr{
				left:     literalExpr{true},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{false},
			},
			want:    false,
			wantErr: nil,
		},
		{
			desc: "OR_nonBoolean",
			input: logicalExpr{
				left:     literalExpr{"string"},
				operator: newTokenNoLiteral(OR),
				right:    literalExpr{true},
			},
			want:    "string",
			wantErr: nil,
		},
		{
			desc: "AND_nonBoolean",
			input: logicalExpr{
				left:     literalExpr{nil},
				operator: newTokenNoLiteral(AND),
				right:    literalExpr{true},
			},
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			got, err := interpreter.visitLogicalExpr(tC.input)
			assert.Equal(t, tC.want, got)
			if tC.wantErr != nil {
				assert.EqualError(t, err, tC.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_interpretCallExpr(t *testing.T) {
	testCases := []struct {
		desc    string
		input   callExpr
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc: "call_simple_function",
			input: callExpr{
				callee: variableExpr{name: newToken(IDENTIFIER, "test", nil, 1)},
				paren:  newToken(RIGHT_PAREN, ")", nil, 1),
				arguments: []expr{
					literalExpr{42.0},
				},
			},
			initEnv: map[string]any{
				"test": function{
					declaration: functionStmt{
						name:   newToken(IDENTIFIER, "test", nil, 1),
						params: []token{newToken(IDENTIFIER, "x", nil, 1)},
						body: []stmt{
							returnStmt{
								keyword: newToken(RETURN, "return", nil, 1),
								value:   variableExpr{name: newToken(IDENTIFIER, "x", nil, 1)},
							},
						},
					},
				},
			},
			want: 42.0,
			err:  nil,
		},
		{
			desc: "call_undefined_function",
			input: callExpr{
				callee: variableExpr{name: newToken(IDENTIFIER, "undefined", nil, 1)},
				paren:  newToken(RIGHT_PAREN, ")", nil, 1),
				arguments: []expr{
					literalExpr{42.0},
				},
			},
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(IDENTIFIER, "undefined", nil, 1), "Undefined variable 'undefined'."),
		},
		{
			desc: "call_non_function",
			input: callExpr{
				callee: variableExpr{name: newToken(IDENTIFIER, "notfunc", nil, 1)},
				paren:  newToken(RIGHT_PAREN, ")", nil, 1),
				arguments: []expr{
					literalExpr{42.0},
				},
			},
			initEnv: map[string]any{
				"notfunc": "string",
			},
			want: nil,
			err:  NewRuntimeError(newToken(RIGHT_PAREN, ")", nil, 1), "Can only call functions and classes."),
		},
		{
			desc: "wrong_arity",
			input: callExpr{
				callee: variableExpr{name: newToken(IDENTIFIER, "test", nil, 1)},
				paren:  newToken(RIGHT_PAREN, ")", nil, 1),
				arguments: []expr{
					literalExpr{42.0},
					literalExpr{43.0},
				},
			},
			initEnv: map[string]any{
				"test": function{
					declaration: functionStmt{
						name:   newToken(IDENTIFIER, "test", nil, 1),
						params: []token{newToken(IDENTIFIER, "x", nil, 1)},
						body: []stmt{
							returnStmt{
								keyword: newToken(RETURN, "return", nil, 1),
								value:   variableExpr{name: newToken(IDENTIFIER, "x", nil, 1)},
							},
						},
					},
				},
			},
			want: nil,
			err:  NewRuntimeError(newToken(RIGHT_PAREN, ")", nil, 1), "Expected 1 arguments but got 2."),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			got, err := interpreter.visitCallExpr(tC.input)
			assert.Equal(t, tC.want, got)
			if tC.err != nil {
				assert.EqualError(t, err, tC.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_interpretVarStmt(t *testing.T) {
	testCases := []struct {
		desc        string
		input       varStmt
		wantEnvVal  any
		wantEnvName string
		err         error
	}{
		{
			desc: "without_initializer",
			input: varStmt{
				name:        newToken(IDENTIFIER, "x", nil, 1),
				initializer: nil,
			},
			wantEnvVal:  nil,
			wantEnvName: "x",
			err:         nil,
		},
		{
			desc: "with_initializer",
			input: varStmt{
				name:        newToken(IDENTIFIER, "y", nil, 1),
				initializer: literalExpr{42.0},
			},
			wantEnvVal:  42.0,
			wantEnvName: "y",
			err:         nil,
		},
		{
			desc: "with_string_initializer",
			input: varStmt{
				name:        newToken(IDENTIFIER, "z", nil, 1),
				initializer: literalExpr{"hello"},
			},
			wantEnvVal:  "hello",
			wantEnvName: "z",
			err:         nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			err := interpreter.visitVarStmt(tC.input)

			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			} else {
				val, ok := interpreter.env.values[tC.wantEnvName]
				assert.True(t, ok)
				assert.Equal(t, tC.wantEnvVal, val)
			}
		})
	}
}

func Test_interpretBlockStmt(t *testing.T) {
	testCases := []struct {
		desc    string
		stmts   []stmt
		initEnv map[string]any
		wantEnv map[string]any
		err     error
	}{
		{
			desc: "access_and_modify_global_var",
			stmts: []stmt{
				varStmt{
					name:        newToken(IDENTIFIER, "local", nil, 1),
					initializer: literalExpr{42.0},
				},
				exprStmt{
					expr: assignExpr{
						name:  newToken(IDENTIFIER, "global", nil, 1),
						value: literalExpr{100.0},
					},
				},
			},
			initEnv: map[string]any{"global": 50.0},
			wantEnv: map[string]any{"global": 100.0},
			err:     nil,
		},
		{
			desc: "nested_blocks_access_global",
			stmts: []stmt{
				blockStmt{
					statements: []stmt{
						varStmt{
							name:        newToken(IDENTIFIER, "a", nil, 1),
							initializer: literalExpr{1.0},
						},
						blockStmt{
							statements: []stmt{
								exprStmt{
									expr: assignExpr{
										name:  newToken(IDENTIFIER, "global", nil, 1),
										value: literalExpr{200.0},
									},
								},
							},
						},
					},
				},
			},
			initEnv: map[string]any{"global": 100.0},
			wantEnv: map[string]any{"global": 200.0},
			err:     nil,
		},
		{
			desc: "local_var_not_accessible_after_block",
			stmts: []stmt{
				varStmt{
					name:        newToken(IDENTIFIER, "local", nil, 1),
					initializer: literalExpr{42.0},
				},
			},
			initEnv: map[string]any{},
			wantEnv: map[string]any{},
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			block := blockStmt{statements: tC.stmts}
			err := interpreter.visitBlockStmt(block)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}

			// Check global environment matches expected
			for k, v := range tC.wantEnv {
				val, exists := interpreter.env.values[k]
				assert.True(t, exists)
				assert.Equal(t, v, val)
			}

			// Check local variables are not accessible
			for _, stmt := range tC.stmts {
				if varStmt, ok := stmt.(varStmt); ok {
					_, exists := interpreter.env.values[varStmt.name.lexeme]
					assert.False(t, exists)
				}
			}
		})
	}
}

func Test_interpretReturnStmt(t *testing.T) {
	testCases := []struct {
		desc    string
		input   returnStmt
		initEnv map[string]any
		wantVal any
		wantErr error
	}{
		{
			desc: "return_literal_number",
			input: returnStmt{
				keyword: newToken(RETURN, "return", nil, 1),
				value:   literalExpr{42.0},
			},
			initEnv: map[string]any{},
			wantVal: 42.0,
			wantErr: nil,
		},
		{
			desc: "return_binary_expr",
			input: returnStmt{
				keyword: newToken(RETURN, "return", nil, 1),
				value: binaryExpr{
					left:     literalExpr{10.0},
					operator: newTokenNoLiteral(PLUS),
					right:    literalExpr{5.0},
				},
			},
			initEnv: map[string]any{},
			wantVal: 15.0,
			wantErr: nil,
		},
		{
			desc: "return_variable",
			input: returnStmt{
				keyword: newToken(RETURN, "return", nil, 1),
				value:   variableExpr{name: newToken(IDENTIFIER, "x", nil, 1)},
			},
			initEnv: map[string]any{"x": "hello"},
			wantVal: "hello",
			wantErr: nil,
		},
		{
			desc: "return_nil",
			input: returnStmt{
				keyword: newToken(RETURN, "return", nil, 1),
				value:   literalExpr{nil},
			},
			initEnv: map[string]any{},
			wantVal: nil,
			wantErr: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err := interpreter.visitReturnStmt(tC.input)
			if tC.wantErr != nil {
				assert.EqualError(t, err, tC.wantErr.Error())
			} else {
				retErr, ok := err.(*returnValue)
				assert.True(t, ok)
				assert.Equal(t, tC.wantVal, retErr.value)
			}
		})
	}
}

func Test_interpretWhileStmt(t *testing.T) {
	testCases := []struct {
		desc      string
		condition expr
		body      stmt
		initEnv   map[string]any
		wantEnv   map[string]any
		err       error
	}{
		{
			desc:      "false_condition_no_iteration",
			condition: literalExpr{false},
			body: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{1.0},
				},
			},
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 0.0},
			err:     nil,
		},
		{
			desc: "counter_loop",
			condition: binaryExpr{
				left:     variableExpr{name: newToken(IDENTIFIER, "counter", nil, 1)},
				operator: newTokenNoLiteral(LESS),
				right:    literalExpr{3.0},
			},
			body: blockStmt{
				statements: []stmt{
					exprStmt{
						expr: assignExpr{
							name: newToken(IDENTIFIER, "counter", nil, 1),
							value: binaryExpr{
								left:     variableExpr{name: newToken(IDENTIFIER, "counter", nil, 1)},
								operator: newTokenNoLiteral(PLUS),
								right:    literalExpr{1.0},
							},
						},
					},
					exprStmt{
						expr: assignExpr{
							name: newToken(IDENTIFIER, "sum", nil, 1),
							value: binaryExpr{
								left:     variableExpr{name: newToken(IDENTIFIER, "sum", nil, 1)},
								operator: newTokenNoLiteral(PLUS),
								right:    literalExpr{1.0},
							},
						},
					},
				},
			},
			initEnv: map[string]any{"counter": 0.0, "sum": 0.0},
			wantEnv: map[string]any{"counter": 3.0, "sum": 3.0},
			err:     nil,
		},
		{
			desc:      "non_boolean_condition",
			condition: variableExpr{name: newToken(IDENTIFIER, "cond", nil, 1)},
			body: blockStmt{statements: []stmt{
				exprStmt{expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{1.0},
				}},
				exprStmt{expr: assignExpr{
					name:  newToken(IDENTIFIER, "cond", nil, 1),
					value: literalExpr{false},
				}},
			}},
			initEnv: map[string]any{"cond": "not a boolean", "x": 0.0},
			wantEnv: map[string]any{"cond": false, "x": 1.0},
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			stmt := whileStmt{
				condition: tC.condition,
				body:      tC.body,
			}

			err := interpreter.visitWhileStmt(stmt)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}

			for k, v := range tC.wantEnv {
				val, exists := interpreter.env.values[k]
				assert.True(t, exists)
				assert.Equal(t, v, val)
			}
		})
	}
}

func Test_interpretIfStmt(t *testing.T) {
	testCases := []struct {
		desc      string
		condition expr
		thenStmt  stmt
		elseStmt  stmt
		initEnv   map[string]any
		wantEnv   map[string]any
		err       error
	}{
		{
			desc:      "true_condition_no_else",
			condition: literalExpr{true},
			thenStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{1.0},
				},
			},
			elseStmt: nil,
			initEnv:  map[string]any{"x": 0.0},
			wantEnv:  map[string]any{"x": 1.0},
			err:      nil,
		},
		{
			desc:      "false_condition_no_else",
			condition: literalExpr{false},
			thenStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{1.0},
				},
			},
			elseStmt: nil,
			initEnv:  map[string]any{"x": 0.0},
			wantEnv:  map[string]any{"x": 0.0},
			err:      nil,
		},
		{
			desc:      "true_condition_with_else",
			condition: literalExpr{true},
			thenStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{1.0},
				},
			},
			elseStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{2.0},
				},
			},
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 1.0},
			err:     nil,
		},
		{
			desc:      "false_condition_with_else",
			condition: literalExpr{false},
			thenStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{1.0},
				},
			},
			elseStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{2},
				},
			},
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 2},
			err:     nil,
		},
		{
			desc:      "non_boolean_condition",
			condition: literalExpr{"not a boolean"},
			thenStmt: exprStmt{
				expr: assignExpr{
					name:  newToken(IDENTIFIER, "x", nil, 1),
					value: literalExpr{"hit then"},
				},
			},
			elseStmt: nil,
			initEnv:  map[string]any{"x": 0.0},
			wantEnv:  map[string]any{"x": "hit then"},
			err:      nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			stmt := ifStmt{
				condition:  tC.condition,
				thenBranch: tC.thenStmt,
				elseBranch: tC.elseStmt,
			}

			err := interpreter.visitIfStmt(stmt)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}

			for k, v := range tC.wantEnv {
				val, exists := interpreter.env.values[k]
				assert.True(t, exists)
				assert.Equal(t, v, val)
			}
		})
	}
}

func Test_interpretFunctionStmt(t *testing.T) {
	testCases := []struct {
		desc    string
		input   functionStmt
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc: "basic_function_declaration",
			input: functionStmt{
				name: newToken(IDENTIFIER, "add", nil, 1),
				params: []token{
					newToken(IDENTIFIER, "a", nil, 1),
					newToken(IDENTIFIER, "b", nil, 1),
				},
				body: []stmt{
					returnStmt{
						keyword: newToken(RETURN, "return", nil, 1),
						value: binaryExpr{
							left:     variableExpr{name: newToken(IDENTIFIER, "a", nil, 1)},
							operator: newTokenNoLiteral(PLUS),
							right:    variableExpr{name: newToken(IDENTIFIER, "b", nil, 1)},
						},
					},
				},
			},
			initEnv: map[string]any{},
			want:    nil,
			err:     nil,
		},
		{
			desc: "empty_function_declaration",
			input: functionStmt{
				name:   newToken(IDENTIFIER, "empty", nil, 1),
				params: []token{},
				body:   []stmt{},
			},
			initEnv: map[string]any{},
			want:    nil,
			err:     nil,
		},
		{
			desc: "redefined_function",
			input: functionStmt{
				name:   newToken(IDENTIFIER, "existing", nil, 1),
				params: []token{},
				body:   []stmt{},
			},
			initEnv: map[string]any{
				"existing": "some_value",
			},
			want: nil,
			err:  nil,
		},
		{
			desc: "fibonacci_function",
			input: functionStmt{
				name: newToken(IDENTIFIER, "fib", nil, 1),
				params: []token{
					newToken(IDENTIFIER, "n", nil, 1),
				},
				body: []stmt{
					ifStmt{
						condition: binaryExpr{
							left:     variableExpr{name: newToken(IDENTIFIER, "n", nil, 1)},
							operator: newTokenNoLiteral(LESS_EQUAL),
							right:    literalExpr{1.0},
						},
						thenBranch: returnStmt{
							keyword: newToken(RETURN, "return", nil, 1),
							value:   variableExpr{name: newToken(IDENTIFIER, "n", nil, 1)},
						},
					},
					returnStmt{
						keyword: newToken(RETURN, "return", nil, 1),
						value: binaryExpr{
							left: callExpr{
								callee: variableExpr{name: newToken(IDENTIFIER, "fib", nil, 1)},
								paren:  newToken(RIGHT_PAREN, ")", nil, 1),
								arguments: []expr{
									binaryExpr{
										left:     variableExpr{name: newToken(IDENTIFIER, "n", nil, 1)},
										operator: newTokenNoLiteral(MINUS),
										right:    literalExpr{2.0},
									},
								},
							},
							operator: newTokenNoLiteral(PLUS),
							right: callExpr{
								callee: variableExpr{name: newToken(IDENTIFIER, "fib", nil, 1)},
								paren:  newToken(RIGHT_PAREN, ")", nil, 1),
								arguments: []expr{
									binaryExpr{
										left:     variableExpr{name: newToken(IDENTIFIER, "n", nil, 1)},
										operator: newTokenNoLiteral(MINUS),
										right:    literalExpr{1.0},
									},
								},
							},
						},
					},
				},
			},
			initEnv: map[string]any{},
			want:    nil,
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			interpreter := NewInterpreter()
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err := interpreter.visitFunctionStmt(tC.input)
			if tC.err != nil {
				assert.EqualError(t, err, tC.err.Error())
			} else {
				assert.NoError(t, err)
				// Verify function was defined in environment
				val, exists := interpreter.env.values[tC.input.name.lexeme]
				assert.True(t, exists)
				_, ok := val.(function)
				assert.True(t, ok)
			}
		})
	}
}
