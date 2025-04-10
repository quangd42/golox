package lox

import (
	"testing"
)

func Test_function_call(t *testing.T) {
	tests := []struct {
		name     string
		fn       function
		args     []any
		want     any
		wantErr  bool
		retValue *functionReturn
		err      error
	}{
		{
			name: "simple function with no return",
			fn: function{
				name: newToken(IDENTIFIER, "test", "test", 0, 1),
				literal: functionExpr{
					params: []token{},
					body:   []stmt{},
				},
			},
			args:    []any{},
			want:    nil,
			wantErr: false,
			err:     nil,
		},
		{
			name: "function with return value",
			fn: function{
				name: newToken(IDENTIFIER, "test", "test", 0, 1),
				literal: functionExpr{
					params: []token{},
					body: []stmt{
						returnStmt{
							value: literalExpr{value: "hello"},
						},
					},
				},
			},
			args:     []any{},
			want:     "hello",
			wantErr:  false,
			retValue: &functionReturn{value: "hello"},
			err:      nil,
		},
		{
			name: "function with parameters",
			fn: function{
				name: newToken(IDENTIFIER, "test", "test", 0, 1),
				literal: functionExpr{
					params: []token{
						newToken(IDENTIFIER, "x", "x", 0, 1),
						newToken(IDENTIFIER, "y", "y", 0, 1),
					},
					body: []stmt{},
				},
			},
			args:    []any{1, 2},
			want:    nil,
			wantErr: false,
			err:     nil,
		},
		{
			name: "runtime error in function body",
			fn: function{
				name: newToken(IDENTIFIER, "test", "test", 0, 1),
				literal: functionExpr{
					params: []token{},
					body: []stmt{
						exprStmt{
							expr: binaryExpr{
								left:     literalExpr{value: 1},
								operator: newTokenNoLiteralType(MINUS, 1, 32),
								right:    literalExpr{value: "string"}, // Invalid operation
							},
						},
					},
				},
			},
			args:    []any{},
			want:    nil,
			wantErr: true,
			err:     NewRuntimeError(newToken(IDENTIFIER, "test", "test", 0, 1), "Operand must be a number."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter(nil)
			got, err := tt.fn.call(interpreter, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("function.call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("function.call() = %v, want %v", got, tt.want)
			}
		})
	}
}
