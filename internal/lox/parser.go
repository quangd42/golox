package lox

import (
	"fmt"
	"slices"
)

type Parser struct {
	tokens  []token
	current int
}

func NewParser(tokens []token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

/*
	Statements
*/

// program → declaration* EOF ;
func (p *Parser) Parse() ([]stmt, error) {
	out := make([]stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			continue
		}
		out = append(out, stmt)
	}
	return out, nil
}

// declaration → varDecl | statement ;
func (p *Parser) declaration() (out stmt, err error) {
	if p.match(VAR) {
		out, err = p.varDecl()
	} else {
		out, err = p.statement()
	}
	// print error and synchronize at statement level
	if err != nil {
		fmt.Println(err.Error())
		p.synchronize()
		return nil, err
	}
	return out, nil
}

// varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDecl() (stmt, error) {
	p.advance()
	name, ok := p.matchConsume(IDENTIFIER)
	if !ok {
		return nil, NewParseError(p.peek(), "Expect variable name.")
	}
	var initializer expr
	var err error
	if p.match(EQUAL) {
		p.advance()
		initializer, err = p.expression()
		if err != nil {
			return nil, NewParseError(p.peek(), "Expect expression.")
		}
	}
	if _, ok := p.matchConsume(SEMICOLON); !ok {
		return nil, NewParseError(p.peek(), "Expect ';' after variable declaration.")
	}
	return varStmt{name: name, initializer: initializer}, nil
}

// statement → exprStmt | printStmt | block ;
func (p *Parser) statement() (stmt, error) {
	switch {
	case p.match(PRINT):
		p.advance()
		return p.printStmt()
	case p.match(LEFT_BRACE):
		p.advance()
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return blockStmt{stmts}, nil
	default:
		return p.exprStmt()
	}
}

// exprStmt → expression ";" ;
func (p *Parser) exprStmt() (stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, ok := p.matchConsume(SEMICOLON); !ok {
		return nil, NewParseError(p.peek(), "Expect ';' after expression.")
	}
	return exprStmt{expr: expr}, nil
}

// printStmt → "print" expression ";" ;
func (p *Parser) printStmt() (stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, ok := p.matchConsume(SEMICOLON); !ok {
		return nil, NewParseError(p.peek(), "Expect ';' after expression.")
	}
	return printStmt{expr: expr}, nil
}

// block → "{" declaration* "}" ;
func (p *Parser) block() ([]stmt, error) {
	out := make([]stmt, 0)
	for !p.match(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			continue
		}
		out = append(out, stmt)
	}
	if _, ok := p.matchConsume(RIGHT_BRACE); !ok {
		return nil, NewParseError(p.peek(), "Expect '}' after block.")
	}
	return out, nil
}

/*
	Expressions
*/

// expression → assignment ( "," assignment )* ;
func (p *Parser) expression() (expr, error) {
	out, err := p.assignment()
	if err != nil {
		return nil, err
	}
	for p.match(COMMA) {
		oper, err := p.advance()
		if err != nil {
			return nil, err
		}
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}
		out = binaryExpr{left: out, operator: oper, right: right}
	}
	return out, nil
}

// assignment → IDENTIFIER "=" assignment | ternary ;
func (p *Parser) assignment() (expr, error) {
	out, err := p.ternary()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		tok, _ := p.advance() // this must succeed in a if p.match()
		val, err := p.assignment()
		if err != nil {
			return nil, err
		}
		varExpr, ok := out.(variableExpr)
		if !ok {
			return nil, NewRuntimeError(tok, "Invalid assignment target.")
		}
		out = assignExpr{name: varExpr.name, value: val}
	}
	return out, nil
}

// ternary → equality ( "?" ternary ":" ternary )? ;
func (p *Parser) ternary() (expr, error) {
	out, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match(QUESTION) {
		lOper, err := p.advance()
		if err != nil {
			return nil, err
		}
		trueExpr, err := p.ternary()
		if err != nil {
			return nil, err
		}
		rOper, err := p.advance()
		if err != nil {
			return nil, err
		}
		if !rOper.hasType(COLON) {
			return nil, NewParseError(lOper, "expect ':' after expression")
		}
		falseExpr, err := p.ternary()
		if err != nil {
			return nil, err
		}
		out = binaryExpr{
			left: binaryExpr{
				left:     out,
				operator: lOper,
				right:    trueExpr,
			},
			operator: rOper,
			right:    falseExpr,
		}
	}
	return out, nil
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
		out = binaryExpr{left: out, operator: oper, right: right}
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
		out = binaryExpr{left: out, operator: oper, right: right}
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
		out = binaryExpr{left: out, operator: oper, right: right}
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
		out = binaryExpr{left: out, operator: oper, right: right}
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

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
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
	case tok.hasType(IDENTIFIER):
		return variableExpr{tok}, nil
	case tok.hasType(NUMBER, STRING):
		return literalExpr{tok.literal}, nil
	case tok.hasType(LEFT_PAREN):
		// TODO: parse empty group ()
		out, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, ok := p.matchConsume(RIGHT_PAREN); !ok {
			return nil, NewParseError(tok, "expect ')' after expression")
		}
		return groupingExpr{out}, nil
	case tok.hasType(SLASH, STAR, MINUS, PLUS, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, BANG, BANG_EQUAL):
		_, err := p.expression()
		if err != nil {
			return nil, err
		}
		return nil, NewParseError(tok, "expect left operand")
	default:
		return nil, NewParseError(tok, "expect an expression")
	}
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		tok, err := p.advance()
		if err != nil {
			return
		}
		if tok.hasType(SEMICOLON) {
			return
		}
		if p.match(CLASS, FN, VAR, FOR, IF, WHILE, PRINT, RETURN) {
			return
		}
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

func (p *Parser) matchConsume(tokenType tokenType) (token, bool) {
	if p.match(tokenType) {
		out, err := p.advance()
		if err != nil {
			return token{}, false
		}
		return out, true
	}
	return token{}, false
}

// match peeks at the current token to see if it is one of the expected tokens
func (p Parser) match(expected ...tokenType) bool {
	return slices.Contains(expected, p.peek().tokenType)
}

// isAtEnd returns whether there is more token to parse
func (p Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

// peek returns the current token without consuming it. Returns the last token
// if there is no more token to peek at
func (p Parser) peek() token {
	if p.current >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current]
}
