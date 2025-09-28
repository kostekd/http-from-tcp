package server

import (
	"fmt"
	"net"
	"httpfromtcp/internal/request"
)
type Server struct {
	listener net.Listener
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		connection, err := s.listener.Accept()
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
		fmt.Printf("Body:\n%s\n", string(request.Body))
		connection.Close()
	}
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":" + fmt.Sprintf("%d", port))
	if err != nil {
		fmt.Printf("Failed to start TCP listener on :%d\n", port)
	}
	return &Server{
		listener: l,
	}, nil
}