package vib_test

import (
	"github.com/alexandremahdhaoui/vib"
	"reflect"
	"testing"
)

func TestResolverSpec_Render(t *testing.T) {
	aliasResolver, err := vib.AliasResolver()
	Err(t, err)

	functionResolver, err := vib.FunctionResolver()
	Err(t, err)

	envResolver, err := vib.EnvironmentResolver()
	Err(t, err)

	exportedEnvResolver, err := vib.ExportedEnvironmentResolver()
	Err(t, err)

	for _, tc := range []struct {
		Name               string
		ResourceDefinition *vib.ResourceDefinition
		Key, Value         string
		Want               string
	}{
		//{
		//	Name:               "ExecResolver",
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
		resolver, ok := tc.ResourceDefinition.Spec.(vib.ResolverSpec)
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
