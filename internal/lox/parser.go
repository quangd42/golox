package lox

import (
	"errors"
	"fmt"
	"slices"
)

type Parser struct {
	tokens  []token
	current int
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (expr, error) {
	return p.expression()
}

func (p *Parser) expression() (expr, error) {
	return p.equality()
}

func (p *Parser) equality() (expr, error) {
	// TODO:this should be comparision()
	out, err := p.factor()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ... some more levels

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (expr, error) {
	out, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(SLASH, STAR) {
		oper, err := p.advance()
		if err != nil {
			return nil, err
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		out = binaryExpr{
			left:     out,
			operator: oper,
			right:    right,
		}
	}
	return out, nil
}

// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() (expr, error) {
	if p.match(BANG, MINUS) {
		oper, err := p.advance()
		if err != nil {
			return nil, err
		}
		next, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unaryExpr{operator: oper, right: next}, nil
	}
	return p.primary()
}

// TODO: reduce duplicate p.advance() in each case

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() (expr, error) {
	switch {
	case p.match(TRUE):
		p.advance()
		return literalExpr{true}, nil
	case p.match(FALSE):
		p.advance()
		return literalExpr{false}, nil
	case p.match(NIL):
		p.advance()
		return literalExpr{nil}, nil
	case p.match(NUMBER, STRING):
		out, err := p.advance()
		if err != nil {
			return nil, err
		}
		return literalExpr{out.literal}, nil
	case p.match(LEFT_PAREN):
		p.advance()
		err := p.consumeForward(RIGHT_PAREN)
		if err != nil {
			return nil, errors.New("expect ')' after expression")
		}
		out, err := p.expression()
		if err != nil {
			return nil, err
		}
		return groupingExpr{out}, nil
	default:
		return nil, fmt.Errorf("invalid token processed by primary(): %v", p.peek())
	}
}

// advance **consumes** the current token and returns it
func (p *Parser) advance() (token, error) {
	if p.isAtEnd() {
		return token{}, ErrEOF
	}
	out := p.peek()
	p.current++
	return out, nil
}

// PERF: this is not effiecient
func (p *Parser) consumeForward(tokenType tokenType) error {
	current := p.current
	for !p.isAtEnd() {
		if p.match(tokenType) {
			p.tokens = slices.Delete(p.tokens, p.current, p.current+1)
			p.current = current
			return nil
		}
		p.current++
	}
	return ErrEOF
}

// match peeks at the current token to see if it is one of the expected tokens
func (p Parser) match(expected ...tokenType) bool {
	for _, tt := range expected {
		if p.peek().tokenType == tt {
			return true
		}
	}
	return false
}

// peek returns the current token without consuming it
func (p Parser) peek() token {
	if p.isAtEnd() {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current]
}

// isAtEnd returns whether there is more token to parse
func (p Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}
