package lox

import (
	"errors"
	"fmt"
)

type TokenType int

const (
	// Single-character tokens.

	LEFT_PAREN TokenType = iota + 1 // reserve 0 for empty value
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
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

func (tt TokenType) String() string {
	names := []string{
		// Emtpy value
		"",

		// Single-character tokens.

		"LEFT_PAREN",
		"RIGHT_PAREN",
		"LEFT_BRACE",
		"RIGHT_BRACE",
		"COMMA",
		"DOT",
		"MINUS",
		"PLUS",
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

type Token struct {
	Type     TokenType
	Lexeme   string
	Literals any
	Line     int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s", t.Type, t.Lexeme)
}

// getKeywords returns the TokenType if the provided lexeme
// is a reserved keyword, and returns an error otherwise
func getKeywords(lex string) (TokenType, error) {
	keywords := map[string]TokenType{
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
