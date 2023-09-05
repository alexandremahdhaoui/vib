package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func RenderProfile(resource *api.ResourceDefinition, server api.APIServer) (string, error) {
	switch resource.APIVersion {
	case apis.V1Alpha1:
		profile := new(v1alpha1.ProfileSpec)
		err := mapstructure.Decode(resource.Spec, profile)
		if err != nil {
			return "", err
		}

		buffer := ""
		for _, ref := range profile.SetRefs {
			if err = api.ValidateResourceName(ref); err != nil {
				return "", err
			}

			results := make([]api.ResourceDefinition, 0)
			supportedKinds := []api.Kind{
				apis.SetKind,
				apis.ExpressionSetKind,
			}
			for _, supportedKind := range supportedKinds {
				res, err := server.Get(nil, supportedKind, &ref)
				if err != nil {
					return "", err
				}
				results = append(results, res...)
			}

			if len(results) == 0 {
				return "", api.ErrReference(ref,
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
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("APIVersion %q is not supported", resource.APIVersion))
	}
}
