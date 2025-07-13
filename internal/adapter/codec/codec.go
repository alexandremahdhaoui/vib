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
	_ types.Codec = &jsonCodec{}
	_ types.Codec = &yamlCodec{}
)

func New(encoding types.Encoding) (types.Codec, error) {
	switch encoding {
	case types.JSONEncoding:
		return &jsonCodec{}, nil
	case types.YAMLEncoding:
		return &yamlCodec{}, nil
	default:
		return nil, flaterrors.Join(
			types.ErrEncoding,
			fmt.Errorf("unrecognized encoding %q", encoding),
		)
	}
}

type jsonCodec struct{}

func (s *jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (s *jsonCodec) Encoding() types.Encoding {
	return types.JSONEncoding
}

type yamlCodec struct{}

func (s *yamlCodec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (s *yamlCodec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (s *yamlCodec) Encoding() types.Encoding {
	return types.YAMLEncoding
}
