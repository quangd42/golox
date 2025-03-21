package lox

import (
	"fmt"
)

type ErrorReporter interface {
	report(line int, where, msg string)
	HadError() bool
	HadRuntimeError() bool
	ResetError()
	ResetRuntimeError()
	ScanError(line int, msg string)
	ParseError(token token, msg string) ParseError
	RuntimeError(e RuntimeError)
}

type LoxErrorReporter struct {
	hadError        bool
	hadRuntimeError bool
}

func NewLoxErrorReporter() *LoxErrorReporter {
	return &LoxErrorReporter{}
}

func (l *LoxErrorReporter) report(line int, where, msg string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, msg)
	l.hadError = true
}

func (l *LoxErrorReporter) ScanError(line int, msg string) {
	l.report(line, "", msg)
}

func (l *LoxErrorReporter) HadError() bool {
	return l.hadError
}

func (l *LoxErrorReporter) ResetError() {
	l.hadError = false
}

func (l *LoxErrorReporter) ParseError(token token, msg string) ParseError {
	switch token.tokenType {
	case EOF:
		l.report(token.line, " at end", msg)
	default:
		l.report(token.line, fmt.Sprintf(" at '%s'", token.lexeme), msg)
	}

	return ParseError{Token: token, Msg: msg}
}

type ParseError struct {
	Token token
	Msg   string
}

func NewParseError(token token, msg string) ParseError {
	return ParseError{token, msg}
}

func (e ParseError) Error() string {
	switch e.Token.tokenType {
	case EOF:
		return fmt.Sprintf("[line %d] Error at end: %s", e.Token.line, e.Msg)
	default:
		return fmt.Sprintf("[line %d] Error at '%s': %s", e.Token.line, e.Token.lexeme, e.Msg)
	}
}

type RuntimeError struct {
	Token token
	Msg   string
}

func NewRuntimeError(token token, msg string) RuntimeError {
	return RuntimeError{token, msg}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Runtime Error at '%s': %s", e.Token.line, e.Token.lexeme, e.Msg)
}

func (l *LoxErrorReporter) RuntimeError(err RuntimeError) {
	l.hadRuntimeError = true
	fmt.Println(err.Error())
}

func (l *LoxErrorReporter) HadRuntimeError() bool {
	return l.hadRuntimeError
}

func (l *LoxErrorReporter) ResetRuntimeError() {
	l.hadRuntimeError = false
}
