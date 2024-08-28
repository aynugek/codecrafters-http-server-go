package main

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	CRLF           = "\r\n"
	StatusOK       = "200 OK"
	StatusNotFound = "404 Not Found"
)

type Request struct {
	Method      string
	Target      string
	HTTPVersion string
	Headers     map[string]string
	Body        string
}

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
	reqStr := string(reqBytes)

	var req Request
	req.Headers = make(map[string]string)

	parts := strings.Split(reqStr, CRLF)
	requestLine := strings.Fields(parts[0])
	req.Method = requestLine[0]
	req.Target = requestLine[1]
	req.HTTPVersion = requestLine[2]
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		key, value, found := strings.Cut(part, ":")
		if found {
			req.Headers[strings.ToLower(key)] = strings.TrimSpace(strings.ToLower(value))
			continue
		}
		req.Body = strings.TrimSpace(strings.Join(parts[i:], ""))
		break
	}

	var b bytes.Buffer

	if req.Target == "/" {
		conn.Write([]byte(req.HTTPVersion + " " + StatusOK + CRLF + CRLF))
	} else if req.Target == "/user-agent" {
		if ua, exists := req.Headers["user-agent"]; exists {
			b.WriteString(req.HTTPVersion + " " + StatusOK + CRLF)
			b.WriteString("Content-Type: text/plain" + CRLF)
			b.WriteString("Content-Length: " + strconv.Itoa(len(ua)) + CRLF)
			b.WriteString(CRLF)
			b.WriteString(ua)
			conn.Write(b.Bytes())
		}
	} else if strings.HasPrefix(req.Target, "/echo/") {
		endpoint := strings.TrimPrefix(req.Target, "/echo/")
		b.WriteString(req.HTTPVersion + " " + StatusOK + CRLF)
		b.WriteString("Content-Type: text/plain" + CRLF)
		b.WriteString("Content-Length: " + strconv.Itoa(len(endpoint)) + CRLF)
		b.WriteString(CRLF)
		b.WriteString(endpoint)
		conn.Write(b.Bytes())
	} else {
		conn.Write([]byte(req.HTTPVersion + " " + StatusNotFound + CRLF + CRLF))
	}
}
