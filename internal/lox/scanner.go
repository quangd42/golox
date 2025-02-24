package lox

import (
	"errors"
	"strconv"
)

var ErrEOF = errors.New("EOF")

type Scanner struct {
	Tokens  []token
	source  []byte
	line    int
	current int
	start   int
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{source: source}
}

func (s *Scanner) ScanTokens() ([]token, error) {
	for !s.isAtEnd() {
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.Tokens = append(s.Tokens, newToken(EOF, "", nil, s.line))
	return s.Tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	defer s.startNextLexeme()
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
		if s.matchConsume('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.matchConsume('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '>':
		if s.matchConsume('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '<':
		if s.matchConsume('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}

	case '/':
		if s.matchConsume('/') {
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

	// String
	case '"':
		return s.addTokenString()

	default:
		switch {
		case isDigit(char):
			return s.addTokenNumber()
		case isAlpha(char):
			return s.addTokenIdentifier()
		default:
			return errUnsupportedCharacter(char, s.line)
		}
	}
	return nil
}

func (s *Scanner) startNextLexeme() {
	s.start = s.current
}

func (s Scanner) makeLexeme() string {
	return string(s.source[s.start:s.current])
}

// advance **consumes** a character and returns it
func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current-1])
}

// addToken appends a new token to the scanner's internal tokens
func (s *Scanner) addToken(t tokenType, literal any) {
	s.Tokens = append(s.Tokens, newToken(t, s.makeLexeme(), literal, s.line))
}

// peek returns the current rune without consuming it
func (s Scanner) peek() (rune, error) {
	if s.isAtEnd() {
		return 0, ErrEOF
	}
	return rune(s.source[s.current]), nil
}

// matchConsume peeks at the current rune, if the current rune matches expected it is consumed.
// returns whether expected rune was matched and consumed.
func (s *Scanner) matchConsume(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] == byte(expected) {
		s.advance()
		return true
	}
	return false
}

func (s *Scanner) addTokenString() error {
	for {
		c, err := s.peek()
		if errors.Is(err, ErrEOF) {
			return errUnterminatedString(s.line)
		}
		if c == '\n' {
			s.line++
		}
		if c == '"' {
			// consume closing '"'
			s.advance()
			// trim surrounding quotes
			s.addToken(STRING, string(s.source[s.start+1:s.current-1]))
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
			break
		}
		if !isDigit(c) && c != '.' {
			break
		}
		if c == '.' {
			if !isFloat {
				isFloat = true
			} else {
				return errInvalidNumber(s.line)
			}
		}
		s.advance()
	}
	lex := s.makeLexeme()
	var num any
	var err error
	if isFloat {
		num, err = strconv.ParseFloat(lex, 64)
	} else {
		num, err = strconv.Atoi(lex)
	}
	if err != nil {
		return errInvalidNumber(s.line)
	}
	s.addToken(NUMBER, num)
	return nil
}

func (s *Scanner) addTokenIdentifier() error {
	for {
		c, err := s.peek()
		if errors.Is(err, ErrEOF) {
			break
		}
		if !isAlphaNum(c) {
			break
		}
		s.advance()
	}
	lex := s.makeLexeme()
	tt, err := getKeywords(lex)
	if err != nil {
		s.addToken(IDENTIFIER, nil)
	} else {
		s.addToken(tt, nil)
	}
	return nil
}
