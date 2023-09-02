package vib

import (
	"errors"
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

var (
	ErrType          = errors.New("ERRTYPE")
	ErrFileExtension = errors.New("ERRFILEEXTENSION: unsupported file extension")
	ErrEncoding      = errors.New("ERRENCODING: unsupported encoding")
	ErrAlreadyExist  = errors.New("ERRALREADYEXIST")
	ErrNotFound      = errors.New("ERRNOTFOUND")
	ErrArgs          = errors.New("ERRARGS")
	ErrValidation    = errors.New("ERRVALIDATION")
)

func NewErr(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

// NewErrAndLog takes an error and a message, concatenate it, log it and return the err
func NewErrAndLog(err error, msg string) error {
	err = NewErr(err, msg)
	logger.SugaredLoggerWithCallerSkip(2).Errorf(err.Error()) //nolint:gomnd

	return err
}

func errKind(kind Kind) error {
	return fmt.Errorf("%w: kind %q is not supported", ErrType, kind)
}

func errApiVersion(apiVersion APIVersion, kind Kind) error {
	return fmt.Errorf("%w: APIVersion %q for Kind %q is not supported", ErrType, apiVersion, kind)
}
