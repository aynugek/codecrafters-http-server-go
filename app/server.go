package main

import (
	"log"
	"net"
)

const (
	CRLF     = "\r\n"
	StatusOK = "HTTP/1.1 200 OK" + CRLF
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatal("Failed to bind to port 4221")
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatal("Error accepting connection: ", err.Error())
	}
	defer conn.Close()

	conn.Write([]byte(StatusOK + CRLF))
}
