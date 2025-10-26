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

// ExpressionSetSpec defines the desired state of an ExpressionSet.
// It contains a set of expressions that can be rendered into a desired output and referenced in a profile.
type ExpressionSetSpec struct {
	// ArbitraryKeys is used for special resolvers, such as "plain", that do not require associated values.
	// ArbitraryKeys are always rendered before KeyValues.
	ArbitraryKeys []string `json:"arbitraryKeys"`

	// KeyValues uses a list of maps to avoid reordered key-values.
	KeyValues []map[string]string `json:"keyValues"`

	// ResolverRef is a reference to the Resolver that should be used to render this ExpressionSet.
	ResolverRef types.NamespacedName `json:"resolverRef"`
}

// APIVersion returns the APIVersion of the ExpressionSetSpec.
// It implements the types.DefinedResource interface.
func (e ExpressionSetSpec) APIVersion() types.APIVersion {
	return APIVersion
}

// Kind returns the Kind of the ExpressionSetSpec.
// It implements the types.DefinedResource interface.
func (e ExpressionSetSpec) Kind() types.Kind {
	return ExpressionSetKind
}

// Render renders the ExpressionSetSpec using the specified storage to resolve a resolver.
// It implements the types.Renderer interface.
func (e *ExpressionSetSpec) Render(storage types.Storage) (string, error) {
	resolverRef := defaultRef(e.ResolverRef)
	if err := types.ValidateNamespacedName(resolverRef); err != nil {
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

func defaultRef(nsName types.NamespacedName) types.NamespacedName {
	if nsName.Namespace == "" {
		nsName.Namespace = types.DefaultNamespace
	}
	return nsName
}
