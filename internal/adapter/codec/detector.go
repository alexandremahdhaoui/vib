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
	"io"

	"github.com/alexandremahdhaoui/vib/internal/types"
)

func NewDetector(reader io.Reader) types.Codec {
	return detector{}
}

type detector struct{}

// Encoding implements types.Codec.
func (d detector) Encoding() types.Encoding {
	panic("unimplemented")
}

// Marshal implements types.Codec.
func (d detector) Marshal(v any) ([]byte, error) {
	panic("unimplemented")
}

// Unmarshal implements types.Codec.
func (d detector) Unmarshal(b []byte, v any) error {
	panic("unimplemented")
}

func test() {
	json.NewDecoder().Decode()
}
