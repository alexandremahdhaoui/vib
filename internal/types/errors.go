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

package types

import (
	"errors"
	"fmt"
)

var (
	// ErrType is returned when an unsupported type is encountered.
	ErrType = errors.New("ERRTYPE: unsupported type")
	// ErrFile is returned when an unsupported file extension is encountered.
	ErrFile = errors.New("ERRFILE: unsupported file extension")
	// ErrEncoding is returned when an unsupported encoding is encountered.
	ErrEncoding = errors.New("ERRENC: unsupported encoding")
	// ErrExists is returned when a resource already exists.
	ErrExists = errors.New("ERREXISTS: resource already exist")
	// ErrNotFound is returned when a resource cannot be found.
	ErrNotFound = errors.New("ERRNOTFOUND: resource cannot be found")
	// ErrArgs is returned when unexpected arguments are provided.
	ErrArgs = errors.New("ERRARGS: unexpected argument")
	// ErrVal is returned when input cannot be validated.
	ErrVal = errors.New("ERRVAL: input cannot be validated")
	// ErrRef is returned when an unexpected reference is encountered.
	ErrRef = errors.New("ERRREF: unexpected reference")
)

// ErrAtIndex returns an error with the given index.
func ErrAtIndex(i int) error {
	return fmt.Errorf("error is propably located at index %d", i)
}
