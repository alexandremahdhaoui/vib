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

package codecadapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/alexandremahdhaoui/vib/internal/types"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"go.yaml.in/yaml/v3"
)

type drd struct {
	apiServer types.APIServer
}

// Instantiate a new detector. A detector is a special codec that can unmarshal one or many documents from
// one or many types.
func NewDynamicResourceDecoder(
	apiServer types.APIServer,
) types.DynamicDecoder[types.APIVersionKind] {
	return drd{
		apiServer: apiServer,
	}
}

var (
	errInputMustBeJsonOrYaml = errors.New("input must be json or yaml")
	errDecodingInput         = errors.New("error decoding input")
	errInputMustNotBeEmpty   = errors.New("input must not be empty")
)

type resourceList[T any] struct {
	Items []types.Resource[T] `json:"items"`
}

// Decode will try to decode the input as json and if it fails as yaml.
// There are no detection mechanisms, the algorithm is really not optimized.
//
// The decoding operation is implemented as such:
// 1. For enc := range [ json yaml ]
// 2. For doc := range reader
// 3. For type := range [ resourceList[Resource[T]] Resource[T] ]
// 4. Decode type until input is exhausted
func (d drd) Decode(reader io.Reader) ([]types.Resource[types.APIVersionKind], error) {
	// not really cool clean but we want to copy the buffer
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, flaterrors.Join(errInputMustNotBeEmpty, errDecodingInput)
	}

	bCopy := make([]byte, len(b))
	copied := copy(bCopy, b)
	if copied < len(b) {
		return nil, flaterrors.Join(errors.New("internal error"), errDecodingInput)
	}

	// -- [ json yaml ]
	jsonBuf := bytes.NewBuffer(b)
	yamlBuf := bytes.NewBuffer(bCopy)

	for _, supportedDecoder := range []struct {
		Decoder       decoder
		UnmarshalFunc func(b []byte, v any) error
	}{
		{
			Decoder:       json.NewDecoder(jsonBuf),
			UnmarshalFunc: json.Unmarshal,
		}, {
			Decoder:       yaml.NewDecoder(yamlBuf),
			UnmarshalFunc: yaml.Unmarshal,
		},
	} {
		rawRes := make([]types.Resource[json.RawMessage], 0)
		// -- resourceList[T]
		rls, rlsErr := decode[resourceList[json.RawMessage]](supportedDecoder.Decoder)
		for _, rl := range rls {
			rawRes = append(rawRes, rl.Items...)
		}

		// -- types.Resource[T]
		rs, rsErr := decode[types.Resource[json.RawMessage]](supportedDecoder.Decoder)
		rawRes = append(rawRes, rs...)

		// -- input is json/yaml but received error while parsing
		if len(rawRes) > 0 && rlsErr != nil && rsErr != nil {
			return nil, flaterrors.Join(
				errDecodingInput,
				rsErr,
				rlsErr,
			)
		}

		// -- input not decoded: try another supported decoder
		if len(rawRes) == 0 {
			continue
		}

		// -- unmarshal inner objects
		out := make([]types.Resource[types.APIVersionKind], 0)
		for i, r := range rawRes {
			v, err := d.apiServer.Get(types.NewAVKFromResource(r))
			if err != nil {
				return nil, flaterrors.Join(err, fmtErrAtIndex(i))
			}

			if err := supportedDecoder.UnmarshalFunc(r.Spec, v); err != nil {
				return nil, flaterrors.Join(err, fmtErrAtIndex(i))
			}

			out = append(out, v)
		}

		// -- At this point, the output can be safely returned
		return out, nil
	}

	// -- input must be invalid
	return nil, flaterrors.Join(errInputMustBeJsonOrYaml, errDecodingInput)
}

func decode[T any](d decoder) ([]T, error) {
	out := make([]T, 0)
	done := false
	i := 0
	for !done {
		v := new(T)
		if err := d.Decode(v); errors.Is(err, io.EOF) {
			done = true // End of file/stream
		} else if err != nil {
			return nil, flaterrors.Join(err, fmtErrAtIndex(i))
		}
		i++
		out = append(out, *v)
	}
	return out, nil
}

func fmtErrAtIndex(i int) error {
	return fmt.Errorf("error is propably located at index %d", i)
}

// -- utility interface

type decoder interface {
	Decode(v any) error
}
