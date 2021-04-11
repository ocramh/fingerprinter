package meta

import (
	"errors"
	"fmt"
)

var (
	ErrStatusNotOK = errors.New("the response status code was not ok")
)

type HTTPError struct {
	message string
	code    int
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("http status: %d. error message: %s", h.code, h.message)
}
