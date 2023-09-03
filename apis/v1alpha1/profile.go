package v1alpha1

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

const (
	ProfileKind api.Kind = "Profile"
)

type ProfileSpec struct {
	// SetRefs is a list of reference to vib.Set.
	SetRefs []string `json:"setRefs" yaml:"setRefs"`
}

func (p *ProfileSpec) Render(server api.APIServer) (string, error) {
	var rendered string

	for _, ref := range p.SetRefs {
		if err := api.ValidateResourceName(ref); err != nil {
			return "", err
		}

		results, err := server.Get(nil, SetKind, &ref)
		if err != nil {
			return "", err
		}

		if len(results) == 0 {
			return "", api.ErrReference(ref, SetKind)
		}

		set := new(SetSpec)
		err = mapstructure.Decode(results[0].Spec, set)
		if err != nil {
			return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("cannot assert type SetSpec for resource %#v", results))
		}

		s, err := set.Render(server)
		if err != nil {
			return "", err
		}

		rendered = fmt.Sprintf("%s\n%s", rendered, s)
	}

	return rendered, nil
}

func NewProfile(name string, spec ProfileSpec) *api.ResourceDefinition {
	return &api.ResourceDefinition{
		APIVersion: APIVersion,
		Kind:       ProfileKind,
		Metadata:   api.NewMetadata(name),
		Spec:       spec,
	}
}
