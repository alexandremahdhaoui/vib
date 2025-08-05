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

	"github.com/alexandremahdhaoui/vib/internal/types"

	// INFO: Package "sigs.k8s.io/yaml" ensures that json tags specifying
	// omitempty are respected
	"sigs.k8s.io/yaml"
)

var (
	_ types.Codec = &jsonCodec{}
	_ types.Codec = &yamlCodec{}
)

func NewJSON() types.Codec {
	return &jsonCodec{}
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

func NewYAML() types.Codec {
	return &yamlCodec{}
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
