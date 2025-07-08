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

import "github.com/alexandremahdhaoui/vib/internal/types"

const (
	ExpressionKind    types.Kind = "Expression"
	ExpressionSetKind types.Kind = "ExpressionSet"
	SetKind           types.Kind = "Set"
	ResolverKind      types.Kind = "Resolver"
	ProfileKind       types.Kind = "Profile"

	APIVersion types.APIVersion = "vib.alexandre.mahdhaoui.com/v1alpha1"
)

var (
	_ types.APIVersionKind = ExpressionSpec{}
	_ types.APIVersionKind = ExpressionSetSpec{}
	_ types.APIVersionKind = ResolverSpec{}
	_ types.APIVersionKind = ProfileSpec{}
)

var kindMap = map[types.Kind]func() types.APIVersionKind{
	ExpressionKind:    func() types.APIVersionKind { return &ExpressionSpec{} },
	ExpressionSetKind: func() types.APIVersionKind { return &ExpressionSetSpec{} },
	ResolverKind:      func() types.APIVersionKind { return &ResolverSpec{} },
	ProfileKind:       func() types.APIVersionKind { return &ProfileSpec{} },
}

func RegisterWithManager(mgr types.APIServer) {
	mgr.Register(APIVersion, kindMap)
}
