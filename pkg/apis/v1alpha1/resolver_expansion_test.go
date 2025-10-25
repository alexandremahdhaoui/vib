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

package v1alpha1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateResolverSpec(t *testing.T) {
	// Create a table of test cases
	testCases := []struct {
		name     string
		spec     ResolverSpec
		expected error
	}{
		{
			name: "Exec resolver with nil Exec spec",
			spec: ResolverSpec{
				Type: "exec",
				Exec: nil,
			},
			expected: fmt.Errorf("ResolverSpec.Exec must be set; got: nil"),
		},
		{
			name: "Fmt resolver with nil Fmt spec",
			spec: ResolverSpec{
				Type: "fmt",
				Fmt:  nil,
			},
			expected: fmt.Errorf("ResolverSpec.Fmt must be set; got: nil"),
		},
		{
			name: "Plain resolver with nil Plain spec",
			spec: ResolverSpec{
				Type:  "plain",
				Plain: nil,
			},
			expected: fmt.Errorf("ResolverSpec.Plain must be set; got: nil"),
		},
		{
			name: "Gotemplate resolver with nil GoTemplate spec",
			spec: ResolverSpec{
				Type:       "gotemplate",
				GoTemplate: nil,
			},
			expected: fmt.Errorf("ResolverSpec.Resolver must be set; got: nil"),
		},
		{
			name: "Unknown resolver type",
			spec: ResolverSpec{
				Type: "unknown",
			},
			expected: fmt.Errorf("cannot parse ResolverSpec.Type; got: unknown"),
		},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested
			err := validateResolverSpec(tc.spec)

			// Assert that the error returned by the function is the expected error
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestExecResolverSpec_Resolve(t *testing.T) {
	// Create an ExecResolverSpec
	spec := ExecResolverSpec{
		Command: "echo",
		Args:    []string{"hello"},
	}

	// Assert that calling Resolve panics with the expected message
	assert.PanicsWithValue(t, "ResolveExec is not yet supported", func() {
		_, _ = spec.Resolve("key", "value")
	})
}

func TestGotemplateResolverSpec_Resolve(t *testing.T) {
	// Create a GotemplateResolverSpec
	spec := GotemplateResolverSpec{
		Template: "{{.Key}} {{.Value}}",
	}

	// Assert that calling Resolve panics with the expected message
	assert.PanicsWithValue(t, "unimplemented", func() {
		_, _ = spec.Resolve("key", "value")
	})
}

func TestPlainResolverSpec_Resolve(t *testing.T) {
	// Create a PlainResolverSpec
	spec := PlainResolverSpec(true)

	// Call the Resolve method
	result, err := spec.Resolve("key", "value")

	// Assert that no error was returned
	assert.NoError(t, err)

	// Assert that the result is the expected value
	assert.Equal(t, "key", result)
}