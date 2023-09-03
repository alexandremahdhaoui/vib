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
