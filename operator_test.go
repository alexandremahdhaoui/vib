package vib_test

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib"
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

func strategy(t *testing.T) *vib.FilesystemOperator {
	t.Helper()

	strategy, err := vib.NewFilesystemOperator(vib.V1Alpha1, vib.ExpressionKind, TestingFolder, vib.YAMLEncoding)
	if err != nil {
		t.Error(err)
	}

	return strategy
}

func expressions(t *testing.T) []*vib.ResourceDefinition {
	t.Helper()

	exps := []*vib.ResourceDefinition{
		vib.NewExpression("test0", vib.Expression{
			Key:         "test0",
			Value:       "0",
			ResolverRef: vib.EnvironmentResolverRef,
		}),
		vib.NewExpression("test1", vib.Expression{
			Key:         "test1",
			Value:       "1",
			ResolverRef: vib.ExportedEnvironmentResolverRef,
		}),
	}
	return exps
}

func writeExpressions(t *testing.T, exps []*vib.ResourceDefinition) {
	t.Helper()
	for _, exp := range exps {
		fp := filepath.Join(
			TestingFolder,
			fmt.Sprintf("%s.%s.%s.%s", vib.cleanAPIVersionForFilesystem(exp.APIVersion), exp.Kind, exp.Metadata.Name, vib.YAMLEncoding),
		)

		err := vib.WriteEncodedFile(fp, *exp)
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
			fmt.Sprintf("%s.%s.%s.%s", vib.cleanAPIVersionForFilesystem(exp.APIVersion), exp.Kind, exp.Metadata.Name, vib.YAMLEncoding),
		)

		got, err := vib.ReadEncodedFile(fp)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(got, *exp) {
			t.Errorf("got: %#v; want: %#v", got, *exp)
		}
	}
}
