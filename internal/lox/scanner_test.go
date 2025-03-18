package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var text = `() <> >=
// This is some comment
// And some more
<= != ! ==
"This is some string"
23. "some more string" 156
this and that or never false true
return out
`

func Test_scanToken(t *testing.T) {
	testCases := []struct {
		desc  string
		input []byte
		want  []token
	}{
		{
			desc:  "One_Char__LEFT_PAREN",
			input: []byte("("),
			want:  []token{newToken(LEFT_PAREN, "(", "(", 1, 0)},
		},
		{
			desc:  "Two_Char__BANG_EQUAL",
			input: []byte("!="),
			want:  []token{newToken(BANG_EQUAL, "!=", "!=", 1, 0)},
		},
		{
			desc:  "Comment",
			input: []byte("// This is some comment text"),
			want:  []token{}, // Comment is ignored
		},
		{
			desc:  "New_line_String",
			input: []byte("\n"),
			want:  []token{},
		},
		{
			desc:  "String",
			input: []byte("\"This is some string\""),
			want:  []token{newToken(STRING, "\"This is some string\"", "This is some string", 1, 0)},
		},
		{
			desc:  "Number_FLOAT",
			input: []byte("17.8"),
			want:  []token{newToken(NUMBER, "17.8", 17.8, 1, 0)},
		},
		{
			desc:  "Number_INT",
			input: []byte("178"),
			want:  []token{newToken(NUMBER, "178", 178, 1, 0)},
		},
		{
			desc:  "Identifier_Keyword",
			input: []byte("var"),
			want:  []token{newToken(VAR, "var", "var", 1, 0)},
		},
		{
			desc:  "Identifier__User_Defined",
			input: []byte("golox"),
			want:  []token{newToken(IDENTIFIER, "golox", "golox", 1, 0)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner(tC.input)
			err := scanner.scanToken()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, scanner.Tokens)
		})
	}
}

func TestScanTokens(t *testing.T) {
	testCases := []struct {
		desc  string
		input []byte
		want  []token
	}{
		{
			desc:  "sample text",
			input: []byte(text),
			want: []token{
				newToken(LEFT_PAREN, "(", "(", 1, 0),
				newToken(RIGHT_PAREN, ")", ")", 1, 1),
				newToken(LESS, "<", "<", 1, 3),
				newToken(GREATER, ">", ">", 1, 4),
				newToken(GREATER_EQUAL, ">=", ">=", 1, 6),

				newToken(LESS_EQUAL, "<=", "<=", 4, 50),
				newToken(BANG_EQUAL, "!=", "!=", 4, 53),
				newToken(BANG, "!", "!", 4, 56),
				newToken(EQUAL_EQUAL, "==", "==", 4, 58),

				newToken(STRING, "\"This is some string\"", "This is some string", 5, 61),

				newToken(NUMBER, "23.", 23., 6, 83),
				newToken(STRING, "\"some more string\"", "some more string", 6, 87),
				newToken(NUMBER, "156", 156, 6, 106),

				newToken(THIS, "this", "this", 7, 110),
				newToken(AND, "and", "and", 7, 115),
				newToken(IDENTIFIER, "that", "that", 7, 119),
				newToken(OR, "or", "or", 7, 124),
				newToken(IDENTIFIER, "never", "never", 7, 127),
				newToken(FALSE, "false", "false", 7, 133),
				newToken(TRUE, "true", "true", 7, 139),

				newToken(RETURN, "return", "return", 8, 144),
				newToken(IDENTIFIER, "out", "out", 8, 151),

				newToken(EOF, "", nil, 9, 155),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			scanner := NewScanner(tC.input)
			got, err := scanner.ScanTokens()
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tC.want, got)
		})
	}
}
