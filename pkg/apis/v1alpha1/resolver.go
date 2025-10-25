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
	"os/exec"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"

	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/internal/util"
)

var (
	_ types.APIVersionKind = &ResolverSpec{}
	_ Resolver             = &ResolverSpec{}

	_ Resolver = &ExecResolverSpec{}
	_ Resolver = &FmtResolverSpec{}
	_ Resolver = &GotemplateResolverSpec{}
	_ Resolver = util.Ptr(PlainResolverSpec(true))
)

type Resolver interface {
	Resolve(key, value string) (string, error)
}

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

type ResolverSpec struct {
	Type       string                  `json:"type"`
	Exec       *ExecResolverSpec       `json:"exec,omitempty"`
	Fmt        *FmtResolverSpec        `json:"fmt,omitempty"`
	Plain      *PlainResolverSpec      `json:"plain,omitempty"`
	GoTemplate *GotemplateResolverSpec `json:"gotemplate,omitempty"`
}

// APIVersion implements types.DefinedResource.
func (r ResolverSpec) APIVersion() types.APIVersion {
	return APIVersion
}

// Kind implements types.DefinedResource.
func (r ResolverSpec) Kind() types.Kind {
	return ResolverKind
}

// Render implements types.Renderer.
func (r ResolverSpec) Render(types.APIServer) (string, error) {
	panic("unimplemented")
}

func (r ResolverSpec) Resolve(key string, value string) (string, error) {
	if err := validateResolverSpec(r); err != nil {
		return "", err
	}

	switch r.Type {
	case ExecResolverType:
		return r.Exec.Resolve(key, value)
	case FmtResolverType:
		return r.Fmt.Resolve(key, value)
	case PlainResolverType:
		return r.Plain.Resolve(key, value)
	case GotemplateResolverType:
		return r.GoTemplate.Resolve(key, value)
	default:
		return "", flaterrors.Join(
			types.ErrType,
			fmt.Errorf("Resolver type %q is not supported", r.Type),
		)
	}
}

type (
	ExecResolverSpec struct {
		Command string   `json:"command"`
		Args    []string `json:"args,omitempty"`

		// WARN: What? Was my intention to pipe an input file?
		Stdin string `json:"stdin,omitempty"`
	}

	FmtResolverSpec struct {
		Template string `json:"template"`
		// FmtArguments is a list of FmtArgument, that will be used to format the template
		FmtArguments []FmtArgument `json:"fmtArguments"`
	}

	GotemplateResolverSpec struct {
		Template string `json:"template"`
	}

	PlainResolverSpec bool
)

func (r ExecResolverSpec) Resolve(key, value string) (string, error) {
	// TODO: implement me
	panic("ResolveExec is not yet supported")

	cmd := exec.Command(r.Command, append(r.Args, key, value)...)
	cmd.Stdin = strings.NewReader(r.Stdin)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (r FmtResolverSpec) Resolve(key string, value string) (string, error) {
	args := make([]any, 0)
	for _, fmtArg := range r.FmtArguments {
		var arg string
		if fmtArg == KeyFmtArgument {
			arg = key
		} else {
			arg = value
		}
		args = append(args, arg)
	}
	return fmt.Sprintf(r.Template, args...), nil
}

func (r GotemplateResolverSpec) Resolve(key string, value string) (string, error) {
	// TODO: implement me
	panic("unimplemented")
}

func (r PlainResolverSpec) Resolve(key string, value string) (string, error) {
	return key, nil
}

func NewAVKResolver(
	name, namespace string,
	spec ResolverSpec,
) (types.Resource[types.APIVersionKind], error) {
	if err := validateResolverSpec(spec); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	return types.Resource[types.APIVersionKind]{
			APIVersion: APIVersion,
			Kind:       ResolverKind,
			Metadata:   types.NewMetadata(name, namespace),
			Spec:       spec,
		},
		nil
}

func validateResolverSpec(spec ResolverSpec) error {
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
		return fmt.Errorf("cannot parse ResolverSpec.Type; got: %s", spec.Type)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// GetDefaultResolver
//
// This list of resolver is used to populate the ~/.config/vib/resources directory on init.
//----------------------------------------------------------------------------------------------------------------------

func DefaultAVKResolver() []types.Resource[types.APIVersionKind] {
	return []types.Resource[types.APIVersionKind]{
		NewPlainResolver(),
		NewFunctionResolver(),
		NewAliasResolver(),
		NewEnvironmentResolver(),
		NewExportedEnvironmentResolver(),
	}
}

//----------------------------------------------------------------------------------------------------------------------
// PlainResolver
//----------------------------------------------------------------------------------------------------------------------

func NewPlainResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			PlainResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type:  PlainResolverRef,
				Plain: util.Ptr(PlainResolverSpec(true)),
			},
		),
	)
}

//----------------------------------------------------------------------------------------------------------------------
// FunctionResolver
//----------------------------------------------------------------------------------------------------------------------

func NewFunctionResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			FunctionResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type: FmtResolverType,
				Fmt: &FmtResolverSpec{
					Template:     "%s() { %s ; }",
					FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
				},
			},
		),
	)
}

//----------------------------------------------------------------------------------------------------------------------
// AliasResolver
//----------------------------------------------------------------------------------------------------------------------

func NewAliasResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			AliasResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type: FmtResolverType,
				Fmt: &FmtResolverSpec{
					Template:     "alias %s='%s'",
					FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
				},
			},
		),
	)
}

//----------------------------------------------------------------------------------------------------------------------
// EnvironmentResolver
//----------------------------------------------------------------------------------------------------------------------

func NewEnvironmentResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			EnvironmentResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type: FmtResolverType,
				Fmt: &FmtResolverSpec{
					Template:     "%s=%q",
					FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
				},
			},
		),
	)
}

//----------------------------------------------------------------------------------------------------------------------
// ExportedEnvironmentResolver
//----------------------------------------------------------------------------------------------------------------------

func NewExportedEnvironmentResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			ExportedEnvironmentResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type: FmtResolverType,
				Fmt: &FmtResolverSpec{
					Template:     "export %s=%q",
					FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
				},
			},
		),
	)
}
