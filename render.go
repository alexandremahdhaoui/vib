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

package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

type Resolver interface {
	Resolve(key, value string) (string, error)
}

type Renderer interface {
	Render(server api.APIServer) (string, error)
}

func Render(resource *api.ResourceDefinition, server api.APIServer) (string, error) {
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
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("Kind %q does not support Render", resource.Kind))
	}
}
