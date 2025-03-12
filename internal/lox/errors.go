package lox

import (
	"fmt"
)

type BaseError struct {
	Line  int
	Msg   string
	Where string
}

func (e BaseError) Error() string {
	if e.Where == "" {
		return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Msg)
	}
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Line, e.Where, e.Msg)
}

func NewBaseError(line int, where, msg string) error {
	return BaseError{
		Line:  line,
		Msg:   msg,
		Where: where,
	}
}

func errUnsupportedCharacter(line int, c rune) error {
	return NewBaseError(line, string(c), "unsupported character")
}

func errUnterminatedString(line int) error {
	return NewBaseError(line, "", "unterminated string")
}

func errInvalidNumber(line int, lex string) error {
	return NewBaseError(line, lex, "invalid number")
}

type ParseError struct {
	Token token
	Msg   string
}

func (e ParseError) Error() string {
	switch e.Token.tokenType {
	case EOF:
		return fmt.Sprintf("[line %d] Error at end: %s", e.Token.line, e.Msg)
	default:
		return fmt.Sprintf("[line %d] Error at '%s': %s", e.Token.line, e.Token.lexeme, e.Msg)
	}
}

func NewParseError(t token, msg string) error {
	return ParseError{Token: t, Msg: msg}
}

type RuntimeError struct {
	Token token
	Msg   string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d] at '%s'", e.Msg, e.Token.line, e.Token.lexeme)
}

func NewRuntimeError(t token, msg string) error {
	return RuntimeError{Token: t, Msg: msg}
}
