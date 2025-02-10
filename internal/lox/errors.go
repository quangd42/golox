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
	return fmt.Sprintf("[line %d] Error%s: %s", e.Line, e.Where, e.Msg)
}

func NewLoxError(msg, where string, line int) error {
	return LoxError{
		Line:  line,
		Msg:   msg,
		Where: where,
	}
}

func errUnsupportedCharacter(c rune, line int) error {
	return NewLoxError(fmt.Sprintf("unsupported character: %s", string(c)), "", line)
}

func errUnterminatedString(line int) error {
	return NewLoxError("unterminated string", "", line)
}

func errInvalidNumber(line int) error {
	return NewLoxError("invalid number", "", line)
}
