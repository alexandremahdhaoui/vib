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
	"bytes"

	"github.com/alexandremahdhaoui/vib/internal/types"
)

// TODO: create an adapter that can Get resource from CLI arguments, file or stdin.
// That adapter must read resources from "anonymous files" by first unmarshalling a raw resource,
// and then getting the appropriate storage from an APIServer.

// NewByteStream instantiates a new resource adapter that can operate on any
// form of input.
func NewByteStream[T types.APIVersionKind](
	reader bytes.Reader,
) types.Storage[T] {
	// reader bytes.Reader
	return &byteStream[T]{}
}

// Only implements List and Get.
// Does not imlement Create, Update or Delete.
type byteStream[T types.APIVersionKind] struct {
	reader bytes.Reader
}

// List many resources from the input byte stream.
func (b *byteStream[T]) List() ([]types.Resource[T], error) {
	// TODO:
	panic("unimplemented")
}

// Get one resource from the input byte stream.
func (b *byteStream[T]) Get(_ string) (types.Resource[T], error) {
	// TODO:
	panic("unimplemented")
}

func (b *byteStream[T]) Create(types.Resource[T]) error {
	panic("unimplemented")
}

func (b *byteStream[T]) Delete(name string) error {
	panic("unimplemented")
}

func (b *byteStream[T]) Update(oldName string, v types.Resource[T]) error {
	panic("unimplemented")
}
