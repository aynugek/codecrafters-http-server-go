package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CRLF           = "\r\n"
	StatusOK       = "200 OK"
	StatusNotFound = "404 Not Found"
	StatusCreated  = "201 Created"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting connection: ", err.Error())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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

	var sb strings.Builder
	sb.WriteString(req.HTTPVersion + " " + StatusOK + CRLF)

	if req.Target == "/" {
		conn.Write([]byte(req.HTTPVersion + " " + StatusOK + CRLF + CRLF))
	} else if req.Target == "/user-agent" {
		if ua, exists := req.Headers["user-agent"]; exists {
			sb.WriteString("Content-Type: text/plain" + CRLF)
			sb.WriteString("Content-Length: " + strconv.Itoa(len(ua)) + CRLF)
			sb.WriteString(CRLF)
			sb.WriteString(ua)
			conn.Write([]byte(sb.String()))
		}
	} else if strings.HasPrefix(req.Target, "/echo/") {
		endpoint := strings.TrimPrefix(req.Target, "/echo/")
		sb.WriteString("Content-Type: text/plain" + CRLF)
		sb.WriteString("Content-Length: " + strconv.Itoa(len(endpoint)) + CRLF)
		if enc, exists := req.Headers["accept-encoding"]; exists {
			if enc == "gzip" {
				sb.WriteString("Content-Encoding: gzip" + CRLF)
			}
		}
		sb.WriteString(CRLF)
		sb.WriteString(endpoint)
		conn.Write([]byte(sb.String()))
	} else if strings.HasPrefix(req.Target, "/files/") {
		dir := os.Args[2]
		filename := filepath.Join(dir, strings.TrimPrefix(req.Target, "/files/"))

		if req.Method == "GET" {
			data, err := os.ReadFile(filename)
			if err != nil {
				goto NotFound
			}
			sb.WriteString("Content-Type: application/octet-stream" + CRLF)
			sb.WriteString("Content-Length: " + strconv.Itoa(len(data)) + CRLF)
			sb.WriteString(CRLF)
			sb.Write(data)
		} else if req.Method == "POST" {
			err = os.WriteFile(filename, []byte(req.Body), 0644)
			if err != nil {
				log.Fatal(err.Error())
			}
			sb.Reset()
			sb.WriteString(req.HTTPVersion + " " + StatusCreated + CRLF + CRLF)
		}

		conn.Write([]byte(sb.String()))
	}
NotFound:
	conn.Write([]byte(req.HTTPVersion + " " + StatusNotFound + CRLF + CRLF))
}
