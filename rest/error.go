package rest

import (
	"fmt"
)

func NewHttpError(url, msg string, statusCode int) *httpError {
	return &httpError{
		Url:        url,
		StatusCode: statusCode,
		Message:    msg,
	}
}

// httpError is a wrapper over errors returned from external services
type httpError struct {
	Url        string
	StatusCode int
	Message    string
}

func (h httpError) Error() string {
	return fmt.Sprintf("an HTTP error response has been received from: %s {statusCode=%d, response=%s}", h.Url, h.StatusCode, h.Message)
}
