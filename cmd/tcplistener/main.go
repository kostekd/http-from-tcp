package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
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
		fmt.Println("Headers:")
		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		connection.Close()
	}
}