package headers

import (
	"fmt"
	"regexp"
	s "strings"
)

func validateHeaderSyntax(header string) (string, error){
	regex := regexp.MustCompile(`^[A-Za-z0-9-]+:\s*.+$`)
	
	if(!regex.MatchString(header)) { 
		return "", fmt.Errorf("invalid http method")
	}
	
	return header, nil
}

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	if !s.Contains(str, CRLF) {
		return 0, false, nil
	}
	if s.Index(str, CRLF) == 0 {
		return 0, false, nil
	}

	untrimedHeader := s.Split(str, CRLF)[0]
	header := s.TrimSpace(untrimedHeader)

	_, err = validateHeaderSyntax(header)

	if err != nil {
		return 0, false, fmt.Errorf("invalid header syntax")
	}

	keyValue := s.Split(header, ":")

	key := keyValue[0]
	value := s.Join(keyValue[1:], ":")[1:]
	h[key] = value

	bytesParsed := len(untrimedHeader) + len(CRLF)

	return bytesParsed, false, nil
}