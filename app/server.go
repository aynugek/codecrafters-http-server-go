package main

import (
	"log"
	"net"
	"strings"
)

const (
	CRLF           = "\r\n"
	StatusOK       = "HTTP/1.1 200 OK" + CRLF
	StatusNotFound = "HTTP/1.1 404 Not Found" + CRLF
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

	reqBytes := make([]byte, 1024)
	n, err := conn.Read(reqBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	reqBytes = reqBytes[:n]

	target := strings.Fields(string(reqBytes))[1]
	if target == "/" {
		conn.Write([]byte(StatusOK + CRLF))
		return
	}
	conn.Write([]byte(StatusNotFound + CRLF))
}
