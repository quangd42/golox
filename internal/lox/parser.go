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
	_, err := p.consume(VAR, "Expect 'var' at the beginning of variable declaration.")
	if err != nil {
		return nil, err
	}
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer expr
	if p.match(EQUAL) {
		p.advance()
		initializer, err = p.expression()
		if err != nil {
			return nil, NewParseError(p.peek(), "Expect expression.")
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return varStmt{name: name, initializer: initializer}, nil
}

// statement → exprStmt | forStmt | ifStmt | printStmt | whileStmt | block ;
func (p *Parser) statement() (stmt, error) {
	switch {
	case p.match(FOR):
		return p.forStatement()
	case p.match(IF):
		return p.ifStatement()
	case p.match(PRINT):
		return p.printStatement()
	case p.match(WHILE):
		return p.whileStatement()
	case p.match(LEFT_BRACE):
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return blockStmt{stmts}, nil
	default:
		return p.exprStatement()
	}
}

// exprStmt → expression ";" ;
func (p *Parser) exprStatement() (stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return exprStmt{expr: expr}, nil
}

// forStmt → "for" (( varDecl | exprStmt | ";" ) expression? ";" expression?)? block ;
func (p *Parser) forStatement() (stmt, error) {
	var err error
	if _, err := p.consume(FOR, "Expect loop."); err != nil {
		return nil, err
	}
	var cond expr
	if p.match(LEFT_BRACE) {
		cond = literalExpr{true}
		bodyStmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return whileStmt{condition: cond, body: blockStmt{bodyStmts}}, nil
	}

	var initializer stmt
	if p.match(VAR) {
		initializer, err = p.varDecl()
	} else if !p.match(SEMICOLON) {
		initializer, err = p.exprStatement()
	} else {
		p.advance() // consume ';'
	}
	if err != nil {
		return nil, err
	}

	if !p.match(SEMICOLON) {
		cond, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after condition."); err != nil {
		return nil, err
	}

	var inc expr
	if !p.match(LEFT_BRACE) {
		inc, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	bodyStmts, err := p.block()
	if err != nil {
		return nil, err
	}

	if inc != nil {
		bodyStmts = append(bodyStmts, exprStmt{inc})
	}
	if cond == nil {
		cond = literalExpr{true}
	}
	var out stmt = whileStmt{condition: cond, body: blockStmt{bodyStmts}}
	if initializer != nil {
		out = blockStmt{[]stmt{initializer, out}}
	}

	return out, nil
}

// ifStmt → "if" expression block ( "else" block )? ;
func (p *Parser) ifStatement() (stmt, error) {
	if _, err := p.consume(IF, "Expect if statement."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	thenStmts, err := p.block()
	if err != nil {
		return nil, err
	}
	var elseBlock stmt
	if p.match(ELSE) {
		p.advance()
		elseStmts, err := p.block()
		if err != nil {
			return nil, err
		}
		elseBlock = blockStmt{elseStmts}
	}
	return ifStmt{
		condition:  condition,
		thenBranch: blockStmt{thenStmts},
		elseBranch: elseBlock,
	}, nil
}

// printStmt → "print" expression ";" ;
func (p *Parser) printStatement() (stmt, error) {
	if _, err := p.consume(PRINT, "Expect print statement."); err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return printStmt{expr: expr}, nil
}

// whileStmt → "while" expression block ;
func (p *Parser) whileStatement() (stmt, error) {
	if _, err := p.consume(WHILE, "Expect loop."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	stmts, err := p.block()
	if err != nil {
		return nil, err
	}
	return whileStmt{condition: condition, body: blockStmt{stmts}}, nil
}

// block → "{" declaration* "}" ;
func (p *Parser) block() ([]stmt, error) {
	if _, err := p.consume(LEFT_BRACE, "Expect block."); err != nil {
		return nil, err
	}
	out := make([]stmt, 0)
	for !p.match(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			continue
		}
		out = append(out, stmt)
	}
	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
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
		oper, _ := p.advance()
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}
		out = binaryExpr{left: out, operator: oper, right: right}
	}
	return out, nil
}

// assignment → IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (expr, error) {
	out, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		tok, _ := p.advance()
		val, err := p.assignment()
		if err != nil {
			return nil, err
		}
		varExpr, ok := out.(variableExpr)
		if !ok {
			return nil, NewParseError(tok, "Invalid assignment target.")
		}
		out = assignExpr{name: varExpr.name, value: val}
	}
	return out, nil
}

// logic_or → logic_and ( "or" logic_and )* ;
func (p *Parser) or() (expr, error) {
	out, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(OR) {
		tok, _ := p.advance()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		out = logicalExpr{left: out, operator: tok, right: right}
	}
	return out, nil
}

// logic_and → ternary ( "and" ternary )* ;
func (p *Parser) and() (expr, error) {
	out, err := p.ternary()
	if err != nil {
		return nil, err
	}
	for p.match(AND) {
		tok, _ := p.advance()
		right, err := p.ternary()
		if err != nil {
			return nil, err
		}
		out = logicalExpr{left: out, operator: tok, right: right}
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
			return nil, NewParseError(lOper, "Expect ':' after expression.")
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
		oper, _ := p.advance()
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
		oper, _ := p.advance()
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
		oper, _ := p.advance()
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
		oper, _ := p.advance()
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
		out, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return groupingExpr{out}, nil
	case tok.hasType(SLASH, STAR, MINUS, PLUS, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, BANG, BANG_EQUAL):
		_, err := p.expression()
		if err != nil {
			return nil, err
		}
		return nil, NewParseError(tok, "Expect left operand.")
	default:
		return nil, NewParseError(tok, "Expect an expression.")
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

// consume **consumes** the current token if it matches expected token.
// Returns the consumed token if matched.
func (p *Parser) consume(expected tokenType, errMsg string) (token, error) {
	if p.peek().tokenType != expected {
		return token{}, NewParseError(p.peek(), errMsg)
	}
	out, err := p.advance()
	if err != nil {
		return token{}, err
	}
	return out, nil
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
