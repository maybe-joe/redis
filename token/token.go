package token

type Type string

const (
	EOF       Type = "EOF"
	DELIMITER Type = "DELIMITER" // \r\n
	DOLLAR    Type = "DOLLAR"    // $
	ASTERISK  Type = "ASTERISK"  // *
	STRING    Type = "STRING"    // PING
)

type Token struct {
	Type    Type
	Literal string
}

func EndOfFile() Token {
	return Token{Type: EOF, Literal: ""}
}

func Delimiter() Token {
	return Token{Type: DELIMITER, Literal: "\r\n"}
}

func Dollar() Token {
	return Token{Type: DOLLAR, Literal: "$"}
}

func Asterisk() Token {
	return Token{Type: ASTERISK, Literal: "*"}
}

func String(literal string) Token {
	return Token{Type: STRING, Literal: literal}
}

func (t Token) IsEOF() bool {
	return t.Type == EOF
}

func (t Token) IsDelimiter() bool {
	return t.Type == DELIMITER
}

func (t Token) IsDollar() bool {
	return t.Type == DOLLAR
}

func (t Token) IsAsterisk() bool {
	return t.Type == ASTERISK
}

func (t Token) IsString() bool {
	return t.Type == STRING
}
