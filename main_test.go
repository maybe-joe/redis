package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Handle(t *testing.T) {
	testcases := []struct {
		name     string
		given    string
		expected string
	}{
		{name: "PING", given: "*1\r\n$4\r\nPING\r\n", expected: "+PONG\r\n"},
		{name: "ECHO", given: "*1\r\n$4\r\nECHO\r\n", expected: "$0\r\n\r\n"},
		{name: "ECHO hey", given: "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n", expected: "$3\r\nhey\r\n"},
		{name: "ECHO Hello World!", given: "*2\r\n$4\r\nECHO\r\n$12\r\nHello World!\r\n", expected: "$12\r\nHello World!\r\n"},
		{name: "Unknown command", given: "*1\r\n$7\r\nUNKNOWN\r\n", expected: "-ERR unknown command\r\n"},
		{name: "Empty command", given: "*0\r\n", expected: "-ERR empty command\r\n"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var actual bytes.Buffer
			handle(&actual, bytes.NewBufferString(tc.given))
			assert.Equal(t, tc.expected, actual.String())
		})
	}
}
