package response

import (
	"fmt"
	h "httpfromtcp/internal/headers"
	httpErr "httpfromtcp/internal/httperrors"
	"io"
)


type StatusCode int

const (
	HttpStatusOK                  StatusCode = 200
	HttpStatusBadRequest          StatusCode = 400
	HttpStatusInternalServerError StatusCode = 500
)


func WriteStatusLine (w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case HttpStatusOK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return nil
	case HttpStatusBadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return nil
	case HttpStatusInternalServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return nil
	default:
		return httpErr.ExceptionMessages[httpErr.UnrecognizedStatusCode]()
	}
}

func WriteHeaders(w io.Writer, headers h.Headers) error {
	for key, value := range headers {
		header :=  fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(header))

		if err != nil {
			return err
		}
	}
	//write an empty header
	_, err := w.Write([]byte("\r\n"))
	
	if err != nil {
		return err
	}

	return nil
}

func WriteBody(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) h.Headers {
	m := make(h.Headers)

	m["content-length"] = fmt.Sprintf("%d", contentLen)
	m["connection"] = "close"
	m["content-type"] = "text/plain"

	return m
}