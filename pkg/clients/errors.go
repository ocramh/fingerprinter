package clients

import (
	"errors"
	"fmt"
)

var (
	ErrStatusNotOK = errors.New("the response status code was not ok")
)

// HTTPError is the Error interface implementation used for HTTP errors
type HTTPError struct {
	code    int
	message string
}

// NewHTTPError returns a new HTTPError instance
func NewHTTPError(statusCode int, message string) HTTPError {
	return HTTPError{
		code:    statusCode,
		message: message,
	}
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("http status: %d. error message: %s", h.code, h.message)
}
