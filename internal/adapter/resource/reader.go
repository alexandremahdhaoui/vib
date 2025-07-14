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

package resourceadapter

import (
	"errors"
	"io"

	"github.com/alexandremahdhaoui/vib/internal/types"
)

// TODO: create an adapter that can Get resource from CLI arguments, file or stdin.
// That adapter must read resources from "anonymous files" by first unmarshalling a raw resource,
// and then getting the appropriate storage from an APIServer.

// NewReader instantiates a new resource adapter that can operate on any
// form of byte stream input.
func NewReader[T types.APIVersionKind](
	dynDecoder types.DynamicDecoder[T],
	reader io.Reader,
) types.Storage[T] {
	// reader bytes.Reader
	return &byteReader[T]{
		dynDecoder: dynDecoder,
		reader:     reader,
	}
}

// Only implements List and Get.
// Does not implement Create, Update or Delete.
type byteReader[T types.APIVersionKind] struct {
	dynDecoder types.DynamicDecoder[T]
	reader     io.Reader
}

// List many resources from the input byte stream.
func (r *byteReader[T]) List() ([]types.Resource[T], error) {
	return r.dynDecoder.Decode(r.reader)
}

// Get one resource from the input byte stream.
func (r *byteReader[T]) Get(_ string) (types.Resource[T], error) {
	out, err := r.dynDecoder.Decode(r.reader)
	if err != nil {
		return types.Resource[T]{}, err
	}

	if len(out) < 1 {
		return types.Resource[T]{}, errors.New("resource cannot be found in input")
	}

	return out[0], nil
}

func (r *byteReader[T]) Create(types.Resource[T]) error {
	panic("unimplemented")
}

func (r *byteReader[T]) Delete(name string) error {
	panic("unimplemented")
}

func (r *byteReader[T]) Update(oldName string, v types.Resource[T]) error {
	panic("unimplemented")
}
