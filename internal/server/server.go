package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message []byte
}

func (hErr *HandlerError) Write(w io.Writer) {
	headers := response.GetDefaultHeaders(len(hErr.Message))
	response.WriteStatusLine(w, hErr.StatusCode)
	response.WriteHeaders(w, headers)
	response.WriteBody(w, hErr.Message)
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	listener net.Listener
	handler Handler
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
    defer conn.Close()

	request, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			Message: []byte(err.Error()),
			StatusCode: response.HttpStatusInternalServerError,
		}
		hErr.Write(conn)
		return
	}

	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for k, v := range request.Headers {
		fmt.Printf("- %s: %s\n", k, v)
	}
	fmt.Printf("Body:\n%s\n", string(request.Body))
	
	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, request)

	if hErr != nil {
		hErr.Write(conn)
		//close the connection after sending what went wrong
		return
	}

	body := buf.Bytes()
	//for now no body so content-length is 0
	headers := response.GetDefaultHeaders(len(body))

	response.WriteStatusLine(conn, response.HttpStatusOK)
	response.WriteHeaders(conn, headers)
	response.WriteBody(conn, body)
}

func Serve(port int, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":" + fmt.Sprintf("%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{
		listener: l,
		handler: h,
	}
	go server.listen()

	return &server, nil
}
