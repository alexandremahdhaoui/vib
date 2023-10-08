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
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func RenderExpression(resource *api.ResourceDefinition, server api.APIServer) (string, error) {
	switch resource.APIVersion {
	case apis.V1Alpha1:
		expression := new(v1alpha1.ExpressionSpec)
		err := mapstructure.Decode(resource.Spec, expression)
		if err != nil {
			return "", err
		}

		if err := api.ValidateResourceName(expression.ResolverRef); err != nil {
			return "", err
		}

		// apiVersion pointer has to be set to nil, as a V1Alpha1 Expression can legally reference a different
		// APIVersion of a ResolverKind
		results, err := server.Get(nil, apis.ResolverKind, &expression.ResolverRef)
		if err != nil {
			return "", err
		}

		if len(results) == 0 {
			return "", api.ErrReference(expression.ResolverRef, apis.ResolverKind)
		}

		return Resolve(&results[0], expression.Key, expression.Value)
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("APIVersion %q is not supported", resource.APIVersion))
	}
}
