package lox

import (
	"fmt"
	"slices"
)

type Parser struct {
	er      ErrorReporter
	tokens  []token
	current int
}

func NewParser(er ErrorReporter, tokens []token) *Parser {
	return &Parser{er: er, tokens: tokens, current: 0}
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

// declaration → classDecl | fnDecl | varDecl | statement ;
func (p *Parser) declaration() (out stmt, err error) {
	switch {
	case p.match(CLASS):
		out, err = p.classDecl()
	case p.match(VAR):
		out, err = p.varDecl()
	case p.match(FN):
		out, err = p.function(fnTypeFUNCTION)
	default:
		out, err = p.statement()
	}
	// synchronize at statement level
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return out, nil
}

// classDecl → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;
func (p *Parser) classDecl() (stmt, error) {
	_, err := p.consume(CLASS, "Expect 'class' at the beginning of variable declaration.")
	if err != nil {
		return nil, err
	}
	name, err := p.consume(IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}
	var superclass variableExpr
	if p.match(LESS) {
		p.advance()
		superclassTok, err := p.consume(IDENTIFIER, "Expect superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = variableExpr{superclassTok}
	}
	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}
	methods := make([]functionStmt, 0)
	for !p.match(RIGHT_BRACE) && !p.isAtEnd() {
		method, err := p.function(fnTypeMETHOD)
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}
	return classStmt{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}, nil
}

// fnDecl → "fn" function ;
// function → IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(ft fnType) (functionStmt, error) {
	if ft == fnTypeFUNCTION {
		_, err := p.consume(FN, "Expect 'fn' at the beginning of function declaration.")
		if err != nil {
			return functionStmt{}, err
		}
	}
	name, err := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", ft))
	if err != nil {
		return functionStmt{}, err
	}
	fnLiteral, err := p.functionLiteral(ft)
	if err != nil {
		return functionStmt{}, err
	}
	return functionStmt{name: name, literal: fnLiteral}, nil
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
			return nil, p.er.ParseError(p.peek(), "Expect expression.")
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return varStmt{name: name, initializer: initializer}, nil
}

/*
statement → exprStmt | forStmt | ifStmt | printStmt | returnStmt | whileStmt
| breakStmt | continueStmt | block ;
*/
func (p *Parser) statement() (stmt, error) {
	switch {
	case p.match(FOR):
		return p.forStatement()
	case p.match(IF):
		return p.ifStatement()
	case p.match(PRINT):
		return p.printStatement()
	case p.match(RETURN):
		return p.returnStatement()
	case p.match(WHILE):
		return p.whileStatement()
	case p.match(BREAK):
		return p.breakStatement()
	case p.match(CONTINUE):
		return p.continueStatement()
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

// exprStmt → expression ";" | labeledLoopStmt ;
func (p *Parser) exprStatement() (stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.match(COLON) {
		return p.labeledLoopStatement(expr)
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return exprStmt{expr: expr}, nil
}

// labeledLoopStmt → IDENTIFIER ":" ( forStmt | whileStmt ) ;
func (p *Parser) labeledLoopStatement(labelExpr expr) (stmt, error) {
	tok, _ := p.advance()
	var label variableExpr
	var ok bool
	if label, ok = labelExpr.(variableExpr); !ok {
		return nil, p.er.ParseError(tok, "Invalid label expression.")
	}
	var loop stmt
	var err error
	switch {
	case p.match(FOR):
		loop, err = p.forStatement()
		if err != nil {
			return nil, err
		}
		forLoop, ok := loop.(forStmt)
		if !ok {
			return nil, p.er.ParseError(tok, "Expect for loop statement.")
		}
		forLoop.whileBody.label = label.name
		return forLoop, nil
	case p.match(WHILE):
		loop, err = p.whileStatement()
		if err != nil {
			return nil, err
		}
		whileLoop, ok := loop.(whileStmt)
		if !ok {
			return nil, p.er.ParseError(tok, "Expect while loop statement.")
		}
		whileLoop.label = label.name
		return whileLoop, nil
	default:
		return nil, p.er.ParseError(tok, "Expect loop statement after label.")
	}
}

// forStmt → "for" ( varDecl | exprStmt | ";" ) expression? ";" expression? block ;
func (p *Parser) forStatement() (stmt, error) {
	var err error
	if _, err := p.consume(FOR, "Expect loop."); err != nil {
		return nil, err
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
	var condition expr
	if !p.match(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after condition."); err != nil {
		return nil, err
	}
	var increment expr
	if !p.match(LEFT_BRACE) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	bodyStmts, err := p.block()
	if err != nil {
		return nil, err
	}

	if condition == nil {
		condition = literalExpr{true}
	}
	var incStmt stmt
	if increment == nil {
		incStmt = nil
	} else {
		incStmt = exprStmt{increment}
	}

	var out stmt = forStmt{
		initializer: initializer,
		whileBody: whileStmt{
			condition: condition,
			body:      blockStmt{bodyStmts},
			increment: incStmt,
		},
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

// returnStmt → "return" expression? ";" ;
func (p *Parser) returnStatement() (stmt, error) {
	tok, err := p.consume(RETURN, "Expect return statement.")
	if err != nil {
		return nil, err
	}
	var value expr
	if !p.match(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(SEMICOLON, "Expect ';' after return value."); err != nil {
		return nil, err
	}
	return returnStmt{keyword: tok, value: value}, nil
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

// breakStmt → "break" ";"
func (p *Parser) breakStatement() (stmt, error) {
	tok, err := p.consume(BREAK, "Expect 'break' at the beginning of breakStatement.")
	if err != nil {
		return nil, err
	}
	var label token
	if !p.match(SEMICOLON) {
		label, err = p.consume(IDENTIFIER, "Expect loop label.")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after break.")
	if err != nil {
		return nil, err
	}
	return breakStmt{keyword: tok, label: label}, nil
}

// continueStmt → "continue" ";"
func (p *Parser) continueStatement() (stmt, error) {
	tok, err := p.consume(CONTINUE, "Expect 'continue' at the beginning of continueStatement.")
	if err != nil {
		return nil, err
	}
	var label token
	if !p.match(SEMICOLON) {
		label, err = p.consume(IDENTIFIER, "Expect loop label.")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after continue.")
	if err != nil {
		return nil, err
	}
	return continueStmt{keyword: tok, label: label}, nil
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

// assignment → ( call "." )? IDENTIFIER "=" assignment | logic_or ;
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
		if varExpr, ok := out.(variableExpr); ok {
			out = assignExpr{name: varExpr.name, value: val}
		} else if getExpr, ok := out.(getExpr); ok {
			out = setExpr{object: getExpr.object, name: getExpr.name, value: val}
		} else {
			return nil, p.er.ParseError(tok, "Invalid assignment target.")
		}
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
		thenExpr, err := p.ternary()
		if err != nil {
			return nil, err
		}
		rOper, err := p.advance()
		if err != nil {
			return nil, err
		}
		if !rOper.hasType(COLON) {
			return nil, p.er.ParseError(lOper, "Expect ':' after expression.")
		}
		elseExpr, err := p.ternary()
		if err != nil {
			return nil, err
		}
		out = ternaryExpr{
			condition: out,
			thenExpr:  thenExpr,
			elseExpr:  elseExpr,
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

// unary → ( "!" | "-" ) unary | call ;
func (p *Parser) unary() (expr, error) {
	if p.match(BANG, MINUS) {
		oper, _ := p.advance()
		next, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unaryExpr{operator: oper, right: next}, nil
	}
	return p.call()
}

/*
primary → "true" | "false" | "nil" | "this"
| NUMBER | STRING | IDENTIFIER | "(" expression ")"
| "super" "." IDENTIFIER | "fn" "(" parameters? ")" block ;
*/
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
	case tok.hasType(THIS):
		return thisExpr{tok}, nil
	case tok.hasType(NUMBER, STRING):
		return literalExpr{tok.literal}, nil
	case tok.hasType(IDENTIFIER):
		return variableExpr{tok}, nil
	case tok.hasType(LEFT_PAREN):
		out, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return groupingExpr{out}, nil
	case tok.hasType(SUPER):
		_, err := p.consume(DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(IDENTIFIER, "Expect superclass method name.")
		if err != nil {
			return nil, err
		}
		return superExpr{keyword: tok, method: method}, nil
	case tok.hasType(FN):
		return p.functionLiteral(fnTypeANONYMOUS)
	case tok.hasType(LEFT_BRACKET):
		return p.arrayLiteral()
	case tok.hasType(SLASH, STAR, MINUS, PLUS, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, BANG, BANG_EQUAL):
		_, err := p.expression()
		if err != nil {
			return nil, err
		}
		return nil, p.er.ParseError(tok, "Expect left operand.")
	default:
		return nil, p.er.ParseError(tok, "Expect an expression.")
	}
}

// call → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
// arguments → expression ( "," expression )* ;
func (p *Parser) call() (expr, error) {
	out, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(LEFT_PAREN) {
			out, err = p.finishCall(out)
			if err != nil {
				return nil, err
			}
		} else if p.match(LEFT_BRACKET) {
			out, err = p.index(out)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			p.advance()
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			out = getExpr{object: out, name: name}
		} else {
			break
		}
	}
	return out, nil
}

func (p *Parser) finishCall(callee expr) (expr, error) {
	tok, err := p.consume(LEFT_PAREN, "Expect '(' at call.")
	if err != nil {
		return nil, err
	}
	args := make([]expr, 0)
	if !p.match(RIGHT_PAREN) {
		for {
			if len(args) >= 255 {
				return nil, p.er.ParseError(p.peek(), "Can't have more than 255 arguments.")
			}
			// TODO: only allowing 'assignment' expression or higher in function call
			// because the comma operator is not allowed in function call (can be confused with
			// parameter seperator comma). To be disallowed with resolver.
			arg, err := p.assignment()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			if p.match(COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return callExpr{callee: callee, paren: tok, arguments: args}, nil
}

// functionLiteral → "fn" "(" parameters? ")" block ;
func (p *Parser) functionLiteral(ft fnType) (functionExpr, error) {
	var errMsg string
	switch ft {
	case fnTypeFUNCTION, fnTypeMETHOD:
		errMsg = fmt.Sprintf("Expect '(' after %s name.", ft)
	case fnTypeANONYMOUS:
		errMsg = "Expect '(' in anonymous function."
	}
	_, err := p.consume(LEFT_PAREN, errMsg)
	if err != nil {
		return functionExpr{}, err
	}
	parameters := make([]token, 0)
	if !p.match(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return functionExpr{}, p.er.ParseError(p.peek(), "Can't have more than 255 parameters.")
			}
			param, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return functionExpr{}, err
			}
			parameters = append(parameters, param)
			if p.match(COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return functionExpr{}, err
	}
	bodyStmts, err := p.block()
	if err != nil {
		return functionExpr{}, err
	}
	return functionExpr{params: parameters, body: bodyStmts}, nil
}

// arrayLiteral → "[" arrayItems "]" ;
// arrayItems → expression ( "," expression )* ;
func (p *Parser) arrayLiteral() (arrayExpr, error) {
	array := make([]expr, 0)
	for !p.match(RIGHT_BRACKET) {
		if len(array) >= 255 {
			return arrayExpr{}, p.er.ParseError(p.peek(), "Can't have more than 255 items in array literal.")
		}
		item, err := p.assignment()
		if err != nil {
			return arrayExpr{}, err
		}
		array = append(array, item)
		if p.match(COMMA) {
			p.advance()
		} else {
			break
		}
	}
	_, err := p.consume(RIGHT_BRACKET, "Expect ']' after array items.")
	if err != nil {
		return arrayExpr{}, err
	}
	return arrayExpr{value: array}, nil
}

func (p *Parser) index(callee expr) (expr, error) {
	tok, err := p.consume(LEFT_BRACKET, "Expect '[' at indexing.")
	if err != nil {
		return nil, err
	}
	index, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RIGHT_BRACKET, "Expect ']' after index expression.")
	if err != nil {
		return nil, err
	}
	return indexExpr{callee: callee, bracket: tok, index: index}, nil
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
// Returns the consumed token if matched, or a ParseError otherwise.
func (p *Parser) consume(expected tokenType, errMsg string) (token, error) {
	if p.peek().tokenType != expected {
		return token{}, p.er.ParseError(p.peek(), errMsg)
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
func (p *Parser) match(expected ...tokenType) bool {
	return slices.Contains(expected, p.peek().tokenType)
}

// isAtEnd returns whether there is more token to parse
func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

// peek returns the current token without consuming it. Returns the last token
// if there is no more token to peek at
func (p *Parser) peek() token {
	if p.current >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current]
}
