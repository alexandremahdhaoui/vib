//go:build unit

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

package storageadapter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/alexandremahdhaoui/vib/pkg/apis/v1alpha1"
)

const TestingFolder = "/tmp/vib-test"

//----------------------------------------------------------------------------------------------------------------------
// FilesystemOperator
//----------------------------------------------------------------------------------------------------------------------

func folder(t *testing.T) func() {
	t.Helper()

	err := os.Mkdir(TestingFolder, 0755)
	if err != nil {
		t.Error(err)
	}

	return func() { _ = os.RemoveAll(TestingFolder) }
}

func strategy(t *testing.T) *FilesystemOperator {
	t.Helper()

	strategy, err := NewFilesystemOperator(
		apis.V1Alpha1,
		apis.ExpressionKind,
		TestingFolder,
		YAMLEncoding,
	)
	if err != nil {
		t.Error(err)
	}

	return strategy
}

func expressions(t *testing.T) []*ResourceDefinition {
	t.Helper()

	exps := []*ResourceDefinition{
		vib.NewExpression("test0", v1alpha1.ExpressionSpec{
			Key:         "test0",
			Value:       "0",
			ResolverRef: v1alpha1.EnvironmentResolverRef,
		}),
		vib.NewExpression("test1", v1alpha1.ExpressionSpec{
			Key:         "test1",
			Value:       "1",
			ResolverRef: v1alpha1.ExportedEnvironmentResolverRef,
		}),
	}
	return exps
}

func writeExpressions(t *testing.T, exps []*ResourceDefinition) {
	t.Helper()
	for _, exp := range exps {
		fp := filepath.Join(
			TestingFolder,
			fmt.Sprintf(
				"%s.%s.%s.%s",
				cleanAPIVersionForFilesystem(exp.APIVersion),
				exp.Kind,
				exp.Metadata.Name,
				YAMLEncoding,
			),
		)

		err := WriteEncodedFile(fp, *exp)
		if err != nil {
			t.Error(err)
		}
	}
}

// Test the behavior of FilesystemOperator[T]

func TestNewFilesystemOperatorStrategy(t *testing.T) {
	defer folder(t)()
	strategy(t)
}

func TestFilesystemOperatorStrategy_Get(t *testing.T) {
	defer folder(t)()
	str := strategy(t)
	exps := expressions(t)
	writeExpressions(t, exps)

	// Test get with name
	name := exps[0].Metadata.Name

	exp, err := str.Get(&name)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*exps[0], exp[0]) {
		t.Errorf("got: %#v; want: %#v", *exps[0], exp[0])
	}

	// Test failing test case
	name = "something_else"

	_, err = str.Get(&name)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error("got: `err == nil`; want: `err == \"\"`")
	}
}

func TestFilesystemOperatorStrategy_Create(t *testing.T) {
	defer folder(t)()
	str := strategy(t)
	exps := expressions(t)

	for _, exp := range exps {
		err := str.Create(exp)
		if err != nil {
			t.Error(err)
		}

		fp := filepath.Join(
			TestingFolder,
			fmt.Sprintf(
				"%s.%s.%s.%s",
				cleanAPIVersionForFilesystem(exp.APIVersion),
				exp.Kind,
				exp.Metadata.Name,
				YAMLEncoding,
			),
		)

		got, err := ReadEncodedFile(fp)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(got, *exp) {
			t.Errorf("got: %#v; want: %#v", got, *exp)
		}
	}
}
