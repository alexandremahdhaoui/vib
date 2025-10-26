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
package main

import (
	"fmt"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	codecadapter "github.com/alexandremahdhaoui/vib/internal/adapter/codec"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

// NewCodec returns a new codec for the given encoding.
func NewCodec(encoding types.Encoding) (types.Codec, error) {
	switch encoding {
	case types.JSONEncoding:
		return codecadapter.NewJSON(), nil
	case types.YAMLEncoding:
		return codecadapter.NewYAML(), nil
	default:
		return nil, flaterrors.Join(
			types.ErrEncoding,
			fmt.Errorf("unrecognized encoding %q", encoding),
		)
	}
}
