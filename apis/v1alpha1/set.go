package v1alpha1

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

const (
	SetKind = "Set"
)

type SetSpec struct {
	// ExpressionRefs is a list of reference to vib.Expression
	ExpressionRefs []string `json:"expressionRefs" yaml:"expressionRefs"`
}

func (p *SetSpec) Render(server api.APIServer) (string, error) {
	var rendered string

	for _, ref := range p.ExpressionRefs {
		if err := api.ValidateResourceName(ref); err != nil {
			return "", err
		}

		results, err := server.Get(nil, ExpressionKind, &ref)
		if err != nil {
			return "", err
		}

		if len(results) == 0 {
			return "", api.ErrReference(ref, ExpressionKind)
		}

		expression := new(ExpressionSpec)
		err = mapstructure.Decode(results[0].Spec, expression)
		if err != nil {
			return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("cannot assert type ExpressionSpec for resource %#v", results))
		}

		s, err := expression.Render(server)
		if err != nil {
			return "", err
		}

		rendered = fmt.Sprintf("%s\n%s", rendered, s)
	}

	return rendered, nil

}

func NewSet(name string, spec SetSpec) *api.ResourceDefinition {
	return &api.ResourceDefinition{
		APIVersion: APIVersion,
		Kind:       SetKind,
		Metadata:   api.NewMetadata(name),
		Spec:       spec,
	}
}
