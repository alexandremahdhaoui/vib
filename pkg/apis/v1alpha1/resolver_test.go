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

package v1alpha1_test

import (
	"reflect"
	"testing"

	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/apis/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func TestResolverSpec_Render(t *testing.T) {
	for _, tc := range []struct {
		Name       string
		Resource   *types.Resource
		Key, Value string
		Want       string
	}{
		//{
		//	Name:               "ExecResolverSpec",
		//	ResourceDefinition: nil,
		//	Key:                "echo",
		//	Value:              "this value",
		//	Want:               "echo this value",
		//},

		{
			Name:     "AliasResolver",
			Resource: Must(t, v1alpha1.NewAliasResolver),
			Key:      "test",
			Value:    "echo 0",
			Want:     "alias test='echo 0'",
		},

		{
			Name:     "FunctionResolver",
			Resource: Must(t, v1alpha1.NewFunctionResolver),
			Key:      "test1",
			Value:    "echo $1",
			Want:     "test1() { echo $1 ; }",
		},

		{
			Name:     "EnvironmentResolver",
			Resource: Must(t, v1alpha1.NewEnvironmentResolver),
			Key:      "TEST",
			Value:    "2",
			Want:     "TEST=\"2\"",
		},

		{
			Name:     "ExportedEnvironmentResolver",
			Resource: Must(t, v1alpha1.NewExportedEnvironmentResolver),
			Key:      "TEST",
			Value:    "3",
			Want:     "export TEST=\"3\"",
		},
	} {
		resolver, ok := tc.Resource.Spec.(v1alpha1.ResolverSpec)
		if !ok {
			t.Error("cannot type assert spec as vib.ResolverSpec")
		}

		got, err := resolver.Resolve(tc.Key, tc.Value)
		assert.NoError(t, err)

		if !reflect.DeepEqual(got, tc.Want) {
			t.Errorf("got: %#v; want: %#v", got, tc.Want)
		}
	}
}

func Must[T any](t *testing.T, f func() (T, error)) T {
	t.Helper()
	out, err := f()
	assert.NoError(t, err)
	return out
}
