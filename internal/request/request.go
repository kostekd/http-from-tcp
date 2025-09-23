package request

import (
	"fmt"
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

func parseRequestLine(str string, request *Request) (int, error) {
	fmt.Printf("String to parse %s\n", str)
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

	return len(requestLine), nil
}
type Request struct {
	RequestLine RequestLine
	// State values:
	// 0 - initialized
	// 1 - done
	State 		int
}

func (r *Request) parse(data []byte) (int, error) { 
	n, err := parseRequestLine(string(data), r)
	if n > 0 {
		r.State = 1
	}
	return n, err
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		State: 0,
	}

	bytesRead := 0
	bytesParsed := 0
	buf := make([]byte, BUFFER_SIZE)

	for request.State != 1 {
		chunk, err := reader.Read(buf[bytesRead:])
		if err != nil {
			return nil, err
		}

		n, err := request.parse(buf)
		if err != nil {
			return nil, err
		}

		if n == NOTHING_PARSED {
			bytesRead += chunk
			biggerBuf := make([]byte, len(buf) * 2)
			copy(biggerBuf, buf)
			buf = biggerBuf
		} else {
			bytesParsed += n
		}
	}
	return request, nil
}