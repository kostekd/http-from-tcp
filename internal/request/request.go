package request

import (
	"fmt"
	"io"
	"regexp"
	s "strings"
)

type Request struct {
	RequestLine RequestLine
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

func parseRequestLine(str string) (*RequestLine, error) {
	requestLine := s.Split(str, " ")

	if len(requestLine) != 3 {
		return nil, fmt.Errorf(`invalid request line format. Requires 3 properties. Received %d`, len(requestLine))	
	}

	method, err := httpMethodParser(requestLine[0])
	if err != nil {
		return nil, err
	}

	requestTarget, err := httpRequestTargetParser(requestLine[1])
	if err != nil {
		return nil, err
	}

	version, err := httpVersionParser(requestLine[2])
	if err != nil {
		return nil, err
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    input, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	strs := s.Split(string(input), "\r\n")

	if len(strs) < 1 {
		return nil, fmt.Errorf("invalid format: Too few lines")
	}
	requestLine, err := parseRequestLine(strs[0])

	if err != nil {
		return nil, err
	}

	return &Request{*requestLine}, nil
}