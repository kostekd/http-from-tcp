package request

import (
	"fmt"
	"io"
	"regexp"
	s "strings"
)

type State int

const BUFFER_SIZE = 8

type Request struct {
	RequestLine RequestLine
	// State values:
	// 0 - initialized
	// 1 - done
	State 		int
}

//TODO: To implement
func (r *Request) parse(data []byte) (int, error) { 
	return 0, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

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

	if !s.Contains(str, "\r\n") {
		return 0, nil
	}

	requestLine := s.Split(str, " ")

	if len(requestLine) != 3 {
		return 0, fmt.Errorf(`invalid request line format. Requires 3 properties. Received %d`, len(requestLine))	
	}

	method, err := httpMethodParser(requestLine[0])
	if err != nil {
		return 0, err
	}

	requestTarget, err := httpRequestTargetParser(requestLine[1])
	if err != nil {
		return 0, err
	}

	version, err := httpVersionParser(requestLine[2])
	if err != nil {
		return 0, err
	}
	request.RequestLine = RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}

	return len(str), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    //TODO: Do not load the whole input at once
	// input, err := io.ReadAll(reader)

	request := &Request{
		State: 0,
	}
	readToIndex := 0

	for {
		buf := make([]byte, BUFFER_SIZE, BUFFER_SIZE)
		chunk, err := reader.Read(buf)

		if err != nil {
			return nil, err
		}
		strs := s.Split(string(buf), "\r\n")
	}


	if len(strs) < 1 {
		return nil, fmt.Errorf("invalid format: Too few lines")
	}
	n, err := parseRequestLine(strs[0], request)
	fmt.Print(n)
	if err != nil {
		return nil, err
	}
	return request, nil
}