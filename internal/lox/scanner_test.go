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
			want:  []token{newToken(LEFT_PAREN, "(", nil, 0)},
		},
		{
			desc:  "Two_Char__BANG_EQUAL",
			input: []byte("!="),
			want:  []token{newToken(BANG_EQUAL, "!=", nil, 0)},
		},
		{
			desc:  "Comment",
			input: []byte("// This is some comment text"),
			want:  nil, // Comment is ignored
		},
		{
			desc:  "New_line_String",
			input: []byte("\n"),
			want:  nil,
		},
		{
			desc:  "String",
			input: []byte("\"This is some string\""),
			want:  []token{newToken(STRING, "\"This is some string\"", "This is some string", 0)},
		},
		{
			desc:  "Number_FLOAT",
			input: []byte("17.8"),
			want:  []token{newToken(NUMBER, "17.8", 17.8, 0)},
		},
		{
			desc:  "Number_INT",
			input: []byte("178"),
			want:  []token{newToken(NUMBER, "178", 178, 0)},
		},
		{
			desc:  "Identifier_Keyword",
			input: []byte("var"),
			want:  []token{newToken(VAR, "var", nil, 0)},
		},
		{
			desc:  "Identifier__User_Defined",
			input: []byte("golox"),
			want:  []token{newToken(IDENTIFIER, "golox", nil, 0)},
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
				newToken(LEFT_PAREN, "(", nil, 0), newToken(RIGHT_PAREN, ")", nil, 0),
				newToken(LESS, "<", nil, 0), newToken(GREATER, ">", nil, 0), newToken(GREATER_EQUAL, ">=", nil, 0), newToken(LESS_EQUAL, "<=", nil, 3),
				newToken(BANG_EQUAL, "!=", nil, 3), newToken(BANG, "!", nil, 3), newToken(EQUAL_EQUAL, "==", nil, 3),
				newToken(STRING, "\"This is some string\"", "This is some string", 4),
				newToken(NUMBER, "23.", 23., 5), newToken(STRING, "\"some more string\"", "some more string", 5), newToken(NUMBER, "156", 156, 5),
				newToken(THIS, "this", nil, 6), newToken(AND, "and", nil, 6), newToken(IDENTIFIER, "that", nil, 6),
				newToken(OR, "or", nil, 6), newToken(IDENTIFIER, "never", nil, 6), newToken(FALSE, "false", nil, 6), newToken(TRUE, "true", nil, 6),
				newToken(RETURN, "return", nil, 7), newToken(IDENTIFIER, "out", nil, 7), newToken(EOF, "", nil, 8),
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
