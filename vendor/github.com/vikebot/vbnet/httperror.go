package vbnet

import (
	"fmt"
)

// HTTPError collects useful meta informations for network applications relying
// on the HTTP protocol
type HTTPError interface {
	// Message is a textual description of the problem ready to send to the
	// user
	Message() string
	// HTTPCode defines the code used when this error is returned over the HTTP
	// protocol
	HTTPCode() int
	// Code is a unique internal error code that can be used to exactly
	// aggregate single errors.
	Code() int
	// Inner is the internal error that caused this httpError
	Inner() error

	error
}

type httpErr struct {
	message  string
	httpCode int
	code     int
	inner    error
}

func (err httpErr) Message() string {
	return err.message
}

func (err httpErr) HTTPCode() int {
	return err.httpCode
}

func (err httpErr) Code() int {
	return err.code
}

func (err httpErr) Inner() error {
	return err.inner
}

func (err httpErr) Error() string {
	str := fmt.Sprintf("vbnet.%d: %s (HTTP %d)", err.code, err.message, err.httpCode)
	if err.inner != nil {
		str += fmt.Sprintf(", due-to: %v", err.inner)
	}
	return str
}

// NewHTTPError creates a new instance that implements the HTTPError interface
// and returns it
func NewHTTPError(message string, httpCode int, code int, inner error) HTTPError {
	return httpErr{
		message:  message,
		httpCode: httpCode,
		code:     code,
		inner:    inner,
	}
}
