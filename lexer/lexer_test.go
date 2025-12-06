package lexer

import (
	"testing"

	"github.com/maybe-joe/redis/token"
	"github.com/stretchr/testify/assert"
)

func Test_Lexer_Next(t *testing.T) {
	testcases := []struct {
		name     string
		given    string
		expected token.Token
	}{
		{name: "eof", given: "", expected: token.EndOfFile()},
		{name: "delimiter", given: "\r\n", expected: token.Delimiter()},
		{name: "dollar", given: "$", expected: token.Dollar()},
		{name: "asterisk", given: "*", expected: token.Asterisk()},
		{name: "string", given: "PING", expected: token.String("PING")},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, New(tc.given).Next())
		})
	}
}

func Test_Lexer_Lex(t *testing.T) {
	testcases := []struct {
		name     string
		given    string
		expected []token.Token
	}{

		{
			name:  "PING",
			given: "*1\r\n$4\r\nPING\r\n",
			expected: []token.Token{
				token.Asterisk(),
				token.String("1"),
				token.Delimiter(),
				token.Dollar(),
				token.String("4"),
				token.Delimiter(),
				token.String("PING"),
				token.Delimiter(),
				token.EndOfFile(),
			},
		},
		{
			name:  "ECHO hey",
			given: "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n",
			expected: []token.Token{
				token.Asterisk(),
				token.String("2"),
				token.Delimiter(),
				token.Dollar(),
				token.String("4"),
				token.Delimiter(),
				token.String("ECHO"),
				token.Delimiter(),
				token.Dollar(),
				token.String("3"),
				token.Delimiter(),
				token.String("hey"),
				token.Delimiter(),
				token.EndOfFile(),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, New(tc.given).Lex())
		})
	}
}
