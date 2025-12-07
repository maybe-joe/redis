package parser

import (
	"errors"
	"strconv"

	"github.com/maybe-joe/redis/token"

	"github.com/maybe-joe/redis/command"
)

var (
	ErrUnexpectedToken         = errors.New("unexpected token")
	ErrUnknownCommand          = errors.New("unknown command")
	ErrUnableToParseArraySize  = errors.New("unable to parse array size")
	ErrUnableToParseBulkLength = errors.New("unable to parse bulk string length")
	ErrMissingArgument         = errors.New("missing argument")
)

type Lexer interface {
	Next() token.Token
}

type Parser struct {
	lexer Lexer
}

func New(lexer Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() (command.Command, error) {
	next := p.lexer.Next()

	if next.IsAsterisk() {
		// This is an array, the next token should be the size
		next = p.lexer.Next()
		if !next.IsString() {
			return command.Unknown(), ErrUnexpectedToken
		}

		// Attempt to parse the size of the array into an integer
		// so we can range over the elements
		size, err := strconv.Atoi(next.Literal)
		if err != nil {
			return command.Unknown(), ErrUnableToParseArraySize
		}

		// If the size is zero, return an unknown command error
		// This prevents panics later on when we try set the capacity of the args slice
		if size == 0 {
			return command.Unknown(), ErrUnknownCommand
		}

		// The next token should be a delimiter
		next = p.lexer.Next()
		if !next.IsDelimiter() {
			return command.Unknown(), ErrUnexpectedToken
		}

		// The first element should be the command name
		// all other elements are arguments
		var (
			typ  = command.UNKNOWN
			args = make([]string, 0, size-1)
		)

		// NOTE: I'm assuming all elements of the array will be bulk strings for now
		for i := range size {
			// The next token should be a dollar sign
			next = p.lexer.Next()

			if !next.IsDollar() {
				return command.Unknown(), ErrUnexpectedToken
			}

			// The next token should be a string literal representing the length of the bulk string
			next = p.lexer.Next()
			if !next.IsString() {
				return command.Unknown(), ErrUnexpectedToken
			}

			// Attempt to parse the length of the bulk string into an integer
			length, err := strconv.Atoi(next.Literal)
			if err != nil {
				return command.Unknown(), ErrUnableToParseBulkLength
			}

			// TODO: This should be part of the parser logic to actually read the bulk string data
			_ = length

			// The next token should be a delimiter
			next = p.lexer.Next()
			if !next.IsDelimiter() {
				return command.Unknown(), ErrUnexpectedToken
			}

			// The next token should be the actual string data
			next = p.lexer.Next()
			if !next.IsString() {
				return command.Unknown(), ErrUnexpectedToken
			}

			// If this is the first element, it should be the command name
			if i == 0 {
				typ = command.Type(next.Literal)
			} else {
				args = append(args, next.Literal)
			}

			// Finally, the next token should be a delimiter
			next = p.lexer.Next()
			if !next.IsDelimiter() {
				return command.Unknown(), ErrUnexpectedToken
			}
		}

		switch typ {
		case command.PING:
			return command.Ping(), nil
		case command.ECHO:
			msg := ""

			if len(args) > 0 {
				// NOTE: Extra arguments are ignored
				msg = args[0]
			}

			return command.Echo(msg), nil
		default:
			return command.Unknown(), ErrUnknownCommand
		}
	}

	return command.Unknown(), ErrUnexpectedToken
}
