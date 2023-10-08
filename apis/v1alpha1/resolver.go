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

package v1alpha1

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/pkg/api"
)

type FmtArgument string

const (
	ExecResolverType       = "exec"
	FmtResolverType        = "fmt"
	PlainResolverType      = "plain"
	GotemplateResolverType = "gotemplate"

	PlainResolverRef               = "plain"
	FunctionResolverRef            = "function"
	AliasResolverRef               = "alias"
	EnvironmentResolverRef         = "environment"
	ExportedEnvironmentResolverRef = "environment-exported"

	KeyFmtArgument   FmtArgument = "key"
	ValueFmtArgument FmtArgument = "value"
)

type (
	ResolverSpec struct {
		Type       string                  `json:"type"                 yaml:"type"`
		Exec       *ExecResolverSpec       `json:"exec,omitempty"       yaml:"exec,omitempty"`
		Fmt        *FmtResolverSpec        `json:"fmt,omitempty"        yaml:"fmt,omitempty"`
		Plain      *PlainResolverSpec      `json:"plain,omitempty"      yaml:"plain,omitempty"`
		GoTemplate *GotemplateResolverSpec `json:"gotemplate,omitempty" yaml:"gotemplate,omitempty"`
	}

	ExecResolverSpec struct {
		Command string   `json:"command" yaml:"command"`
		Args    []string `json:"args,omitempty" yaml:"args,omitempty"`

		// Stdin is a format-able
		Stdin string `json:"stdin,omitempty" yaml:"stdin,omitempty"`
	}

	FmtResolverSpec struct {
		Template string `json:"template" yaml:"template"`
		// FmtArguments is a list of FmtArgument, that will be used to format the template
		FmtArguments []FmtArgument `json:"fmtArguments" yaml:"fmtArguments"`
	}

	PlainResolverSpec bool

	GotemplateResolverSpec struct {
		Template string `json:"template" yaml:"template"`
	}
)

func NewResolver(name string, spec ResolverSpec) (*api.ResourceDefinition, error) {
	if err := validateResolverSpec(&spec); err != nil {
		return nil, err
	}

	return &api.ResourceDefinition{
			APIVersion: apis.V1Alpha1,
			Kind:       apis.ResolverKind,
			Metadata:   api.NewMetadata(name),
			Spec:       spec,
		},
		nil
}

func validateResolverSpec(spec *ResolverSpec) error {
	switch spec.Type {
	case ExecResolverType:
		if spec.Exec == nil {
			return fmt.Errorf("ResolverSpec.Exec must be set; got: nil")
		}
	case FmtResolverType:
		if spec.Fmt == nil {
			return fmt.Errorf("ResolverSpec.Fmt must be set; got: nil")
		}
	case PlainResolverType:
		if spec.Plain == nil {
			return fmt.Errorf("ResolverSpec.Plain must be set; got: nil")
		}
	case GotemplateResolverType:
		if spec.GoTemplate == nil {
			return fmt.Errorf("ResolverSpec.Resolver must be set; got: nil")
		}
	default:
		return fmt.Errorf("couldn't parse ResolverSpec.Type; got: %s", spec.Type)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// PlainResolver
//----------------------------------------------------------------------------------------------------------------------

func PlainResolver() (*api.ResourceDefinition, error) {
	return NewResolver(
		PlainResolverRef,
		ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
			Type:  PlainResolverRef,
			Plain: api.ToPointer(PlainResolverSpec(true)),
		},
	)
}

//----------------------------------------------------------------------------------------------------------------------
// FunctionResolver
//----------------------------------------------------------------------------------------------------------------------

func FunctionResolver() (*api.ResourceDefinition, error) {
	return NewResolver(
		FunctionResolverRef,
		ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
			Type: FmtResolverType,
			Fmt: &FmtResolverSpec{
				Template:     "%s() {\n%s\n}",
				FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
			},
		},
	)
}

//----------------------------------------------------------------------------------------------------------------------
// AliasResolver
//----------------------------------------------------------------------------------------------------------------------

func AliasResolver() (*api.ResourceDefinition, error) {
	return NewResolver(
		AliasResolverRef,
		ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
			Type: FmtResolverType,
			Fmt: &FmtResolverSpec{
				Template:     "alias %s='%s'",
				FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
			},
		},
	)
}

//----------------------------------------------------------------------------------------------------------------------
// EnvironmentResolver
//----------------------------------------------------------------------------------------------------------------------

func EnvironmentResolver() (*api.ResourceDefinition, error) {
	return NewResolver(
		EnvironmentResolverRef,
		ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
			Type: FmtResolverType,
			Fmt: &FmtResolverSpec{
				Template:     "%s=%q",
				FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
			},
		},
	)
}

//----------------------------------------------------------------------------------------------------------------------
// ExportedEnvironmentResolver
//----------------------------------------------------------------------------------------------------------------------

func ExportedEnvironmentResolver() (*api.ResourceDefinition, error) {
	return NewResolver(
		ExportedEnvironmentResolverRef,
		ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
			Type: FmtResolverType,
			Fmt: &FmtResolverSpec{
				Template:     "export %s=%q",
				FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
			},
		},
	)
}
