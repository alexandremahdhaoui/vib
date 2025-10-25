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

func TestProfileSpec_Render(t *testing.T) {
	// Create a mock storage
	mockStorage := &MockStorage{
		listResources: func(avk types.APIVersionKind, namespace string) ([]types.Resource[types.APIVersionKind], error) {
			if namespace == "default" {
				return []types.Resource[types.APIVersionKind]{
					{
						APIVersion: APIVersion,
						Kind:       ExpressionSetKind,
						Metadata: types.Metadata{
							Name:      "my-expression-set",
							Namespace: "default",
						},
						Spec: &ExpressionSetSpec{
							ArbitraryKeys: []string{"key1"},
							ResolverRef: types.NamespacedName{
								Name:      "my-resolver",
								Namespace: "default",
							},
						},
					},
				}, nil
			}
			return nil, nil
		},
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
						Type:  PlainResolverType,
						Plain: new(PlainResolverSpec),
					},
				}, nil
			}
			return types.Resource[types.APIVersionKind]{}, fmt.Errorf("resolver not found")
		},
	}

	// Create a ProfileSpec
	spec := &ProfileSpec{
		Refs: []types.NamespacedName{
			{
				Name:      "my-expression-set",
				Namespace: "default",
			},
		},
	}

	// Call the Render method
	result, err := spec.Render(mockStorage)

	// Assert that no error was returned
	assert.NoError(t, err)

	// Assert that the result is the expected value
	assert.Equal(t, "key1", result)
}

func TestProfileSpec_Render_ListError(t *testing.T) {
	// Create a mock storage
	mockStorage := &MockStorage{
		listResources: func(avk types.APIVersionKind, namespace string) ([]types.Resource[types.APIVersionKind], error) {
			return nil, fmt.Errorf("list error")
		},
	}

	// Create a ProfileSpec
	spec := &ProfileSpec{
		Refs: []types.NamespacedName{
			{
				Name:      "my-expression-set",
				Namespace: "default",
			},
		},
	}

	// Call the Render method
	_, err := spec.Render(mockStorage)

	// Assert that an error was returned
	assert.Error(t, err)
}

func TestProfileSpec_Render_RenderError(t *testing.T) {
	// Create a mock storage
	mockStorage := &MockStorage{
		listResources: func(avk types.APIVersionKind, namespace string) ([]types.Resource[types.APIVersionKind], error) {
			if namespace == "default" {
				return []types.Resource[types.APIVersionKind]{
					{
						APIVersion: APIVersion,
						Kind:       ExpressionSetKind,
						Metadata: types.Metadata{
							Name:      "my-expression-set",
							Namespace: "default",
						},
						Spec: &ExpressionSetSpec{
							ResolverRef: types.NamespacedName{
								Name:      "my-resolver",
								Namespace: "default",
							},
						},
					},
				}, nil
			}
			return nil, nil
		},
		getResource: func(avk types.APIVersionKind, namespacedName types.NamespacedName) (types.Resource[types.APIVersionKind], error) {
			return types.Resource[types.APIVersionKind]{}, fmt.Errorf("resolver not found")
		},
	}

	// Create a ProfileSpec
	spec := &ProfileSpec{
		Refs: []types.NamespacedName{
			{
				Name:      "my-expression-set",
				Namespace: "default",
			},
		},
	}

	// Call the Render method
	_, err := spec.Render(mockStorage)

	// Assert that an error was returned
	assert.Error(t, err)
}