package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"net"
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
		go s.handle(connection)
	}
}

func (s *Server) handle(conn net.Conn) {
	request, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Print("Error: request parsing failed\n")
	}

	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for k, v := range request.Headers {
		fmt.Printf("- %s: %s\n", k, v)
	}
	fmt.Printf("Body:\n%s\n", string(request.Body))
	

	//for now no body so content-length is 0
	headers := response.GetDefaultHeaders(0)

	response.WriteStatusLine(conn, response.HttpStatusOK)
	response.WriteHeaders(conn, headers)

	conn.Close()
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":" + fmt.Sprintf("%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to start TCP listener on :%d", port)
	}
	server := Server{
		listener: l,
	}
	go server.listen()

	return &server, nil
}