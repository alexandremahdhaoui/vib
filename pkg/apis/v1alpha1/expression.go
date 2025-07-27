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
	"github.com/alexandremahdhaoui/vib/internal/types"
)

type ExpressionSpec struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	ResolverRef string `json:"resolverRef"`
}

// APIVersion implements types.DefinedResource.
func (e ExpressionSpec) APIVersion() types.APIVersion {
	return APIVersion
}

// Kind implements types.DefinedResource.
func (e ExpressionSpec) Kind() types.Kind {
	return ExpressionKind
}

// Render implements types.Renderer.
func (e *ExpressionSpec) Render(storage types.Storage) (string, error) {
	// TODO: validate resource name (e.ResolverRef must be valid)

	// TODO: create types.GetFromStorageTyped[T]
	resolver, err := types.GetTypedResourceFromStorage[ResolverSpec](storage, e.ResolverRef)
	if err != nil {
		return "", err
	}

	return resolver.Spec.Resolve(e.Key, e.Value)
}
