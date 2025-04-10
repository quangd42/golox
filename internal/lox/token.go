package lox

import (
	"fmt"
	"slices"
)

type tokenType string

const (
	// Single-character tokens.

	LEFT_PAREN    tokenType = "("
	RIGHT_PAREN   tokenType = ")"
	LEFT_BRACKET  tokenType = "["
	RIGHT_BRACKET tokenType = "]"
	LEFT_BRACE    tokenType = "{"
	RIGHT_BRACE   tokenType = "}"
	COMMA         tokenType = ","
	COLON         tokenType = ":"
	DOT           tokenType = "."
	MINUS         tokenType = "-"
	PLUS          tokenType = "+"
	QUESTION      tokenType = "?"
	SEMICOLON     tokenType = ";"
	SLASH         tokenType = "/"
	STAR          tokenType = "*"

	// One or two character tokens.

	BANG          tokenType = "!"
	BANG_EQUAL    tokenType = "!="
	EQUAL         tokenType = "="
	EQUAL_EQUAL   tokenType = "=="
	GREATER       tokenType = ">"
	GREATER_EQUAL tokenType = ">="
	LESS          tokenType = "<"
	LESS_EQUAL    tokenType = "<="

	// Literals.

	IDENTIFIER tokenType = "IDENTIFIER"
	STRING     tokenType = "STRING"
	NUMBER     tokenType = "NUMBER"

	// Keywords.

	AND      tokenType = "and"
	CLASS    tokenType = "class"
	ELSE     tokenType = "else"
	FALSE    tokenType = "false"
	FN       tokenType = "fn"
	FOR      tokenType = "for"
	IF       tokenType = "if"
	NIL      tokenType = "nil"
	OR       tokenType = "or"
	PRINT    tokenType = "print"
	RETURN   tokenType = "return"
	SUPER    tokenType = "super"
	THIS     tokenType = "this"
	TRUE     tokenType = "true"
	VAR      tokenType = "var"
	WHILE    tokenType = "while"
	BREAK    tokenType = "break"
	CONTINUE tokenType = "continue"

	EOF tokenType = "EOF"
)

func (tt tokenType) String() string {
	return string(tt)
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

// newTokenNoLiteralType is a convenient function to create a token with no literal value
// this is mostly useful for tests
func newTokenNoLiteralType(tokenType tokenType, line, offset int) token {
	switch tokenType {
	case IDENTIFIER, STRING, NUMBER:
		return token{}
	case EOF:
		return newToken(tokenType, "", nil, line, offset)
	default:
		return newToken(tokenType, tokenType.String(), tokenType.String(), line, offset)
	}
}

func (t token) String() string {
	return fmt.Sprintf("%s %s", t.tokenType, t.lexeme)
}

// hasType returns whether the tokenType is one of the expected.
func (t token) hasType(expected ...tokenType) bool {
	return slices.Contains(expected, t.tokenType)
}

// lookupIdentifier returns the tokenType if the provided lexeme
// is a reserved keyword, and returns IDENTIFIER tokenType otherwise
func lookupIdentifier(lex string) tokenType {
	keywords := map[string]tokenType{
		"and":      AND,
		"class":    CLASS,
		"else":     ELSE,
		"false":    FALSE,
		"fn":       FN,
		"for":      FOR,
		"if":       IF,
		"nil":      NIL,
		"or":       OR,
		"print":    PRINT,
		"return":   RETURN,
		"super":    SUPER,
		"this":     THIS,
		"true":     TRUE,
		"var":      VAR,
		"while":    WHILE,
		"break":    BREAK,
		"continue": CONTINUE,
	}
	tt, ok := keywords[lex]
	if !ok {
		return IDENTIFIER
	}
	return tt
}
