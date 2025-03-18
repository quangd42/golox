package lox

import (
	"errors"
	"fmt"
	"slices"
)

type tokenType int

const (
	// Single-character tokens.

	LEFT_PAREN tokenType = iota + 1 // reserve 0 for empty value
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	COLON
	DOT
	MINUS
	PLUS
	QUESTION
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.

	IDENTIFIER
	STRING
	NUMBER

	// Keywords.

	AND
	CLASS
	ELSE
	FALSE
	FN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

func (tt tokenType) String() string {
	names := []string{
		// Emtpy value
		"",

		// Single-character tokens.

		"LEFT_PAREN",
		"RIGHT_PAREN",
		"LEFT_BRACE",
		"RIGHT_BRACE",
		"COMMA",
		"COLON",
		"DOT",
		"MINUS",
		"PLUS",
		"QUESTION",
		"SEMICOLON",
		"SLASH",
		"STAR",

		// One or two character tokens.

		"BANG",
		"BANG_EQUAL",
		"EQUAL,",
		"EQUAL_EQUAL",
		"GREATER",
		"GREATER_EQUAL",
		"LESS",
		"LESS_EQUAL",

		// Literals.

		"IDENTIFIER",
		"STRING",
		"NUMBER",

		// Keywords.

		"AND",
		"CLASS",
		"ELSE",
		"FALSE",
		"FN",
		"FOR",
		"IF",
		"NIL",
		"OR",

		"PRINT",
		"RETURN",
		"SUPER",
		"THIS",
		"TRUE",
		"VAR",
		"WHILE",

		"EOF",
	}
	return names[int(tt)]
}

type token struct {
	tokenType tokenType
	lexeme    string
	literal   any
	line      int
	offset    int
}

func newToken(tokenType tokenType, lexeme string, literal any, line, offset int) token {
	return token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
		offset:    offset,
	}
}

// newTokenNoLiteral is a convenient function to create a token with no literal value
// this is mostly useful for tests
func newTokenNoLiteral(tokenType tokenType) token {
	return newToken(tokenType, tokenType.String(), nil, 0, 0)
}

func (t token) String() string {
	return fmt.Sprintf("%s %s", t.tokenType, t.lexeme)
}

// hasType returns whether the tokenType is one of the expected.
func (t token) hasType(expected ...tokenType) bool {
	return slices.Contains(expected, t.tokenType)
}

// getKeywords returns the TokenType if the provided lexeme
// is a reserved keyword, and returns an error otherwise
func getKeywords(lex string) (tokenType, error) {
	keywords := map[string]tokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"fn":     FN,
		"for":    FOR,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
	tt, ok := keywords[lex]
	if !ok {
		return 0, errors.New("not a keyword")
	}
	return tt, nil
}
