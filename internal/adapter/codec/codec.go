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
	"encoding/json"
	"fmt"

	"go.yaml.in/yaml/v3"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

var (
	_ types.Codec = &JSONStrategy{}
	_ types.Codec = &YAMLStrategy{}
)

func New(encoding types.Encoding) (types.Codec, error) {
	switch encoding {
	case types.JSONEncoding:
		return &JSONStrategy{}, nil
	case types.YAMLEncoding:
		return &YAMLStrategy{}, nil
	default:
		return nil, flaterrors.Join(
			types.ErrEncoding,
			fmt.Errorf("unrecognized encoding %q", encoding),
		)
	}
}

type JSONStrategy struct{}

func (s *JSONStrategy) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JSONStrategy) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (s *JSONStrategy) Encoding() types.Encoding {
	return types.JSONEncoding
}

type YAMLStrategy struct{}

func (s *YAMLStrategy) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (s *YAMLStrategy) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (s *YAMLStrategy) Encoding() types.Encoding {
	return types.YAMLEncoding
}
