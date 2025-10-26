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

package service

import (
	"fmt"
	"strings"

	"github.com/alexandremahdhaoui/vib/internal/types"
)

type (
	avkHash string // concatenate APIVersion and Kind
	// leaf contains information about an AVK
	leaf struct {
		// avkFactory function that instantiate the zero-valued struct
		// corresponding to the AVK.
		avkFactory types.AVKFunc
	}
)

// apiServer implements the types.APIServer interface.
type apiServer struct {
	leavesByHash          map[avkHash]leaf
	registeredAPIVersions []types.APIVersion
}

// NewAPIServer returns a new APIServer.
func NewAPIServer() types.APIServer {
	return &apiServer{
		leavesByHash: make(map[avkHash]leaf),
	}
}

// Get implements the types.APIServer interface.
func (a *apiServer) Get(avk types.APIVersionKind) (types.Resource[types.APIVersionKind], error) {
	l, err := a.getLeaf(avk)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	spec := l.avkFactory()
	return types.Resource[types.APIVersionKind]{
		APIVersion: spec.APIVersion(),
		Kind:       spec.Kind(),
		Metadata:   types.Metadata{},
		Spec:       spec,
	}, nil
}

// Register implements the types.APIServer interface.
func (a *apiServer) Register(avkFactory []types.AVKFunc) {
	for _, f := range avkFactory {
		avk := f()
		a.registeredAPIVersions = append(a.registeredAPIVersions, avk.APIVersion())
		a.leavesByHash[a.computeAVKHash(avk)] = leaf{
			avkFactory: f,
		}
	}
}

func (a *apiServer) computeAVKHash(avk types.APIVersionKind) avkHash {
	return avkHash(
		fmt.Sprintf(
			"%s-%s",
			strings.ToLower(avk.APIVersion()),
			strings.ToLower(avk.Kind())),
	)
}

func (a *apiServer) getLeaf(avk types.APIVersionKind) (leaf, error) {
	if v := avk.APIVersion(); v == "" {
		for _, v := range a.registeredAPIVersions {
			newAVK := types.NewAPIVersionKind(v, avk.Kind())
			l, ok := a.leavesByHash[a.computeAVKHash(newAVK)]
			if ok {
				return l, nil
			}
		}

		return leaf{}, types.ErrNotFound
	}

	l, ok := a.leavesByHash[a.computeAVKHash(avk)]
	if !ok {
		return leaf{}, types.ErrNotFound
	}

	return l, nil
}
