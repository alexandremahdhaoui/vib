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
	"github.com/alexandremahdhaoui/vib/internal/util"
)

type ExpressionSetSpec struct {
	// ArbitraryKeys is used for special resolvers, such as plain that does not require associated values
	// ArbitraryKeys are always rendered before KeyValues
	ArbitraryKeys []string `json:"arbitraryKeys"`

	// KeyValues uses a list of map to avoid reordered key-values
	KeyValues []map[string]string `json:"keyValues"`

	// ResolverRef
	ResolverRef types.NamespacedName `json:"resolverRef"`
}

// APIVersion implements types.DefinedResource.
func (e ExpressionSetSpec) APIVersion() types.APIVersion {
	return APIVersion
}

// Kind implements types.DefinedResource.
func (e ExpressionSetSpec) Kind() types.Kind {
	return ExpressionSetKind
}

// Render implements types.Renderer.
func (e *ExpressionSetSpec) Render(storage types.Storage) (string, error) {
	if err := types.ValidateNamespacedName(e.ResolverRef); err != nil {
		return "", err
	}

	resolver, err := types.GetTypedResourceFromStorage(
		storage,
		e.ResolverRef,
		&ResolverSpec{},
	)
	if err != nil {
		return "", err
	}

	buf := ""
	for _, key := range e.ArbitraryKeys {
		s, err := resolver.Spec.Resolve(key, "")
		if err != nil {
			return "", err
		}

		buf = util.JoinLine(buf, s)
	}

	for _, keyValues := range e.KeyValues {
		for k, v := range keyValues {
			s, err := resolver.Spec.Resolve(k, v)
			if err != nil {
				return "", err
			}

			buf = util.JoinLine(buf, s)
		}
	}

	return buf, nil
}
