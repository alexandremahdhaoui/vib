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
	"github.com/alexandremahdhaoui/vib/internal/types"
)

// TODO: implement apiServer

type (
	avkHash string // concatenate APIVersion and Kind

	// leaf contains information about an AVK
	leaf struct {
		// factory function that instantiate the zero-valued struct
		// corresponding to the AVK.
		avkFactory types.AVKFactory
		// ReaderFactory
		decoderFactory func() types.Reader[types.APIVersionKind]
	}
)

type apiServer struct {
	avkMap map[types.APIVersion]types.Kind
	resMap map[avkHash]leaf
}

func New() types.APIServer {
	return &apiServer{}
}

// Get implements types.APIServer.
func (a *apiServer) Get(avk types.APIVersionKind) (types.Resource[types.APIVersionKind], error) {
	panic("unimplemented")
}

// Register implements types.APIServer.
func (a *apiServer) Register(types.APIVersion, map[types.Kind]func() types.APIVersionKind) {
	panic("unimplemented")
}
