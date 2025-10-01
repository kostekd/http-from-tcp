package httperrors

import "fmt"

type Exception string 

const (
	InvalidHttpMethod    Exception = "invalid_http_method"
	InvalidHttpVersion   Exception = "invalid_http_version"
	InvalidRequestTarget Exception = "invalid_request_target"
	InvalidRequestLineFormat Exception = "invalid_request_line_format"
	BodyTooShort         Exception = "body_too_short"
)

var ExceptionMessages = map[Exception]func(args... any) error{
	InvalidHttpMethod:    func(args... any) error { return fmt.Errorf("error: invalid HTTP method") },
	InvalidHttpVersion:   func(args... any) error { return fmt.Errorf("error: invalid HTTP version") },
	InvalidRequestTarget: func(args... any) error { return fmt.Errorf("error: Invalid requestTarget") },
	BodyTooShort:         func(args... any) error { return fmt.Errorf("error: body too short") },
	InvalidRequestLineFormat: func(args... any) error {return fmt.Errorf(`invalid request line format. Requires 3 properties. Received %d`, args[0])},
}
