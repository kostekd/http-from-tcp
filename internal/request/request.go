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
    regex := regexp.MustCompile(`^HTTP/\d+(\.\d+)*$`)

	if(!regex.MatchString(str)) {
		return "", fmt.Errorf("invalid HTTP version")
	}

	return s.Split(str, "/")[1], nil;
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    input, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Raw input %s",input)
	strs := s.Split(string(input), "\r\n")
	fmt.Printf("Splits %s", strs[0])

	if len(strs) < 1 {
		return nil, fmt.Errorf("invalid format: Too few lines")
	}
	requestLine := s.Split(strs[0], " ")

	if len(requestLine) != 3 {
		return nil, fmt.Errorf(`invalid request line format. Requires 3 properties. Received %d`, len(requestLine))	
	}

	httpMethod, err := httpMethodParser(requestLine[2])
	if err != nil {
		return nil, err
	}

	request := &Request{
		RequestLine: RequestLine{
			Method:        requestLine[0],
			RequestTarget: requestLine[1],
			HttpVersion:   httpMethod,
		},
	}

	return request, nil
}