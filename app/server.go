package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func debug(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

type Request struct {
	Method  string
	Path    string
	conn    net.Conn
	Headers Headers
}

func NewRequest(conn net.Conn) *Request {
	req := &Request{conn: conn, Headers: make(Headers)}
	req.parse(conn)
	return req
}

func (r *Request) parse(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	firstLine := strings.Split(string(buf), "\r\n")[0]

	r.Path = strings.Split(firstLine, " ")[1]
	r.Method = strings.Split(firstLine, " ")[0]

	for _, line := range strings.Split(string(buf), "\r\n")[1:] {
		debug("Parse line: ", line)
		key := strings.Split(line, ": ")[0]
		if key == "" {
			break
		}
		r.Headers[strings.Split(line, ": ")[0]] = strings.Split(line, ": ")[1]
	}
}

type Headers map[string]string

func (h Headers) String() string {
	var str string
	for k, v := range h {
		str += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return str
}

type Response struct {
	conn       net.Conn
	Headers    Headers
	StatusCode int
	Body       []byte
}

var codeNames = map[int]string{
	200: "OK",
	404: "Not Found",
}

func NewResponse(conn net.Conn) *Response {
	return &Response{conn: conn, Headers: make(map[string]string)}
}

func (r *Response) WriteHeader(k, val string) *Response {
	r.Headers[k] = val
	return r
}

func (r *Response) WriteStatusCode(s int) *Response {
	r.StatusCode = s
	return r
}

func (r *Response) WriteBody(b []byte) *Response {
	r.Body = b

	r.Headers["Content-Length"] = fmt.Sprintf("%d", len(b))
	return r
}

func (r *Response) Send() {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s \r\n", r.StatusCode, codeNames[r.StatusCode])
	headers := r.Headers.String()
	bodyLine := "\r\n" + string(r.Body)
	debug("Sending: ", statusLine+headers+bodyLine+"\r\n")
	r.conn.Write([]byte(statusLine + headers + bodyLine + "\r\n"))

	r.conn.Close()
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func() {

			req := NewRequest(conn)
			res := NewResponse(conn)

			debug("Path is: ", req.Path)

			if req.Path == "/" {
				res.WriteHeader("Content-type", "text/plain").WriteStatusCode(200).Send()
				return
			}

			if strings.HasPrefix(req.Path, "/echo") {
				res.WriteHeader("Content-type", "text/plain").WriteStatusCode(200).WriteBody([]byte(strings.Split(req.Path, "/echo/")[1])).Send()
				return
			}

			if strings.HasPrefix(req.Path, "/user-agent") {
				res.WriteHeader("Content-type", "text/plain").WriteStatusCode(200).WriteBody([]byte(req.Headers["User-Agent"])).Send()
				return
			}

			res.WriteHeader("Content-type", "text/plain").WriteStatusCode(404).Send()
		}()
	}
}
