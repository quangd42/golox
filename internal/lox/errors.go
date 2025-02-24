package lox

import (
	"fmt"
)

type LoxError struct {
	Line  int
	Msg   string
	Where string
}

func (e LoxError) Error() string {
	if e.Where == "" {
		return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Msg)
	}
	return fmt.Sprintf("[line %d] Error at %s: %s", e.Line, e.Where, e.Msg)
}

func NewLoxError(line int, where, msg string) error {
	return LoxError{
		Line:  line,
		Msg:   msg,
		Where: where,
	}
}

func errUnsupportedCharacter(line int, c rune) error {
	return NewLoxError(line, string(c), "unsupported character")
}

func errUnterminatedString(line int) error {
	return NewLoxError(line, "", "unterminated string")
}

func errInvalidNumber(line int, lex string) error {
	return NewLoxError(line, lex, "invalid number")
}
