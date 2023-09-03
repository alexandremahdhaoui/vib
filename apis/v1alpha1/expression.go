package v1alpha1

import (
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/mitchellh/mapstructure"
)

const (
	ExpressionKind = "Expression"
)

type ExpressionSpec struct {
	Key         string `json:"key"         yaml:"key"`
	Value       string `json:"value"       yaml:"value"`
	ResolverRef string `json:"resolverRef" yaml:"resolverRef"`
}

func (e *ExpressionSpec) Render(server api.APIServer) (string, error) {
	if err := api.ValidateResourceName(e.ResolverRef); err != nil {
		return "", err
	}

	results, err := server.Get(nil, ResolverKind, &e.ResolverRef)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", api.ErrReference(e.ResolverRef, ResolverKind)
	}

	resolver := new(ResolverSpec)
	err = mapstructure.Decode(results[0].Spec, resolver)
	if err != nil {
		return "", err
	}

	return resolver.Render(e.Key, e.Value)
}
