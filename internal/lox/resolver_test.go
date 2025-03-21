package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedLocals map[expr]int
	}{
		{
			name:  "simple variable declaration and usage",
			input: `{var x = 1; print x;}`,
			expectedLocals: map[expr]int{
				variableExpr{
					name: newToken(IDENTIFIER, "x", "x", 1, 18),
				}: 0,
			},
		},
		{
			name:  "nested scope variable resolution",
			input: `{ var x = 1; { print x; } }`,
			expectedLocals: map[expr]int{
				variableExpr{
					name: newToken(IDENTIFIER, "x", "x", 1, 21),
				}: 1,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scanner := NewScanner(nil, []byte(tt.input))
			tokens, err := scanner.ScanTokens()
			assert.NoError(t, err)

			parser := NewParser(nil, tokens)
			stmts, err := parser.Parse()
			assert.NoError(t, err)

			interpreter := NewInterpreter(nil)
			resolver := NewResolver(nil, interpreter)

			err = resolver.Resolve(stmts)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedLocals, interpreter.locals)
		})
	}
}
