package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

type ExpressionRenderer interface {
	Render(key, value string) (string, error)
}

type Renderer interface {
	Render(server api.APIServer) (string, error)
}

func Render(resource *api.ResourceDefinition, server api.APIServer) (string, error) {
	switch resource.Kind {
	case v1alpha1.ProfileKind:
		spec := new(v1alpha1.ProfileSpec)
		err := mapstructure.Decode(resource.Spec, spec)
		if err != nil {
			return "", err
		}

		return spec.Render(server)
	case v1alpha1.SetKind:
		spec := new(v1alpha1.SetSpec)
		err := mapstructure.Decode(resource.Spec, spec)
		if err != nil {
			return "", err
		}

		return spec.Render(server)
	case v1alpha1.ExpressionKind:
		spec := new(v1alpha1.ExpressionSpec)
		err := mapstructure.Decode(resource.Spec, spec)
		if err != nil {
			return "", err
		}

		return spec.Render(server)
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("Kind %q does not support Render", resource.Kind))
	}
}
