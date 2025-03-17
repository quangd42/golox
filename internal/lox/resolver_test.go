package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolver_Resolve(t *testing.T) {
	tests := []struct {
		name            string
		input          []stmt
		expectedLocals map[expr]int
	}{
		{
			name: "simple variable declaration and usage",
			input: []stmt{
				varStmt{
					name: newToken(IDENTIFIER, "x", "x", 0),
					initializer: literalExpr{value: 1},
				},
				exprStmt{
					expr: variableExpr{
						name: newToken(IDENTIFIER, "x", "x", 0),
					},
				},
			},
			expectedLocals: map[expr]int{
				variableExpr{
					name: newToken(IDENTIFIER, "x", "x", 0),
				}: 0,
			},
		},
		{
			name: "nested scope variable resolution",
			input: []stmt{
				blockStmt{
					statements: []stmt{
						varStmt{
							name: newToken(IDENTIFIER, "x", "x", 0),
							initializer: literalExpr{value: 1},
						},
						exprStmt{
							expr: variableExpr{
								name: newToken(IDENTIFIER, "x", "x", 0),
							},
						},
					},
				},
			},
			expectedLocals: map[expr]int{
				variableExpr{
					name: newToken(IDENTIFIER, "x", "x", 0),
				}: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interpreter := NewInterpreter()
			resolver := NewResolver(interpreter)

			err := resolver.Resolve(tt.input)
			assert.NoError(t, err)

			// Compare the interpreter's locals with expected values
			assert.Equal(t, tt.expectedLocals, interpreter.locals)
		})
	}
}
