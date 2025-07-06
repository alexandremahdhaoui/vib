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

	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/internal/util"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func RenderExpressionSet(resource *types.Resource, server api.APIServer) (string, error) {
	switch resource.APIVersion {
	case apis.V1Alpha1:
		expressionSet := new(v1alpha1.ExpressionSetSpec)
		err := mapstructure.Decode(resource.Spec, expressionSet)
		if err != nil {
			return "", err
		}

		if err := api.ValidateResourceName(expressionSet.ResolverRef); err != nil {
			return "", err
		}

		results, err := server.Get(nil, apis.ResolverKind, &expressionSet.ResolverRef)
		if err != nil {
			return "", err
		}

		if len(results) == 0 {
			return "", api.ErrReference(expressionSet.ResolverRef, apis.ResolverKind)
		}

		buffer := ""
		resolver := &results[0]

		for _, key := range expressionSet.ArbitraryKeys {
			s, err := Resolve(resolver, key, "")
			if err != nil {
				return "", err
			}

			buffer = util.JoinLine(buffer, s)
		}

		for _, keyValues := range expressionSet.KeyValues {
			for k, v := range keyValues {
				s, err := Resolve(resolver, k, v)
				if err != nil {
					return "", err
				}

				buffer = util.JoinLine(buffer, s)
			}
		}

		return buffer, nil
	default:
		return "", logger.NewErrAndLog(
			logger.ErrType,
			fmt.Sprintf("APIVersion %q is not supported", resource.APIVersion),
		)
	}
}
