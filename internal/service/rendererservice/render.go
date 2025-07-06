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

package rendererservice

import (
	"fmt"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/apis"
)

type Resolver interface {
	Resolve(key, value string) (string, error)
}

type Renderer interface {
	Render(server adapter.APIServer) (string, error)
}

func Render(resource *types.Resource, server api.APIServer) (string, error) {
	// TODO: This must be automated by registering the APIVersions in main.
	//       And registering the Kinds to the APIVersion in pkg/apis/v1alpha1.
	switch resource.APIVersion {
	case apis.V1Alpha1:
	case default:
	}

	// WARN: this is terrible
	// TODO: we must first check the APIVersion then check the kind

	switch resource.Kind {
	case apis.ProfileKind:
		return RenderProfile(resource, server)
	case apis.SetKind:
		return RenderSet(resource, server)
	case apis.ExpressionKind:
		return RenderExpression(resource, server)
	case apis.ExpressionSetKind:
		return RenderExpressionSet(resource, server)
	default:
		return "", flaterrors.Join(
			types.ErrType,
			fmt.Errorf("Kind %q does not support Render", resource.Kind),
		)
	}
}
