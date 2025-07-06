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

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

var (
	ErrType     = errors.New("ERRTYPE: unsupported type")
	ErrFile     = errors.New("ERRFILE: unsupported file extension")
	ErrEncoding = errors.New("ERRENC: unsupported encoding")
	ErrExist    = errors.New("ERREXIST: resource already exist")
	ErrNotFound = errors.New("ERRNOTFOUND: resource cannot be found")
	ErrArgs     = errors.New("ERRARGS: unexpected argument")
	ErrVal      = errors.New("ERRVAL: input cannot be validated")
	ErrRef      = errors.New("ERRREF: unexpected reference")
)

func NewKindErr(kind Kind) error {
	return flaterrors.Join(ErrType, fmt.Errorf("kind %q is not supported", kind))
}

func NewAPIVersionErr(apiVersion APIVersion, kind Kind) error {
	return flaterrors.Join(
		ErrType,
		fmt.Errorf("APIVersion %q for Kind %q is not supported", apiVersion, kind),
	)
}

func NewRefErr(reference string, kind Kind) error {
	return flaterrors.Join(
		ErrRef,
		fmt.Errorf("cannot resolve reference %q to Kind %q", reference, kind),
	)
}
