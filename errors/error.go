package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type Type int

type customError struct {
	errType       Type
	originalError error
}

func (c customError) Error() string {
	return c.originalError.Error()
}

func (t Type) New(msg string) error {
	return customError{
		errType:       t,
		originalError: errors.New(msg),
	}
}

func (t Type) Warp(err error, msg string) error {
	return customError{
		errType:       t,
		originalError: errors.Wrap(err, msg),
	}
}

func (t Type) Warpf(err error, msg string, args ...interface{}) error {
	return customError{
		errType:       t,
		originalError: errors.Wrapf(err, msg, args...),
	}
}

func New(msg string) error {
	return customError{errType: None, originalError: errors.New(msg)}
}

func Newf(msg string, args ...interface{}) error {
	return customError{errType: None, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errType:       customErr.errType,
			originalError: wrappedError,
		}
	}

	return customError{errType: None, originalError: wrappedError}
}

func GetType(err error) Type {
	if customErr, ok := err.(customError); ok {
		return customErr.errType
	}
	return None
}
