package headers

import (
	"fmt"
	"regexp"
	s "strings"
)

var EMPTY_HEADER = ""

func validateHeaderSyntax(header string) (string, error){
	regex := regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+:\s*.+$`)
	
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
	/*
	   Only fully loaded headers should be parsed and validated.
	   For example, given the input:
	   "Host: localhost:42069\r\nU"
	   the partial header "U" should not be parsed, as it is incomplete and not a valid header.
	*/
	lastCRLF := s.LastIndex(str, CRLF)

	headers := s.Split(str[:lastCRLF], CRLF)
	bytesParsed := 0

	for _, header := range headers {
		// it means that the header we are trying to parse is the empty line
		if header == EMPTY_HEADER {
			return bytesParsed, true, nil
		}
		trimedHeader := s.TrimSpace(header)
		 _, err = validateHeaderSyntax(trimedHeader)

		if err != nil {
			return bytesParsed, false, fmt.Errorf("error: invalid header syntax")
		}

	 	keyValue := s.Split(trimedHeader, ":")

		key := s.ToLower(keyValue[0])
		value := s.Join(keyValue[1:], ":")[1:]
		curr, ok := h[key]
		if ok {
			h[key] = s.Join([]string{curr, value}, ",")
		} else {
			h[key] = value
	
		}
		bytesParsed += len(header) + len(CRLF)
	}

	return bytesParsed, false, nil
}

func (h Headers) Get(key string) string {
	return h[s.ToLower(key)]
}