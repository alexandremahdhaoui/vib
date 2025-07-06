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

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

type Encoding string

const (
	JSONEncoding Encoding = "json"
	YAMLEncoding Encoding = "yaml"
)

type Encoder interface {
	Marshal(v *types.Resource) ([]byte, error)
	Unmarshal([]byte) (*types.Resource, error)
	Encoding() Encoding
}

func NewEncoder(encoding Encoding) (Encoder, error) {
	switch encoding {
	case JSONEncoding:
		return &JSONStrategy{}, nil
	case YAMLEncoding:
		return &YAMLStrategy{}, nil
	default:
		return nil, fmt.Errorf("%w: %q", types.ErrEncoding, encoding)
	}
}

func NewEncoderFromFilepath(path string) (Encoder, error) {
	split := strings.Split(path, ".")
	extension := split[len(split)-1]

	encoder, err := NewEncoder(Encoding(extension))
	if err != nil {
		if errors.As(err, &types.ErrEncoding) {
			return nil, flaterrors.Join(err, fmt.Errorf("filepath: %s", path), types.ErrFile)
		}
		return nil, err
	}

	return encoder, nil
}

type JSONStrategy struct{}

func (s *JSONStrategy) Marshal(v *types.Resource) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JSONStrategy) Unmarshal(data []byte) (*types.Resource, error) {
	v := new(types.Resource)
	return v, json.Unmarshal(data, v)
}

func (s *JSONStrategy) Encoding() Encoding {
	return JSONEncoding
}

type YAMLStrategy struct{}

func (s *YAMLStrategy) Marshal(v *types.Resource) ([]byte, error) {
	return yaml.Marshal(v)
}

func (s *YAMLStrategy) Unmarshal(data []byte) (*types.Resource, error) {
	v := new(types.Resource)
	return v, yaml.Unmarshal(data, v)
}

func (s *YAMLStrategy) Encoding() Encoding {
	return YAMLEncoding
}

func ReadEncodedFile(path string) (*types.Resource, error) {
	encoder, err := NewEncoderFromFilepath(path)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return encoder.Unmarshal(b)
}

func WriteEncodedFile(path string, v *types.Resource) error {
	encoder, err := NewEncoderFromFilepath(path)
	if err != nil {
		return err
	}

	b, err := encoder.Marshal(v)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, b, 0666)
	if err != nil {
		return err
	}

	return nil
}
