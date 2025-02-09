package lox

import "fmt"

type TokenType int

const (
	// Single-character tokens.

	LEFT_PAREN TokenType = iota
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
	FUN
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
		"FUN",
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
	Type   TokenType
	Lexeme string
	Line   int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s", t.Type, t.Lexeme)
}
