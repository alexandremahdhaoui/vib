/*
Copyright 2023 Alexandre Mahdhaoui

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
