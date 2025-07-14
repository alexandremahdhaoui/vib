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

package readeradapter

import (
	"errors"
	"io"

	"github.com/alexandremahdhaoui/vib/internal/types"
)

// New instantiates a new resource adapter that can operate on any
// form of byte stream input.
func New(
	dynDecoder types.DynamicDecoder[types.APIVersionKind],
	reader io.Reader,
) types.Reader[types.APIVersionKind] {
	// reader bytes.Reader
	return &byteReader{
		dynDecoder: dynDecoder,
		reader:     reader,
	}
}

// Only implements List and Get.
// Does not implement Create, Update or Delete.
type byteReader struct {
	dynDecoder types.DynamicDecoder[types.APIVersionKind]
	reader     io.Reader
}

// List many resources from the input byte stream.
func (r *byteReader) List() ([]types.Resource[types.APIVersionKind], error) {
	return r.dynDecoder.Decode(r.reader)
}

// Get one resource by name from the input.
func (r *byteReader) Get(name string) (types.Resource[types.APIVersionKind], error) {
	resources, err := r.dynDecoder.Decode(r.reader)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	for _, r := range resources {
		if r.Metadata.Name == name {
			return r, nil
		}
	}

	return types.Resource[types.APIVersionKind]{}, errors.New("resource cannot be found")
}
