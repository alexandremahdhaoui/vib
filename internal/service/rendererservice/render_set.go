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
	"errors"
	"fmt"

	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/apis"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func RenderSet(resource *types.Resource, server api.APIServer) (string, error) {
	switch resource.APIVersion {
	case apis.V1Alpha1:
		set := new(v1alpha1.SetSpec)
		err := mapstructure.Decode(resource.Spec, set)
		if err != nil {
			return "", err
		}

		buffer := ""
		for _, expressionRef := range set.ExpressionRefs {
			if err = api.ValidateResourceName(expressionRef); err != nil {
				return "", err
			}

			results := make([]api.ResourceDefinition, 0)
			supportedKinds := []api.Kind{
				apis.ExpressionSetKind,
				apis.ExpressionKind,
			}
			for _, supportedKind := range supportedKinds {
				res, err := server.Get(nil, supportedKind, &expressionRef)
				if err != nil {
					return "", err
				}

				results = append(results, res...)
			}

			if len(results) == 0 {
				return "", api.ErrReference(expressionRef,
					api.Kind(fmt.Sprintf("%#v", supportedKinds)))
			}

			s, err := Render(&results[0], server)
			if err != nil {
				return "", err
			}

			buffer = JoinLine(buffer, s)
		}

		return buffer, nil
	default:
		return "", errors.Join(
			logger.ErrType,
			fmt.Errorf("APIVersion %q is not supported", resource.APIVersion),
		)
	}
}
