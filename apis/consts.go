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

package apis

import "github.com/alexandremahdhaoui/vib/pkg/api"

const (
	ExpressionKind    api.Kind = "Expression"
	ExpressionSetKind api.Kind = "ExpressionSet"
	SetKind           api.Kind = "Set"
	ResolverKind      api.Kind = "Resolver"
	ProfileKind       api.Kind = "Profile"

	V1Alpha1 api.APIVersion = "vib.alexandre.mahdhaoui.com/v1alpha1"
)
