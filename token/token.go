package token

type Typ string

const (
	EOF       Typ = "EOF"
	DELIMITER Typ = "DELIMITER" // \r\n
	DOLLAR    Typ = "DOLLAR"    // $
	ASTERISK  Typ = "ASTERISK"  // *
	STRING    Typ = "STRING"    // PING
)

type Token struct {
	Type    Typ
	Literal string
}

func EndOfFile() Token {
	return Token{Type: EOF, Literal: ""}
}

func Dollar() Token {
	return Token{Type: DOLLAR, Literal: "$"}
}

func Asterisk() Token {
	return Token{Type: ASTERISK, Literal: "*"}
}

func Delimiter() Token {
	return Token{Type: DELIMITER, Literal: "\r\n"}
}

func String(literal string) Token {
	return Token{Type: STRING, Literal: literal}
}
