package main

import (
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatal("Failed to bind to port 4221")
	}
	defer l.Close()

	_, err = l.Accept()
	if err != nil {
		log.Fatal("Error accepting connection: ", err.Error())
	}
}
