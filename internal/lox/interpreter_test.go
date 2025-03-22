package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_interpretLiteralExpr(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  any
		err   error
	}{
		{
			desc:  "STRING",
			input: `"a string"`,
			want:  "a string",
			err:   nil,
		},
		{
			desc:  "NUMBER_float64",
			input: "158.2",
			want:  158.2,
			err:   nil,
		},
		{
			desc:  "NUMBER_int",
			input: "2389",
			want:  2389,
			err:   nil,
		},
		{
			desc:  "NIL",
			input: "nil",
			want:  nil,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}
			interpreter := NewInterpreter(nil)
			got, err := interpreter.visitLiteralExpr(expr.(literalExpr))
			assert.Equal(t, tC.want, got)
			assert.Equal(t, tC.err, err)
		})
	}
}

func Test_interpretUnaryExpr(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  any
		err   error
	}{
		{
			desc:  "MINUS__NUMBER__Float",
			input: "-189.228",
			want:  -189.228,
			err:   nil,
		},
		{
			desc:  "MINUS__NUMBER__Int",
			input: "-189",
			want:  float64(-189),
			err:   nil,
		},
		{
			desc:  "MINUS__NUMBER__NaN",
			input: `-"NaN"`,
			want:  nil,
			err:   NewRuntimeError(newToken(MINUS, "-", nil, 1, 0), "Operand must be a number."),
		},
		{
			desc:  "BANG__TRUE",
			input: "!true",
			want:  false,
			err:   nil,
		},
		{
			desc:  "BANG__FALSE",
			input: "!false",
			want:  true,
			err:   nil,
		},
		{
			desc:  "BANG__NIL",
			input: "!nil",
			want:  true,
			err:   nil,
		},
		{
			desc:  "BANG__LITERAL",
			input: `!"some string"`,
			want:  false,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}
			interpreter := NewInterpreter(nil)
			got, err := interpreter.visitUnaryExpr(expr.(unaryExpr))
			assert.Equal(t, tC.want, got)
			if err != nil {
				assert.EqualError(t, err, tC.err.Error())
			}
		})
	}
}

func Test_interpretBinaryExpr(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  any
		err   error
	}{
		{
			desc:  "PLUS_float_float",
			input: "5.0 + 3.0",
			want:  8.0,
			err:   nil,
		},
		{
			desc:  "PLUS_int_int",
			input: "5 + 3",
			want:  8.0,
			err:   nil,
		},
		{
			desc:  "PLUS_float_int",
			input: "5.0 + 3",
			want:  8.0,
			err:   nil,
		},
		{
			desc:  "PLUS_strings",
			input: `"hello" + " world"`,
			want:  "hello world",
			err:   nil,
		},
		{
			desc:  "PLUS_invalid",
			input: "true + 5.0",
			want:  nil,
			err:   NewRuntimeError(newToken(PLUS, "+", nil, 1, 5), "Operands must be either numbers or strings."),
		},
		{
			desc:  "MINUS",
			input: "5.0 - 3.0",
			want:  2.0,
			err:   nil,
		},
		{
			desc:  "MINUS_invalid",
			input: `"string" - 5.0`,
			want:  nil,
			err:   NewRuntimeError(newToken(MINUS, "-", nil, 1, 9), "Operands must be numbers."),
		},
		{
			desc:  "MULTIPLY",
			input: "5.0 * 3.0",
			want:  15.0,
			err:   nil,
		},
		{
			desc:  "MULTIPLY_invalid",
			input: "true * 5.0",
			want:  nil,
			err:   NewRuntimeError(newToken(STAR, "*", nil, 1, 5), "Operands must be numbers."),
		},
		{
			desc:  "DIVIDE",
			input: "15.0 / 3.0",
			want:  5.0,
			err:   nil,
		},
		{
			desc:  "DIVIDE_invalid",
			input: `"string" / 5.0`,
			want:  nil,
			err:   NewRuntimeError(newToken(SLASH, "/", nil, 1, 9), "Operands must be numbers."),
		},
		{
			desc:  "GREATER",
			input: "5.0 > 3.0",
			want:  true,
			err:   nil,
		},
		{
			desc:  "GREATER_invalid",
			input: "true > 5.0",
			want:  nil,
			err:   NewRuntimeError(newToken(GREATER, ">", nil, 1, 5), "Operands must be numbers."),
		},
		{
			desc:  "GREATER_EQUAL",
			input: "5.0 >= 5.0",
			want:  true,
			err:   nil,
		},
		{
			desc:  "GREATER_EQUAL_invalid",
			input: `"string" >= 5.0`,
			want:  nil,
			err:   NewRuntimeError(newToken(GREATER_EQUAL, ">=", nil, 1, 9), "Operands must be numbers."),
		},
		{
			desc:  "LESS",
			input: "3.0 < 5.0",
			want:  true,
			err:   nil,
		},
		{
			desc:  "LESS_invalid",
			input: "true < 5.0",
			want:  nil,
			err:   NewRuntimeError(newToken(LESS, "<", nil, 1, 5), "Operands must be numbers."),
		},
		{
			desc:  "LESS_EQUAL",
			input: "5.0 <= 5.0",
			want:  true,
			err:   nil,
		},
		{
			desc:  "LESS_EQUAL_invalid",
			input: `"string" <= 5.0`,
			want:  nil,
			err:   NewRuntimeError(newToken(LESS_EQUAL, "<=", nil, 1, 9), "Operands must be numbers."),
		},
		{
			desc:  "EQUAL_EQUAL",
			input: "5.0 == 5.0",
			want:  true,
			err:   nil,
		},
		{
			desc:  "BANG_EQUAL",
			input: "5.0 != 3.0",
			want:  true,
			err:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}
			interpreter := NewInterpreter(nil)
			got, err := interpreter.visitBinaryExpr(expr.(binaryExpr))
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
		input   string
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc:    "variable_exists",
			input:   "x",
			initEnv: map[string]any{"x": 42.0},
			want:    42.0,
			err:     nil,
		},
		{
			desc:    "variable_undefined",
			input:   "y",
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(IDENTIFIER, "y", nil, 1, 0), "Undefined variable 'y'."),
		},
		{
			desc:    "variable_nil",
			input:   "z",
			initEnv: map[string]any{"z": nil},
			want:    nil,
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			got, err := interpreter.visitVariableExpr(expr.(variableExpr))
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
		input   string
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc:    "valid_assignment",
			input:   "x = 100",
			initEnv: map[string]any{"x": 42.0},
			want:    100,
			err:     nil,
		},
		{
			desc:    "undefined_variable",
			input:   "y = 200",
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(IDENTIFIER, "y", nil, 1, 0), "Undefined variable 'y'."),
		},
		{
			desc:    "assign_string",
			input:   `z = "hello"`,
			initEnv: map[string]any{"z": "world"},
			want:    "hello",
			err:     nil,
		},
		{
			desc:    "assign_string_to_int",
			input:   `q = "string"`,
			initEnv: map[string]any{"q": 42},
			want:    "string",
			err:     nil,
		},
		{
			desc:    "assign_nil",
			input:   "w = nil",
			initEnv: map[string]any{"w": 42.0},
			want:    nil,
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			got, err := interpreter.visitAssignExpr(expr.(assignExpr))
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
		input   string
		want    any
		wantErr error
	}{
		{
			desc:    "OR_leftTrue",
			input:   "true or false",
			want:    true,
			wantErr: nil,
		},
		{
			desc:    "OR_leftFalse",
			input:   "false or true",
			want:    true,
			wantErr: nil,
		},
		{
			desc:    "OR_bothFalse",
			input:   "false or false",
			want:    false,
			wantErr: nil,
		},
		{
			desc:    "AND_bothTrue",
			input:   "true and true",
			want:    true,
			wantErr: nil,
		},
		{
			desc:    "AND_leftFalse",
			input:   "false and true",
			want:    false,
			wantErr: nil,
		},
		{
			desc:    "AND_rightFalse",
			input:   "true and false",
			want:    false,
			wantErr: nil,
		},
		{
			desc:    "OR_nonBoolean",
			input:   `"string" or true`,
			want:    "string",
			wantErr: nil,
		},
		{
			desc:    "AND_nonBoolean",
			input:   "nil and true",
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			got, err := interpreter.visitLogicalExpr(expr.(logicalExpr))
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
		input   string
		code    string
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc:  "call_simple_function",
			input: "test(42)",
			code: `fn test(x) {
				return x;
			}`,
			initEnv: map[string]any{},
			want:    42,
			err:     nil,
		},
		{
			desc:    "call_undefined_function",
			input:   "undefined(42)",
			code:    "",
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(IDENTIFIER, "undefined", nil, 1, 0), "Undefined variable 'undefined'."),
		},
		{
			desc:    "call_non_function",
			input:   "notfunc(42)",
			code:    "",
			initEnv: map[string]any{"notfunc": "string"},
			want:    nil,
			err:     NewRuntimeError(newToken(RIGHT_PAREN, ")", nil, 1, 7), "Can only call functions and classes."),
		},
		{
			desc:  "wrong_arity",
			input: "test(42, 43)",
			code: `fn test(x) {
				return x;
			}`,
			initEnv: map[string]any{},
			want:    nil,
			err:     NewRuntimeError(newToken(RIGHT_PAREN, ")", nil, 1, 11), "Expected 1 arguments but got 2."),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			// Parse and execute function declaration if present
			if tC.code != "" {
				scanner := NewScanner(nil, []byte(tC.code))
				tokens, err := scanner.ScanTokens()
				if err != nil {
					t.Fatal(err)
				}

				parser := NewParser(nil, tokens)
				stmts, err := parser.Parse()
				if err != nil {
					t.Fatal(err)
				}

				resolver := NewResolver(nil, interpreter)
				err = resolver.Resolve(stmts)
				if err != nil {
					t.Fatal(err)
				}

				for _, stmt := range stmts {
					err = interpreter.execute(stmt)
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			// Parse and execute function call
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser(nil, tokens)
			expr, err := parser.expression()
			if err != nil {
				t.Fatal(err)
			}

			resolver := NewResolver(nil, interpreter)
			_, err = resolver.resolveExpr(expr)
			if err != nil {
				t.Fatal(err)
			}

			got, err := interpreter.visitCallExpr(expr.(callExpr))
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
		input       string
		wantEnvVal  any
		wantEnvName string
		err         error
	}{
		{
			desc:        "without_initializer",
			input:       "var x;",
			wantEnvVal:  nil,
			wantEnvName: "x",
			err:         nil,
		},
		{
			desc:        "with_initializer",
			input:       "var y = 42;",
			wantEnvVal:  42,
			wantEnvName: "y",
			err:         nil,
		},
		{
			desc:        "with_string_initializer",
			input:       `var z = "hello";`,
			wantEnvVal:  "hello",
			wantEnvName: "z",
			err:         nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			stmts, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			err = interpreter.visitVarStmt(stmts[0].(varStmt))

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
		input   string
		initEnv map[string]any
		wantEnv map[string]any
		err     error
	}{
		{
			desc: "access_and_modify_global_var",
			input: `{
				var local = 42;
				global = 100;
			}`,
			initEnv: map[string]any{"global": 50.0},
			wantEnv: map[string]any{"global": 100},
			err:     nil,
		},
		{
			desc: "nested_blocks_access_global",
			input: `{
				var a = 1;
				{
					global = 200;
				}
			}`,
			initEnv: map[string]any{"global": 100.0},
			wantEnv: map[string]any{"global": 200},
			err:     nil,
		},
		{
			desc: "local_var_not_accessible_after_block",
			input: `{
				var local = 42;
			}`,
			initEnv: map[string]any{},
			wantEnv: map[string]any{},
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			stmts, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err = interpreter.visitBlockStmt(stmts[0].(blockStmt))
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
			for _, stmt := range stmts[0].(blockStmt).statements {
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
		input   string
		initEnv map[string]any
		wantVal any
		wantErr error
	}{
		{
			desc:    "return_literal_number",
			input:   "return 42;",
			initEnv: map[string]any{},
			wantVal: 42,
			wantErr: nil,
		},
		{
			desc:    "return_binary_expr",
			input:   "return 10 + 5;",
			initEnv: map[string]any{},
			wantVal: 15.0,
			wantErr: nil,
		},
		{
			desc:    "return_variable",
			input:   "return x;",
			initEnv: map[string]any{"x": "hello"},
			wantVal: "hello",
			wantErr: nil,
		},
		{
			desc:    "return_nil",
			input:   "return nil;",
			initEnv: map[string]any{},
			wantVal: nil,
			wantErr: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			stmt, err := parser.statement()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err = interpreter.visitReturnStmt(stmt.(returnStmt))
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
		desc    string
		input   string
		initEnv map[string]any
		wantEnv map[string]any
		err     error
	}{
		{
			desc:    "false_condition_no_iteration",
			input:   "while (false) { x = 1; }",
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 0.0},
			err:     nil,
		},
		{
			desc:    "counter_loop",
			input:   "while (counter < 3) { counter = counter + 1; sum = sum + 1; }",
			initEnv: map[string]any{"counter": 0.0, "sum": 0.0},
			wantEnv: map[string]any{"counter": 3.0, "sum": 3.0},
			err:     nil,
		},
		{
			desc:    "non_boolean_condition",
			input:   "while (cond) { x = 1; cond = false; }",
			initEnv: map[string]any{"cond": "not a boolean", "x": 0.0},
			wantEnv: map[string]any{"cond": false, "x": 1},
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()
			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}
			parser := NewParser(nil, tokens)
			stmt, err := parser.statement()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err = interpreter.visitWhileStmt(stmt.(whileStmt))
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
		desc    string
		input   string
		initEnv map[string]any
		wantEnv map[string]any
		err     error
	}{
		{
			desc:    "true_condition_no_else",
			input:   "if (true) { x = 1; }",
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 1},
			err:     nil,
		},
		{
			desc:    "false_condition_no_else",
			input:   "if (false) { x = 1; }",
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 0.0},
			err:     nil,
		},
		{
			desc:    "true_condition_with_else",
			input:   "if (true) { x = 1; } else { x = 2; }",
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 1},
			err:     nil,
		},
		{
			desc:    "false_condition_with_else",
			input:   "if (false) { x = 1; } else { x = 2; }",
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": 2},
			err:     nil,
		},
		{
			desc:    "non_boolean_condition",
			input:   `if ("not a boolean") { x = "hit then"; }`,
			initEnv: map[string]any{"x": 0.0},
			wantEnv: map[string]any{"x": "hit then"},
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser(nil, tokens)
			stmts, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err = interpreter.visitIfStmt(stmts[0].(ifStmt))
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
		input   string
		initEnv map[string]any
		want    any
		err     error
	}{
		{
			desc: "basic_function_declaration",
			input: `fn add(a, b) {
				return a + b;
			}`,
			initEnv: map[string]any{},
			want:    nil,
			err:     nil,
		},
		{
			desc:    "empty_function_declaration",
			input:   "fn empty() {}",
			initEnv: map[string]any{},
			want:    nil,
			err:     nil,
		},
		{
			desc:    "redefined_function",
			input:   "fn existing() {}",
			initEnv: map[string]any{"existing": "some_value"},
			want:    nil,
			err:     nil,
		},
		{
			desc: "fibonacci_function",
			input: `fn fib(n) {
				if (n <= 1) {
					return n;
				}
				return fib(n - 2) + fib(n - 1);
			}`,
			initEnv: map[string]any{},
			want:    nil,
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser(nil, tokens)
			stmts, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err = interpreter.visitFunctionStmt(stmts[0].(functionStmt))
			if tC.err != nil {
				assert.EqualError(t, err, tC.err.Error())
			} else {
				assert.NoError(t, err)
				val, exists := interpreter.env.values[stmts[0].(functionStmt).name.lexeme]
				assert.True(t, exists)
				_, ok := val.(function)
				assert.True(t, ok)
			}
		})
	}
}
func Test_interpretClassStmt(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		initEnv map[string]any
		err     error
	}{
		{
			desc:    "basic_class_declaration",
			input:   `class Test {}`,
			initEnv: map[string]any{},
			err:     nil,
		},
		{
			desc: "class_with_methods",
			input: `class Test {
    method() {}
    anotherMethod() {}
   }`,
			initEnv: map[string]any{},
			err:     nil,
		},
		{
			desc:    "redefined_class",
			input:   `class Existing {}`,
			initEnv: map[string]any{"Existing": "some_value"},
			err:     nil,
		},
		{
			desc: "class_with_method_return",
			input: `class Test {
    method() {
     return "test";
    }
   }`,
			initEnv: map[string]any{},
			err:     nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			scanner := NewScanner(nil, []byte(tC.input))
			tokens, err := scanner.ScanTokens()
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser(nil, tokens)
			stmts, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			interpreter := NewInterpreter(nil)
			for k, v := range tC.initEnv {
				interpreter.env.define(k, v)
			}

			err = interpreter.visitClassStmt(stmts[0].(classStmt))
			if tC.err != nil {
				assert.EqualError(t, err, tC.err.Error())
			} else {
				assert.NoError(t, err)
				val, exists := interpreter.env.values[stmts[0].(classStmt).name.lexeme]
				assert.True(t, exists)
				_, ok := val.(class)
				assert.True(t, ok)
			}
		})
	}
}
