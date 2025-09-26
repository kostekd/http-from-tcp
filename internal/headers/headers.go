package headers

import (
	"fmt"
	"regexp"
	s "strings"
)

func validateHeaderSyntax(header string) (string, error){
	regex := regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+:\s*.+$`)
	
	if(!regex.MatchString(header)) { 
		return "", fmt.Errorf("invalid http method")
	}
	
	return header, nil
}

type Headers map[string]string

const CRLF = "\r\n"

//TODO: For now I will leave it with a default done as false all the time but tbh I dont understand why.
func (h Headers) Parse(data []byte) (n int, done bool, err error) {	
	str := string(data)
	if !s.Contains(str, CRLF) {	
		return 0, false, nil
	}
	headers := s.Split(str, CRLF)
	bytesParsed := 0

	for _, header := range headers {
		// it means that the header we are trying to parse is the empty line
		if header == "" {
			return bytesParsed, true, nil
		}
		trimedHeader := s.TrimSpace(header)
		 _, err = validateHeaderSyntax(trimedHeader)

		if err != nil {
			return bytesParsed, false, fmt.Errorf("error: invalid syntax")
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