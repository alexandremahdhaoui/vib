package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func RenderSet(resource *api.ResourceDefinition, server api.APIServer) (string, error) {
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

			for _, supportedKind := range []api.Kind{
				apis.ExpressionSetKind,
				apis.ExpressionKind,
			} {

				res, err := server.Get(nil, supportedKind, &expressionRef)
				if err != nil {
					return "", err
				}

				results = append(results, res...)
			}

			if len(results) == 0 {
				return "", api.ErrReference(expressionRef,
					api.Kind(fmt.Sprintf("{%s,%s}", apis.ExpressionKind, apis.ExpressionSetKind)))
			}

			s, err := Render(&results[0], server)
			if err != nil {
				return "", err
			}

			buffer = JoinLine(buffer, s)
		}

		return buffer, nil
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("APIVersion %q is not supported", resource.APIVersion))
	}

}
