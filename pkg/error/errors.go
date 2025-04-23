package error

import "errors"

var (
	ErrTaskNotFound = errors.New("taks not found")

	ErrTypeConversion = errors.New("type conversion error")
)
