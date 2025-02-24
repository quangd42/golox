package lox

import (
	"fmt"
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

// expression → equality ;
func (p *Parser) expression() (expr, error) {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (expr, error) {
	out, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		oper, err := p.advance()
		if err != nil {
			return nil, err
		}
		right, err := p.comparison()
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

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (expr, error) {
	out, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		oper, err := p.advance()
		if err != nil {
			return nil, err
		}
		right, err := p.term()
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

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (expr, error) {
	out, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(MINUS, PLUS) {
		oper, err := p.advance()
		if err != nil {
			return nil, err
		}
		right, err := p.factor()
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

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() (expr, error) {
	tok, err := p.advance()
	if err != nil {
		return nil, err
	}
	switch {
	case tok.hasType(TRUE):
		return literalExpr{true}, nil
	case tok.hasType(FALSE):
		return literalExpr{false}, nil
	case tok.hasType(NIL):
		return literalExpr{nil}, nil
	case tok.hasType(NUMBER, STRING):
		return literalExpr{tok.literal}, nil
	case tok.hasType(LEFT_PAREN):
		// TODO: parse empty group ()
		out, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.matchConsume(RIGHT_PAREN) {
			return nil, NewLoxError(tok.line, fmt.Sprintf("'%s'", tok.lexeme), "expect ')' after expression")
		}
		return groupingExpr{out}, nil
	default:
		return nil, NewLoxError(tok.line, fmt.Sprintf("'%s'", tok.lexeme), "expected an expression")
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

func (p *Parser) matchConsume(tokenType tokenType) bool {
	if p.match(tokenType) {
		p.advance()
		return true
	}
	return false
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

// peek returns the current token without consuming it. Returns the last token
// if there is no more token to peek at
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
