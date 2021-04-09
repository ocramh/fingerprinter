package fingerprint

import (
	"errors"
)

var (
	ErrFileNotExists    = errors.New("file does not exists")
	ErrInvalidFileInput = errors.New("input is not a valid file")
	ErrInvalidFormat    = errors.New("file format is not valid")
)
