package request

import (
	"httpfromtcp/internal/buffer"
	"httpfromtcp/internal/headers"
	httpErr "httpfromtcp/internal/httperrors"
	"io"
	"regexp"
	s "strings"
)

const BUFFER_SIZE = 8
const NOTHING_PARSED = 0
const CRLF = "\r\n"

type ParsingState string

const (
	Initialized ParsingState = "initialized"
	RequestLineParsed ParsingState = "requestLineParsed"
	HeadersParsed ParsingState = "headersParsed"
	Done ParsingState = "Done"
)

func httpMethodParser(str string) (string, error) {
	regex := regexp.MustCompile(`^[A-Z]+$`)

	if(!regex.MatchString(str)) { 
		return "", httpErr.ExceptionMessages[httpErr.InvalidHttpMethod]()
	}

	return str, nil
}

func httpVersionParser(str string) (string, error) {
    regex := regexp.MustCompile(`^HTTP/\d+(\.\d+)*$`)

	if(!regex.MatchString(str)) {
		return "", httpErr.ExceptionMessages[httpErr.InvalidHttpMethod]()
	}

	return s.Split(str, "/")[1], nil;
}

func httpRequestTargetParser(str string) (string, error) {
	regex := regexp.MustCompile(`^/([A-Za-z0-9._~!$&'()*+,;=:@%-]*(/[A-Za-z0-9._~!$&'()*+,;=:@%-]*)*)(\?[A-Za-z0-9._~!$&'()*+,;=:@%/?-]*)?(#[A-Za-z0-9._~!$&'()*+,;=:@%/?-]*)?$`)

	if(!regex.MatchString(str)) {
		return "", httpErr.ExceptionMessages[httpErr.InvalidRequestTarget]()
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
		return 0, httpErr.ExceptionMessages[httpErr.InvalidRequestLineFormat](len(requestLine))
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
	State ParsingState
}

func (r *Request) done() bool {
	return r.State == Done
}

func (r *Request) parse(data *buffer.Buf) (int, error) { 
	switch r.State {
	case Initialized:
		n, err := parseRequestLine(data.B, r)
		if n > 0 {
			r.State = RequestLineParsed
		}
		return n, err
	case RequestLineParsed: 
		n, done, err := r.Headers.Parse(data.B)
		if done {
			contentLength := r.Headers.GetContentLength()
			//skip last step because there is no body to be parsed
			if contentLength == 0 {
				r.State = Done
			} else {
				r.State = HeadersParsed
			}
			return n, err
		}
		return n, err
	case HeadersParsed:
		contentLength := r.Headers.GetContentLength()
		r.Body = append(r.Body, data.B[:data.R]...)
		if len(r.Body) == contentLength {
			r.State = Done
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
		State: Initialized,
		Headers: headers.Headers{},
	}

	buf := buffer.New(BUFFER_SIZE)
	
	for !request.done() {
		chunk, errReader := reader.Read(buf.B[buf.R:])
		buf.R += chunk
		if errReader != nil && errReader != io.EOF {
			return nil, errReader
		}
		n, err := request.parse(buf)

		if errReader == io.EOF && request.Headers.GetContentLength() > len(request.Body) {
			return nil, httpErr.ExceptionMessages[httpErr.BodyTooShort]()
		}

		if err != nil {
			return nil, err
		}
		
		if !(n == NOTHING_PARSED) {
			buf.R -= n
			buf.Free(n)
		}
		
		if(buf.R >= len(buf.B)) {
			buf.Grow()
		}
	}
	return request, nil
}