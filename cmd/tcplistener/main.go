package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)


func main() {
	l, err := net.Listen("tcp", ":42069")
	x	
	if err != nil {
		fmt.Print("Failed to start TCP listener on :42069\n")
	}
	defer l.Close();
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Print("Error: failed to accept connection\n")
		}
		request, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Print("Error: request parsing failed\n")
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		connection.Close()
	}
}