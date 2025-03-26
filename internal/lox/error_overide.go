// These types are used to taken advantage of error control flow and exit early.

package lox

import "fmt"

type functionReturn struct {
	value any
}

func (fr *functionReturn) Error() string {
	return fmt.Sprintf("%s", fr.value)
}

type loopBreak struct {
	keyword token
	label   token
}

func (lb *loopBreak) Error() string {
	return fmt.Sprintf("%s on line %d", lb.keyword.lexeme, lb.keyword.line)
}

type loopContinue struct {
	keyword token
	label   token
}

func (lc *loopContinue) Error() string {
	return fmt.Sprintf("%s on line %d", lc.keyword.lexeme, lc.keyword.line)
}
