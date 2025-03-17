package lox

var while = blockStmt{
	statements: []stmt{
		varStmt{
			name:        token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1},
			initializer: literalExpr{value: 1},
		},
		whileStmt{
			condition: binaryExpr{
				left:     variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 2}},
				operator: token{tokenType: 20, lexeme: "<", literal: interface{}(nil), line: 2},
				right:    literalExpr{value: 1000},
			},
			body: blockStmt{
				statements: []stmt{
					printStmt{expr: variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 3}}},
					exprStmt{expr: assignExpr{
						name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 4},
						value: binaryExpr{
							left:     variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 4}},
							operator: token{tokenType: 13, lexeme: "*", literal: interface{}(nil), line: 4},
							right:    literalExpr{value: 2},
						},
					}},
				},
			},
		},
	},
}

var divide = blockStmt{
	statements: []stmt{
		varStmt{
			name:        token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1},
			initializer: variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1}},
		},
		whileStmt{
			condition: binaryExpr{
				left:     variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1}},
				operator: token{tokenType: 20, lexeme: "<", literal: interface{}(nil), line: 1},
				right:    literalExpr{value: 1000},
			},
			body: blockStmt{
				statements: []stmt{
					printStmt{expr: variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 2}}},
					exprStmt{expr: assignExpr{
						name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1},
						value: binaryExpr{
							left:     variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1}},
							operator: token{tokenType: 13, lexeme: "*", literal: interface{}(nil), line: 1},
							right:    literalExpr{value: 2},
						},
					}},
				},
			},
		},
	},
}

// WARN: Locals missing i from initializer
var locals_for = map[expr]int{
	variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 0}}: 1,
	variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 1}}: 1,
	assignExpr{
		name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 0},
		value: binaryExpr{
			left:     variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 0}},
			operator: token{tokenType: 13, lexeme: "*", literal: interface{}(nil), line: 0},
			right:    literalExpr{value: 2},
		},
	}: 1,
}

var locals_while = map[expr]int{
	variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 2}}: 0,
	variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 3}}: 1,
	variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 4}}: 1,
	assignExpr{
		name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 4},
		value: binaryExpr{
			left:     variableExpr{name: token{tokenType: 22, lexeme: "i", literal: interface{}(nil), line: 4}},
			operator: token{tokenType: 13, lexeme: "*", literal: interface{}(nil), line: 4},
			right:    literalExpr{value: 2},
		},
	}: 1,
}
