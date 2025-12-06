package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Server Crashed! %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func(rwc io.ReadWriteCloser) {
			defer rwc.Close()
			handle(rwc, rwc)
		}(conn)
	}
}

func handle(w io.Writer, r io.Reader) {
	buf := make([]byte, 512)

	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		w.Write([]byte("-ERR internal error\r\n"))
		return
	}

	data := buf[:n]
	fmt.Printf("DEBUG: %q\n", data)

	tokens, err := parse(bytes.NewBuffer(data))
	if err != nil {
		w.Write([]byte("-ERR internal error\r\n"))
		return
	}

	if len(tokens) == 1 && strings.ToUpper(tokens[0]) == "PING" {
		w.Write([]byte("+PONG\r\n"))
		return
	}

	if len(tokens) == 2 && strings.ToUpper(tokens[0]) == "ECHO" {
		w.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(tokens[1]), tokens[1])))
		return
	}

	w.Write([]byte("-ERR unknown command\r\n"))
}

func parse(r io.Reader) ([]string, error) {
	s := NewScanner(r)

	result := []string{}
	for {
		tok, err := s.Scan()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(tok, "*") {
			// TODO: consume array length
			continue
		}

		if strings.HasPrefix(tok, "$") {
			// TODO: consume bulk string length

			// The next token is the actual string
			tok, err = s.Scan()
			if err != nil {
				return nil, err
			}

			result = append(result, tok)

			continue
		}

		return nil, fmt.Errorf("unexpected token: %q", tok)
	}

	return result, nil
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (string, error) {
	tok, err := s.r.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(tok, "\r\n"), nil
}
