package util

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code      int    `json:"-"`
	Message   any    `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	Developer string `json:"developer_message,omitempty"`
	Internal  error  `json:"-"`
}

func NewHTTPError(code int, err error, message ...any) *HTTPError {
	e := &HTTPError{Code: code, Message: http.StatusText(code), Internal: err}
	if len(message) > 0 {
		e.Message = message[0]
	}
	return e
}

func (e *HTTPError) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", e.Code, e.Message, e.Internal)
}

func (e *HTTPError) Unwrap() error {
	return e.Internal
}
