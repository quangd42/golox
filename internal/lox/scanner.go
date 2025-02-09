package lox

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrEOF                = errors.New("EOF")
	ErrUnterminatedString = errors.New("unterminated string")
	ErrInvalidNumber      = errors.New("invalid number")
)

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
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)

	// One or two character tokens.
	case '!':
		if s.match('=') {
			s.addToken2(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken2(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken2(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken2(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}

	case '/':
		if s.match('/') {
			// This is a comment, ignore every character until end of line '\n'
			for {
				s.advance()
				c, err := s.peek()
				if c == '\n' || errors.Is(err, ErrEOF) {
					break
				}
			}
		} else {
			s.addToken(SLASH, nil)
		}

	// Ignore some white space
	case ' ', '\t', '\r':
	// New line
	case '\n':
		s.line++

		// Literals
	case '"':
		return s.addTokenString()
	default:
		switch {
		case isDigit(char):
			return s.addTokenNumber()
		}
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

func (s *Scanner) addToken(t TokenType, literals any) {
	s.Tokens = append(s.Tokens, Token{
		Type:     t,
		Lexeme:   string(s.Source[s.lexemeStart:s.cursor]),
		Literals: literals,
		Line:     s.line,
	})
}

// addToken2 calls addToken and consumes 1 addtional character
func (s *Scanner) addToken2(t TokenType, literals any) {
	s.cursor++
	s.addToken(t, literals)
}

// peek returns the next rune without consuming it
func (s Scanner) peek() (rune, error) {
	if s.cursor >= len(s.Source) {
		return 0, ErrEOF
	}
	return rune(s.Source[s.cursor]), nil
}

// match peeks at the next rune and returns whether it matches expected
func (s Scanner) match(expected rune) bool {
	if s.cursor >= len(s.Source) {
		return false
	}
	return s.Source[s.cursor] == byte(expected)
}

func (s *Scanner) addTokenString() error {
	for {
		c, err := s.peek()
		if errors.Is(err, ErrEOF) {
			return ErrUnterminatedString
		}
		if c == '\n' {
			s.line++
		}
		if c == '"' {
			// consume closing '"'
			s.advance()
			// trim surrounding quotes
			s.addToken(STRING, s.Source[s.lexemeStart+1:s.cursor])
			break
		}
		s.advance()
	}
	return nil
}

func (s *Scanner) addTokenNumber() error {
	var isFloat bool
	for {
		c, err := s.peek()
		if errors.Is(err, ErrEOF) {
			return ErrUnterminatedString
		}
		if !isDigit(c) && c != '.' {
			break
		}
		if c == '.' {
			if !isFloat {
				isFloat = true
			} else {
				// TODO: report error with line number
				return ErrInvalidNumber
			}
		}
		s.advance()
	}
	numStr := string(s.Source[s.lexemeStart:s.cursor])
	var num any
	var err error
	if isFloat {
		num, err = strconv.ParseFloat(numStr, 32)
	} else {
		num, err = strconv.Atoi(numStr)
	}
	if err != nil {
		return ErrInvalidNumber
	}
	s.addToken(NUMBER, num)
	return nil
}
