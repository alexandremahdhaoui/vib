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

	V1Alpha1 types.APIVersion = "vib.alexandre.mahdhaoui.com/v1alpha1"
)

func RegisterToServer(server types.APIServer) {
	server.Register(
		V1Alpha1,
		map[types.Kind]func() any{
			ExpressionKind:    func() any { return ExpressionSpec{} },
			ExpressionSetKind: func() any { return ExpressionSetSpec{} },
			SetKind:           func() any { return SetSpec{} },
			ResolverKind:      func() any { return ResolverSpec{} },
			ProfileKind:       func() any { return ProfileSpec{} },
		},
	)
}
