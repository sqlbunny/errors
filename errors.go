package errors

import (
	"errors"
	"fmt"
)

func New(message string) error {
	return withStackInner(errors.New(message))
}

func Errorf(format string, a ...interface{}) error {
	return withStackInner(fmt.Errorf(format, a...))
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
