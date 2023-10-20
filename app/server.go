package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	path := strings.Split(string(buf), "\r\n")[1]
	fmt.Println(path)
	if path == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	}

	conn.Close()
}
