package parser

import (
	"testing"

	"github.com/maybe-joe/redis/command"
	"github.com/maybe-joe/redis/lexer"
	"github.com/stretchr/testify/assert"
)

func Test_Parse_Parse(t *testing.T) {
	testcases := []struct {
		name     string
		given    Lexer
		expected command.Command
	}{

		{
			name:     "PING",
			given:    lexer.New("*1\r\n$4\r\nPING\r\n"),
			expected: command.Ping(),
		},
		{
			name:     "ECHO hey",
			given:    lexer.New("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"),
			expected: command.Echo("hey"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := New(tc.given).Parse()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
