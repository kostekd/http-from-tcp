package request

import (
	"fmt"
	"httpfromtcp/internal/buffer"
	"httpfromtcp/internal/headers"
	"io"
	"regexp"
	s "strings"
)

type State int

const BUFFER_SIZE = 8
const NOTHING_PARSED = 0
const CRLF = "\r\n"

func httpMethodParser(str string) (string, error) {
	regex := regexp.MustCompile(`^[A-Z]+$`)

	if(!regex.MatchString(str)) { 
		return "", fmt.Errorf("invalid http method")
	}

	return str, nil
}

func httpVersionParser(str string) (string, error) {
    regex := regexp.MustCompile(`^HTTP/\d+(\.\d+)*$`)

	if(!regex.MatchString(str)) {
		return "", fmt.Errorf("invalid HTTP version")
	}

	return s.Split(str, "/")[1], nil;
}

func httpRequestTargetParser(str string) (string, error) {
	regex := regexp.MustCompile(`^/([A-Za-z0-9._~!$&'()*+,;=:@%-]*(/[A-Za-z0-9._~!$&'()*+,;=:@%-]*)*)(\?[A-Za-z0-9._~!$&'()*+,;=:@%/?-]*)?(#[A-Za-z0-9._~!$&'()*+,;=:@%/?-]*)?$`)

	if(!regex.MatchString(str)) {
		return "", fmt.Errorf("invalid request target")
	}

	return str, nil;
}

func parseRequestLine(data []byte, request *Request) (int, error) {
	str := string(data)
	if !s.Contains(str, CRLF) {
		return 0, nil
	}

	requestLine := s.Split(str, CRLF)[0]
	parts := s.Split(requestLine, " ")

	if len(parts) != 3 {
		return 0, fmt.Errorf(`invalid request line format. Requires 3 properties. Received %d`, len(requestLine))	
	}

	method, err := httpMethodParser(parts[0])
	if err != nil {
		return 0, err
	}

	requestTarget, err := httpRequestTargetParser(parts[1])
	if err != nil {
		return 0, err
	}

	version, err := httpVersionParser(parts[2])
	if err != nil {
		return 0, err
	}
	request.RequestLine = RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}

	return len(requestLine) + len(CRLF), nil
}

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	Body []byte
	// State values:
	// 0 - initialized
	// 1 - request line parsed
	// 2 - headers parsed
	// 3 - done (body parsed)
	State 		int
}

func (r *Request) parse(data *buffer.Buf) (int, error) { 
	switch r.State {
	case 0:
		n, err := parseRequestLine(data.B, r)
		if n > 0 {
			r.State = 1
		}
		return n, err
	case 1: 
		n, done, err := r.Headers.Parse(data.B)
		if done {
			contentLength := r.Headers.GetContentLength()
			//skip last step because there is no body to be parsed
			if contentLength == 0 {
				r.State = 3
			} else {
				r.State = 2
			}
			return n, err
		}
		return n, err
	case 2:
		contentLength := r.Headers.GetContentLength()
		r.Body = append(r.Body, data.B[:data.R]...)
		if len(r.Body) == contentLength {
			r.State = 3
		}
		return data.R, nil

	}

	return -1, nil
}
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		State: 0,
		Headers: headers.Headers{},
	}

	buf := buffer.New(BUFFER_SIZE)
	
	for request.State != 3 {
		chunk, err := reader.Read(buf.B[buf.R:])
		buf.R += chunk
		if err != nil && err != io.EOF {
			return nil, err
		}
		n, err := request.parse(buf)

		if err != nil {
			return nil, err
		}
		
		if !(n == NOTHING_PARSED) {
			buf.R -= n
			buf.Free(n)
		}
		
		//grow buffer
		if(buf.R >= len(buf.B)) {
			buf.Grow()
		}
	}
	return request, nil
}