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

package vib_test

import (
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"reflect"
	"testing"
)

func TestResolverSpec_Render(t *testing.T) {
	aliasResolver, err := v1alpha1.AliasResolver()
	Err(t, err)

	functionResolver, err := v1alpha1.FunctionResolver()
	Err(t, err)

	envResolver, err := v1alpha1.EnvironmentResolver()
	Err(t, err)

	exportedEnvResolver, err := v1alpha1.ExportedEnvironmentResolver()
	Err(t, err)

	for _, tc := range []struct {
		Name               string
		ResourceDefinition *api.ResourceDefinition
		Key, Value         string
		Want               string
	}{
		//{
		//	Name:               "ExecResolverSpec",
		//	ResourceDefinition: nil,
		//	Key:                "echo",
		//	Value:              "this value",
		//	Want:               "echo this value",
		//},
		{
			Name:               "AliasResolver",
			ResourceDefinition: aliasResolver,
			Key:                "test",
			Value:              "echo 0",
			Want:               "alias test='echo 0'",
		},
		{
			Name:               "FunctionResolver",
			ResourceDefinition: functionResolver,
			Key:                "test1",
			Value:              "echo $1",
			Want:               "test1() { echo $1 ; }",
		},
		{
			Name:               "EnvironmentResolver",
			ResourceDefinition: envResolver,
			Key:                "TEST",
			Value:              "2",
			Want:               "TEST=\"2\"",
		},
		{
			Name:               "ExportedEnvironmentResolver",
			ResourceDefinition: exportedEnvResolver,
			Key:                "TEST",
			Value:              "3",
			Want:               "export TEST=\"3\"",
		},
	} {
		resolver, ok := tc.ResourceDefinition.Spec.(v1alpha1.ResolverSpec)
		if !ok {
			t.Error("cannot type assert spec as vib.ResolverSpec")
		}
		got, err := resolver.Render(tc.Key, tc.Value)
		Err(t, err)

		if !reflect.DeepEqual(got, tc.Want) {
			t.Errorf("got: %#v; want: %#v", got, tc.Want)
		}
	}
}

func Err(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
