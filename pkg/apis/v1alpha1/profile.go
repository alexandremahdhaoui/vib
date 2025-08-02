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

type ProfileSpec struct {
	// Refs is a list of reference to Expression or ExpressionSet
	// TODO: ExpressionRefs must be validated to ensure no duplication.
	// Name duplication would yield unexpected behavior
	Refs []types.NamespacedName `json:"refs"`
}

// APIVersion implements types.DefinedResource.
func (p ProfileSpec) APIVersion() types.APIVersion {
	return APIVersion
}

// Kind implements types.DefinedResource.
func (p ProfileSpec) Kind() types.Kind {
	return ProfileKind
}

// Render implements types.Renderer.
func (p *ProfileSpec) Render(storage types.Storage) (string, error) {
	refs := make(map[types.NamespacedName]string, len(p.Refs))
	namespaces := make(map[string]struct{})
	for _, ref := range p.Refs {
		if err := types.ValidateNamespacedName(ref); err != nil {
			return "", err
		}

		refs[ref] = ""
		namespaces[ref.Namespace] = struct{}{}
	}

	// -- Expressions
	eList := make([]types.Resource[*ExpressionSpec], 0)
	for ns := range namespaces {
		l, err := types.ListTypedResourceFromStorage(
			storage,
			ns,
			&ExpressionSpec{},
		)
		if err != nil {
			return "", err
		}

		eList = append(eList, l...)
	}

	for _, e := range eList {
		nsName := types.NewNamespacedNameFromMetadata(e.Metadata)
		if _, ok := refs[nsName]; !ok {
			continue
		}

		s, err := e.Spec.Render(storage)
		if err != nil {
			return "", err
		}

		refs[nsName] = s
	}

	// -- Expression sets
	esList := make([]types.Resource[*ExpressionSetSpec], 0)
	for ns := range namespaces {
		l, err := types.ListTypedResourceFromStorage(
			storage,
			ns,
			&ExpressionSetSpec{},
		)
		if err != nil {
			return "", err
		}

		esList = append(esList, l...)
	}

	for _, es := range esList {
		nsName := types.NewNamespacedNameFromMetadata(es.Metadata)
		if _, ok := refs[nsName]; !ok {
			continue
		}

		s, err := es.Spec.Render(storage)
		if err != nil {
			return "", err
		}

		refs[nsName] = s
	}

	buf := ""
	for _, ref := range p.Refs {
		buf = util.JoinLine(buf, refs[ref])
		delete(refs, ref)
	}

	return buf, nil
}
