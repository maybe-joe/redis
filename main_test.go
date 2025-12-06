package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Handle(t *testing.T) {
	t.Run("PING", func(t *testing.T) {
		var (
			given    = "*1\r\n$4\r\nPING\r\n"
			expected = "+PONG\r\n"
			actual   bytes.Buffer
		)

		handle(&actual, bytes.NewBufferString(given))
		assert.Equal(t, expected, actual.String())
	})

	t.Run("ECHO hey", func(t *testing.T) {
		var (
			given    = "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
			expected = "$3\r\nhey\r\n"
			actual   bytes.Buffer
		)

		handle(&actual, bytes.NewBufferString(given))
		assert.Equal(t, expected, actual.String())
	})

	t.Run("ECHO \"Hello World!\"", func(t *testing.T) {
		var (
			given    = "*3\r\n$4\r\nECHO\r\n$12\r\nHello World!\r\n"
			expected = "$12\r\nHello World!\r\n"
			actual   bytes.Buffer
		)

		handle(&actual, bytes.NewBufferString(given))
		assert.Equal(t, expected, actual.String())
	})
}

func Test_Parse(t *testing.T) {
	t.Run("PING", func(t *testing.T) {
		var (
			given    = "*1\r\n$4\r\nPING\r\n"
			expected = []string{"PING"}
		)

		actual, err := parse(bytes.NewBufferString(given))
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Echo hey", func(t *testing.T) {
		var (
			given    = "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
			expected = []string{"ECHO", "hey"}
		)

		actual, err := parse(bytes.NewBufferString(given))
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("ECHO Hello World!", func(t *testing.T) {
		var (
			given    = "*2\r\n$4\r\nECHO\r\n$12\r\nHello World!\r\n"
			expected = []string{"ECHO", "Hello World!"}
		)

		actual, err := parse(bytes.NewBufferString(given))
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_Scanner(t *testing.T) {
	s := NewScanner(bytes.NewBufferString("*1\r\n$4\r\nPING\r\n"))

	tok, err := s.Scan()
	assert.NoError(t, err)
	assert.Equal(t, "*1", tok)

	tok, err = s.Scan()
	assert.NoError(t, err)
	assert.Equal(t, "$4", tok)

	tok, err = s.Scan()
	assert.NoError(t, err)
	assert.Equal(t, "PING", tok)

	tok, err = s.Scan()
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, "", tok)
}
