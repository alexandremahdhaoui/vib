package logger

import (
	"errors"
	"fmt"
)

var (
	ErrType          = errors.New("ERRTYPE")
	ErrFileExtension = errors.New("ERRFILEEXTENSION: unsupported file extension")
	ErrEncoding      = errors.New("ERRENCODING: unsupported encoding")
	ErrAlreadyExist  = errors.New("ERRALREADYEXIST")
	ErrNotFound      = errors.New("ERRNOTFOUND")
	ErrArgs          = errors.New("ERRARGS")
	ErrValidation    = errors.New("ERRVALIDATION")
	ErrReference     = errors.New("ERRREFERENCE")
)

func NewErr(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

// NewErrAndLog takes an error and a message, concatenate it, log it and return the err
func NewErrAndLog(err error, msg string) error {
	err = NewErr(err, msg)
	SugaredLoggerWithCallerSkip(1).Errorf(err.Error()) //nolint:gomnd

	return err
}
