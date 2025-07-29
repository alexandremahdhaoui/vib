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
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

// TODO: write unit tests

// TODO: Not urgent:
// - LTS would be adding Encoder/Decoder iface to "sigs.k8s.io/yaml"

// Instantiate a new dynamic resource decoder. A dynamic resource decoder is a special codec
// that can unmarshal one or many documents from any supported encoding.
func NewDynamicResourceDecoder(
	apiServer types.APIServer,
) types.DynamicDecoder[types.APIVersionKind] {
	return &rawDrd{
		apiServer: apiServer,
	}
}

type rawDrd struct {
	apiServer types.APIServer
}

var (
	errInputMustBeJsonOrYaml = errors.New("input must be json or yaml")
	errDecodingInput         = errors.New("error decoding input")
	errInputMustNotBeEmpty   = errors.New("input must not be empty")
)

type resourceList[T any] struct {
	Items []types.Resource[T] `json:"items" yaml:"items"`
}

// Decode implements types.DynamicDecoder.
func (d *rawDrd) Decode(reader io.Reader) ([]types.Resource[types.APIVersionKind], error) {
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
		IsJSON        bool
		Decoder       decoder
		UnmarshalFunc func(b []byte, v any) error
		MarshalFunc   func(v any) ([]byte, error)
	}{
		{
			IsJSON:        true,
			Decoder:       json.NewDecoder(jsonBuf),
			UnmarshalFunc: json.Unmarshal,
			MarshalFunc:   json.Marshal,
		},

		{
			Decoder:       yaml.NewDecoder(yamlBuf),
			UnmarshalFunc: yaml.Unmarshal,
			MarshalFunc:   yaml.Marshal,
		},
	} {
		objectList, err := decodeRaw(supportedDecoder.Decoder)
		if len(objectList) > 0 && err != nil {
			// -- input is json/yaml but received error while parsing
			return nil, flaterrors.Join(err, errDecodingInput)
		} else if err != nil {
			// -- input not decoded: try another supported decoder
			continue
		}

		out := make([]types.Resource[types.APIVersionKind], 0)
		for i, obj := range objectList {
			if v, ok := obj["items"]; ok {
				// TODO: handle list items
				list, ok := v.([]any)
				if !ok {
					return nil, flaterrors.Join(errors.New("TODO"), fmtErrAtIndex(i))
				}
				for _, raw := range list {
					item, err := d.decodeOne(raw)
					if err != nil {
						return nil, flaterrors.Join(err, fmtErrAtIndex(i))
					}
					out = append(out, item)
				}
			} else {
				item, err := d.decodeOne(obj)
				if err != nil {
					return nil, flaterrors.Join(err, fmtErrAtIndex(i))
				}
				out = append(out, item)
			}
		}

		// -- At this point, the output can be safely returned
		return out, nil
	}

	// -- invalid input
	return nil, flaterrors.Join(errInputMustBeJsonOrYaml, errDecodingInput)
}

func (d *rawDrd) decodeOne(v any) (types.Resource[types.APIVersionKind], error) {
	m, ok := v.(map[string]any)
	if !ok {
		return types.Resource[types.APIVersionKind]{}, errors.New("TODO: decodeOne err")
	}

	// Ignore types assertion as "apiServer.Get" will return ERRNOTFOUND in the worst case scenario
	apiVersion, _ := m["apiVersion"].(string)
	kind, _ := m["kind"].(string)
	avk := types.NewAPIVersionKind(apiVersion, kind)
	out, err := d.apiServer.Get(avk)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	if err = json.Unmarshal(b, &out); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	return out, nil
}

func NewRawDynamicResourceDecoder(
	apiServer types.APIServer,
) types.DynamicDecoder[types.APIVersionKind] {
	return &rawDrd{
		apiServer: apiServer,
	}
}

// raw decoding into map[any]any
func decodeRaw(d decoder) ([]map[any]any, error) {
	out := make([]map[any]any, 0)
	done := false
	i := 0
	for !done {
		var v map[any]any
		if err := d.Decode(&v); errors.Is(err, io.EOF) {
			done = true // End of file/stream
		} else if err != nil {
			return nil, flaterrors.Join(err, fmtErrAtIndex(i))
		}
		i++
		out = append(out, v)
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
