package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type Type int

type CustomError struct {
	errType       Type
	originalError error
}

func (c CustomError) Error() string {
	res := c.originalError.Error()
	if c.errType.String() != "" {
		res = c.errType.String() + ": " + res
	}
	return res
}

func (t Type) New(msg string) error {
	e := CustomError{
		errType:       t,
		originalError: errors.New(msg),
	}
	return e
}

func (t Type) Warp(err error, msg string) error {
	return CustomError{
		errType:       t,
		originalError: errors.Wrap(err, msg),
	}
}

func (t Type) Warpf(err error, msg string, args ...interface{}) error {
	return CustomError{
		errType:       t,
		originalError: errors.Wrapf(err, msg, args...),
	}
}

func New(msg string) error {
	return CustomError{errType: None, originalError: errors.New(msg)}
}

func Newf(msg string, args ...interface{}) error {
	return CustomError{errType: None, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(CustomError); ok {
		return CustomError{
			errType:       customErr.errType,
			originalError: wrappedError,
		}
	}

	return CustomError{errType: None, originalError: wrappedError}
}

func GetType(err error) Type {
	if customErr, ok := err.(CustomError); ok {
		return customErr.errType
	}
	return None
}
