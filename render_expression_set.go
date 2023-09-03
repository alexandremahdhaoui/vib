package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func RenderExpressionSet(resource *api.ResourceDefinition, server api.APIServer) (string, error) {
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

			buffer = JoinLine(buffer, s)
		}

		for _, keyValues := range expressionSet.KeyValues {
			for k, v := range keyValues {
				s, err := Resolve(resolver, k, v)
				if err != nil {
					return "", err
				}

				buffer = JoinLine(buffer, s)
			}
		}

		return buffer, nil
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("APIVersion %q is not supported", resource.APIVersion))
	}
}
