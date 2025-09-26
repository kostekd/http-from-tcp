package request

import (
	"fmt"
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

func parseRequestLine(str string, request *Request) (int, error) {
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
	// State values:
	// 0 - initialized
	// 1 - request line parsed
	// 2 - done
	State 		int
}

func (r *Request) parse(data []byte) (int, error) { 
	switch r.State {
	case 0:
		n, err := parseRequestLine(string(data), r)
		if n > 0 {
			r.State = 1
		}
		return n, err
	case 1: 
		n, done, err := r.Headers.Parse(data)
		if done {
			r.State = 2
		}
		return n, err
	}

	return 0, nil
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

	readToStart := 0
	bytesParsed := 0
	buf := make([]byte, BUFFER_SIZE)
	
	for request.State != 2 {
		chunk, err := reader.Read(buf[readToStart:])
		readToStart += chunk
		fmt.Printf("BEFORE OPERATIONS -------- buf stats -- len: %d, cap: %d, readToStart: %d, bytesParsed: %d, as string: %q, as bytes: %v\n", len(buf), cap(buf), readToStart, bytesParsed, string(buf[:readToStart]), buf)
		
		if err != nil {
			return nil, err
		}
		
		n, err := request.parse(buf)
		if err != nil {
			return nil, err
		}
		
		if !(n == NOTHING_PARSED) {
			bytesParsed += n
			readToStart -= n
			buf = shiftBuffer(buf, n)
		}
		
		//grow buffer
		if(readToStart >= len(buf)) {
			buf = growBuffer(buf)
		}

		fmt.Printf("AFTER OPERATIONS -------- buf stats -- len: %d, cap: %d, readToStart: %d, bytesParsed: %d, as string: %q, as bytes: %v\n", len(buf), cap(buf), readToStart, bytesParsed, string(buf[:readToStart]), buf)
	}
	return request, nil
}

func growBuffer(buf []byte) []byte {
	biggerBuf := make([]byte, len(buf) * 2)
	copy(biggerBuf, buf)
	buf = biggerBuf
	return buf
}
func shiftBuffer(buf []byte, offset int) []byte {
	shiftBuf := make([]byte, len(buf))
	copy(shiftBuf, buf[offset:])
	buf = shiftBuf
	return buf
}