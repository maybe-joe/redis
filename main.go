package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/maybe-joe/redis/command"
	"github.com/maybe-joe/redis/lexer"
	"github.com/maybe-joe/redis/parser"
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

	data := string(buf[:n])
	fmt.Printf("DEBUG: %q\n", data)

	var (
		l = lexer.New(data)
		p = parser.New(l)
	)

	cmd, err := p.Parse()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)

		switch err {
		case parser.ErrUnknownCommand:
			w.Write([]byte("-ERR unknown command\r\n"))
		default:
			w.Write([]byte("-ERR parse error\r\n"))
		}

		return
	}

	switch cmd.Type {
	default:
		// This should never happen
		w.Write([]byte("-ERR internal error\r\n"))
	case command.UNKNOWN:
		w.Write([]byte("-ERR unknown command\r\n"))
	case command.PING:
		w.Write([]byte("+PONG\r\n"))
	case command.ECHO:
		msg := cmd.Args[0]
		if len(msg) > 0 {
			fmt.Fprintf(w, "$%d\r\n%s\r\n", len(msg), msg)
		} else {
			w.Write([]byte("$0\r\n\r\n"))
		}
	}
}
