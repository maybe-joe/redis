package lexer

import (
	"github.com/maybe-joe/redis/token"
)

// Lexer converts RESP input into tokens for parsing
// See: https://redis.io/docs/latest/develop/reference/protocol-spec
type Lexer struct {
	// data, the RESP input to be lexed
	data string
	// current, the position of the lexer in the input data
	current int
	// next, the next position of the lexer in the input data
	next int
	// ch, the character at the index of current in data
	ch byte
}

func New(data string) *Lexer {
	l := &Lexer{data: data}
	l.Advance()
	return l
}

// Advance advances the lexer by one character
func (l *Lexer) Advance() {
	if l.next >= len(l.data) {
		// If we've reached the end of the input data, reset ch to 0
		l.ch = 0
	} else {
		// Otherwise, advance current and next, and set ch to the next character
		l.ch = l.data[l.next]
	}

	// Question: What is the point of advancing if we've already reached the end of the data?
	l.current = l.next
	l.next++
}

// Peek returns the next character without advancing the lexer
func (l *Lexer) Peek() byte {
	if l.next >= len(l.data) {
		return 0
	} else {
		return l.data[l.next]
	}
}

// Next returns the next token from the input data
func (l *Lexer) Next() token.Token {
	var t token.Token

	switch l.ch {
	case 0:
		t = token.EndOfFile()
	case '$':
		t = token.Dollar()
	case '*':
		t = token.Asterisk()
	case '\r':
		if l.Peek() == '\n' {
			l.Advance() // consume '\n'
			t = token.Delimiter()
		}
	default:
		start := l.current
		for l.ch != 0 && (l.ch != '\r' && l.Peek() != '\n') {
			l.Advance()
		}
		t = token.String(l.data[start:l.current])
		return t
	}

	l.Advance()

	return t
}

// Lex converts the entire input data into a slice of tokens
// Intended for testing full RESP commands are lexed correctly
func (l *Lexer) Lex() []token.Token {
	var tokens []token.Token

	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}
