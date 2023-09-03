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
