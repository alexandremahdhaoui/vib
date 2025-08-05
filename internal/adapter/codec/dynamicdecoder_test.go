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

package codecadapter_test

import (
	"bytes"
	"testing"

	codecadapter "github.com/alexandremahdhaoui/vib/internal/adapter/codec"
	"github.com/alexandremahdhaoui/vib/internal/service"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestDynamicResourceDecoder(t *testing.T) {
	var (
		drd types.DynamicDecoder[types.APIVersionKind]
		_   = ""
	)

	setup := func(t *testing.T) {
		t.Helper()

		apiServer := service.NewAPIServer()
		v1alpha1.RegisterWithManager(apiServer)

		drd = codecadapter.NewDynamicResourceDecoder(apiServer)
	}

	t.Run("Decode", func(t *testing.T) {
		setup(t)

		var expected any
		buf := bytes.NewBufferString(`
---
items:
  - apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
    kind: ExpressionSet
    metadata:
      name: ItemFromExpressionSet
  - apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
    kind: ExpressionSet
    metadata:
      name: test-test
`)

		actual, err := drd.Decode(buf)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
