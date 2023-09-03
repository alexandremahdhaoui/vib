package api

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"os"
	"path/filepath"
	"reflect"
	"testing"
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

	strategy, err := NewFilesystemOperator(v1alpha1.APIVersion, v1alpha1.ExpressionKind, TestingFolder, YAMLEncoding)
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
			fmt.Sprintf("%s.%s.%s.%s", cleanAPIVersionForFilesystem(exp.APIVersion), exp.Kind, exp.Metadata.Name, YAMLEncoding),
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
			fmt.Sprintf("%s.%s.%s.%s", cleanAPIVersionForFilesystem(exp.APIVersion), exp.Kind, exp.Metadata.Name, YAMLEncoding),
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
