package lox

import "fmt"

type Scanner struct {
	Source      []byte
	Tokens      []Token
	line        int
	cursor      int
	lexemeStart int
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{
		Source: source,
	}
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.Tokens = append(s.Tokens, Token{
		Type:   EOF,
		Lexeme: "",
		Line:   s.line,
	})
	return s.Tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.cursor >= len(s.Source)
}

func (s *Scanner) scanToken() error {
	char := s.advance()
	switch char {
	// Single-character tokens.
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)

	// One or two character tokens.
	case '!':
		if s.match('=') {
			s.addToken2(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken2(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '>':
		if s.match('=') {
			s.addToken2(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '<':
		if s.match('=') {
			s.addToken2(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}

	default:
		return fmt.Errorf("unexpected character: %s", string(char))
	}
	s.lexemeStart = s.cursor
	return nil
}

// advance **consumes** a character and returns it
func (s *Scanner) advance() rune {
	out := s.Source[s.cursor]
	s.cursor++
	return rune(out)
}

func (s *Scanner) addToken(t TokenType) {
	s.Tokens = append(s.Tokens, Token{
		Type:   t,
		Lexeme: string(s.Source[s.lexemeStart:s.cursor]),
		Line:   s.line,
	})
}

// addToken2 calls addToken and consumes 1 addtional character
func (s *Scanner) addToken2(t TokenType) {
	s.cursor++
	s.addToken(t)
}

// match peeks at the next rune and returns whether it matches expected
func (s Scanner) match(expected rune) bool {
	if s.cursor >= len(s.Source) {
		return false
	}
	return s.Source[s.cursor] == byte(expected)
}
