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
		retValue *returnValue
		err      error
	}{
		{
			name: "simple function with no return",
			fn: function{
				declaration: functionStmt{
					name:   newToken(IDENTIFIER, "test", "test", 0),
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
				declaration: functionStmt{
					name:   newToken(IDENTIFIER, "test", "test", 0),
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
			retValue: &returnValue{value: "hello"},
			err:      nil,
		},
		{
			name: "function with parameters",
			fn: function{
				declaration: functionStmt{
					name: newToken(IDENTIFIER, "test", "test", 0),
					params: []token{
						newToken(IDENTIFIER, "x", "x", 0),
						newToken(IDENTIFIER, "y", "y", 0),
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
				declaration: functionStmt{
					name:   newToken(IDENTIFIER, "test", "test", 0),
					params: []token{},
					body: []stmt{
						exprStmt{
							expr: binaryExpr{
								left:     literalExpr{value: 1},
								operator: newTokenNoLiteral(MINUS),
								right:    literalExpr{value: "string"}, // Invalid operation
							},
						},
					},
				},
			},
			args:    []any{},
			want:    nil,
			wantErr: true,
			err:     NewRuntimeError(newToken(IDENTIFIER, "test", "test", 0), "Operand must be a number."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
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
