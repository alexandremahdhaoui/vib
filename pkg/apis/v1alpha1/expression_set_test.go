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
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockStorage is a mock implementation of the types.Storage interface
type MockStorage struct {
	getResource    func(avk types.APIVersionKind, namespacedName types.NamespacedName) (types.Resource[types.APIVersionKind], error)
	listResources  func(avk types.APIVersionKind, namespace string) ([]types.Resource[types.APIVersionKind], error)
	createResource func(resource types.Resource[types.APIVersionKind]) error
	updateResource func(resource types.Resource[types.APIVersionKind]) error
	deleteResource func(avk types.APIVersionKind, namespacedName types.NamespacedName) error
}

func (m *MockStorage) Get(avk types.APIVersionKind, namespacedName types.NamespacedName) (types.Resource[types.APIVersionKind], error) {
	if m.getResource != nil {
		return m.getResource(avk, namespacedName)
	}
	return types.Resource[types.APIVersionKind]{}, fmt.Errorf("Get not implemented")
}

func (m *MockStorage) List(avk types.APIVersionKind, namespace string) ([]types.Resource[types.APIVersionKind], error) {
	if m.listResources != nil {
		return m.listResources(avk, namespace)
	}
	return nil, fmt.Errorf("List not implemented")
}

func (m *MockStorage) Create(resource types.Resource[types.APIVersionKind]) error {
	if m.createResource != nil {
		return m.createResource(resource)
	}
	return fmt.Errorf("Create not implemented")
}

func (m *MockStorage) Update(resource types.Resource[types.APIVersionKind]) error {
	if m.updateResource != nil {
		return m.updateResource(resource)
	}
	return fmt.Errorf("Update not implemented")
}

func (m *MockStorage) Delete(avk types.APIVersionKind, namespacedName types.NamespacedName) error {
	if m.deleteResource != nil {
		return m.deleteResource(avk, namespacedName)
	}
	return fmt.Errorf("Delete not implemented")
}

func TestExpressionSetSpec_Render(t *testing.T) {
	// Create a mock storage
	mockStorage := &MockStorage{
		getResource: func(avk types.APIVersionKind, namespacedName types.NamespacedName) (types.Resource[types.APIVersionKind], error) {
			if namespacedName.Name == "my-resolver" && namespacedName.Namespace == "default" {
				return types.Resource[types.APIVersionKind]{
					APIVersion: APIVersion,
					Kind:       ResolverKind,
					Metadata: types.Metadata{
						Name:      "my-resolver",
						Namespace: "default",
					},
					Spec: &ResolverSpec{
						Type: FmtResolverType,
						Fmt: &FmtResolverSpec{
							Template:     "key=%s, value=%s",
							FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
						},
					},
				}, nil
			}
			return types.Resource[types.APIVersionKind]{}, fmt.Errorf("resolver not found")
		},
	}

	// Create an ExpressionSetSpec
	spec := &ExpressionSetSpec{
		ArbitraryKeys: []string{"key1", "key2"},
		KeyValues: []map[string]string{
			{"key3": "value3"},
			{"key4": "value4"},
		},
		ResolverRef: types.NamespacedName{
			Name:      "my-resolver",
			Namespace: "default",
		},
	}

	// Call the Render method
	result, err := spec.Render(mockStorage)

	// Assert that no error was returned
	assert.NoError(t, err)

	// Assert that the result is the expected value
	expected := "key=key1, value=\nkey=key2, value=\nkey=key3, value=value3\nkey=key4, value=value4"
	assert.Equal(t, expected, result)
}

func TestExpressionSetSpec_Render_Error(t *testing.T) {
	// Create a mock storage
	mockStorage := &MockStorage{
		getResource: func(avk types.APIVersionKind, namespacedName types.NamespacedName) (types.Resource[types.APIVersionKind], error) {
			return types.Resource[types.APIVersionKind]{}, fmt.Errorf("resolver not found")
		},
	}

	// Create an ExpressionSetSpec
	spec := &ExpressionSetSpec{
		ResolverRef: types.NamespacedName{
			Name:      "my-resolver",
			Namespace: "default",
		},
	}

	// Call the Render method
	_, err := spec.Render(mockStorage)

	// Assert that an error was returned
	assert.Error(t, err)
}

func TestDefaultRef(t *testing.T) {
	// Create a NamespacedName with a namespace
	nsName := types.NamespacedName{
		Name:      "my-name",
		Namespace: "my-namespace",
	}

	// Call the defaultRef function
	result := defaultRef(nsName)

	// Assert that the namespace is unchanged
	assert.Equal(t, "my-namespace", result.Namespace)

	// Create a NamespacedName without a namespace
	nsName = types.NamespacedName{
		Name: "my-name",
	}

	// Call the defaultRef function
	result = defaultRef(nsName)

	// Assert that the namespace is the default namespace
	assert.Equal(t, types.DefaultNamespace, result.Namespace)
}