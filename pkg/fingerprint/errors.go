package fingerprint

import (
	"errors"
)

var (
	ErrInvalidPath      = errors.New("file does not exists")
	ErrInvalidFileInput = errors.New("invalid input file")
	ErrInvalidFormat    = errors.New("invalid file format")
)
