package lox

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrEOF = errors.New("EOF")

type Scanner struct {
	er      ErrorReporter
	tokens  []token
	source  []byte
	line    int
	current int
	start   int
}

func NewScanner(er ErrorReporter, source []byte) *Scanner {
	return &Scanner{
		er:     er,
		tokens: make([]token, 0),
		source: source,
		line:   1,
	}
}

func (s *Scanner) ScanTokens() ([]token, error) {
	var err error
	for !s.isAtEnd() {
		err = s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.tokens = append(s.tokens, newToken(EOF, "", nil, s.line, s.start))
	return s.tokens, nil
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
		s.addToken(LEFT_PAREN, "(")
	case ')':
		s.addToken(RIGHT_PAREN, ")")
	case '{':
		s.addToken(LEFT_BRACE, "{")
	case '}':
		s.addToken(RIGHT_BRACE, "}")
	case ',':
		s.addToken(COMMA, ",")
	case ':':
		s.addToken(COLON, ":")
	case '.':
		s.addToken(DOT, ".")
	case '-':
		s.addToken(MINUS, "-")
	case '+':
		s.addToken(PLUS, "+")
	case '?':
		s.addToken(QUESTION, "?")
	case ';':
		s.addToken(SEMICOLON, ";")
	case '*':
		s.addToken(STAR, "*")

	// One or two character tokens.
	case '!':
		if s.matchConsume('=') {
			s.addToken(BANG_EQUAL, "!=")
		} else {
			s.addToken(BANG, "!")
		}
	case '=':
		if s.matchConsume('=') {
			s.addToken(EQUAL_EQUAL, "==")
		} else {
			s.addToken(EQUAL, "=")
		}
	case '>':
		if s.matchConsume('=') {
			s.addToken(GREATER_EQUAL, ">=")
		} else {
			s.addToken(GREATER, ">")
		}
	case '<':
		if s.matchConsume('=') {
			s.addToken(LESS_EQUAL, "<=")
		} else {
			s.addToken(LESS, "<")
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
			s.addToken(SLASH, "/")
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
			s.er.ScanError(s.line, fmt.Sprintf("unsupported character '%s'", string(char)))
			return nil
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
	s.tokens = append(s.tokens, newToken(t, s.makeLexeme(), literal, s.line, s.start))
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
			s.er.ScanError(s.line, "unterminated string")
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
			isFloat = true
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
		s.er.ScanError(s.line, fmt.Sprintf("invalid number '%s'", lex))
		return nil
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
		s.addToken(IDENTIFIER, lex)
	} else {
		s.addToken(tt, tt.String())
	}
	return nil
}
