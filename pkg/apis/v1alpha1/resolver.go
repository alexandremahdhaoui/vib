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

// Resolver is the interface that all resolvers must implement.
type Resolver interface {
	// Resolve takes a key and a value and returns the resolved string.
	Resolve(key, value string) (string, error)
}

// FmtArgument is a string that represents a format argument.
type FmtArgument string

const (
	// ExecResolverType is the type for exec resolvers.
	ExecResolverType = "exec"
	// FmtResolverType is the type for fmt resolvers.
	FmtResolverType = "fmt"
	// PlainResolverType is the type for plain resolvers.
	PlainResolverType = "plain"
	// GotemplateResolverType is the type for go-template resolvers.
	GotemplateResolverType = "gotemplate"

	// PlainResolverRef is the name of the plain resolver.
	PlainResolverRef = "plain"
	// FunctionResolverRef is the name of the function resolver.
	FunctionResolverRef = "function"
	// AliasResolverRef is the name of the alias resolver.
	AliasResolverRef = "alias"
	// EnvironmentResolverRef is the name of the environment resolver.
	EnvironmentResolverRef = "environment"
	// ExportedEnvironmentResolverRef is the name of the exported environment resolver.
	ExportedEnvironmentResolverRef = "environment-exported"

	// KeyFmtArgument is the key format argument.
	KeyFmtArgument FmtArgument = "key"
	// ValueFmtArgument is the value format argument.
	ValueFmtArgument FmtArgument = "value"
)

// ResolverSpec defines the desired state of a Resolver.
// It specifies the type of the resolver and its configuration.
type ResolverSpec struct {
	// Type is the type of the resolver.
	Type string `json:"type"`
	// Exec is the configuration for an exec resolver.
	Exec *ExecResolverSpec `json:"exec,omitempty"`
	// Fmt is the configuration for a fmt resolver.
	Fmt *FmtResolverSpec `json:"fmt,omitempty"`
	// Plain is the configuration for a plain resolver.
	Plain *PlainResolverSpec `json:"plain,omitempty"`
	// GoTemplate is the configuration for a go-template resolver.
	GoTemplate *GotemplateResolverSpec `json:"gotemplate,omitempty"`
}

// APIVersion returns the APIVersion of the ResolverSpec.
// It implements the types.DefinedResource interface.
func (r ResolverSpec) APIVersion() types.APIVersion {
	return APIVersion
}

// Kind returns the Kind of the ResolverSpec.
// It implements the types.DefinedResource interface.
func (r ResolverSpec) Kind() types.Kind {
	return ResolverKind
}

// Render implements types.Renderer.
// It is not implemented for ResolverSpec and will panic if called.
func (r ResolverSpec) Render(types.APIServer) (string, error) {
	panic("unimplemented")
}

// Resolve resolves the given key-value pair using the appropriate resolver.
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
	// ExecResolverSpec defines the configuration for an exec resolver.
	ExecResolverSpec struct {
		// Command is the command to execute.
		Command string `json:"command"`
		// Args is a list of arguments to pass to the command.
		Args []string `json:"args,omitempty"`
		// Stdin is a string to be piped to the command's stdin.
		// WARN: What? Was my intention to pipe an input file?
		Stdin string `json:"stdin,omitempty"`
	}

	// FmtResolverSpec defines the configuration for a fmt resolver.
	FmtResolverSpec struct {
		// Template is the fmt template string.
		Template string `json:"template"`
		// FmtArguments is a list of FmtArgument, that will be used to format the template.
		FmtArguments []FmtArgument `json:"fmtArguments"`
	}

	// GotemplateResolverSpec defines the configuration for a go-template resolver.
	GotemplateResolverSpec struct {
		// Template is the go-template string.
		Template string `json:"template"`
	}

	// PlainResolverSpec defines the configuration for a plain resolver.
	PlainResolverSpec bool
)

// Resolve executes the command and returns the output.
// It is not yet implemented and will panic if called.
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

// Resolve formats the template with the given key and value.
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

// Resolve executes the go-template with the given key and value.
// It is not yet implemented and will panic if called.
func (r GotemplateResolverSpec) Resolve(key string, value string) (string, error) {
	// TODO: implement me
	panic("unimplemented")
}

// Resolve returns the key.
func (r PlainResolverSpec) Resolve(key string, value string) (string, error) {
	return key, nil
}

// NewAVKResolver creates a new resolver resource.
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

// DefaultAVKResolver returns a list of default resolver resources.
// This list of resolver is used to populate the ~/.config/vib/resources directory on init.
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

// NewPlainResolver creates a new plain resolver resource.
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

// NewFunctionResolver creates a new function resolver resource.
func NewFunctionResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			FunctionResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type: FmtResolverType,
				Fmt: &FmtResolverSpec{
					Template:     "function %s() {\n%s\n}",
					FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
				},
			},
		),
	)
}

//----------------------------------------------------------------------------------------------------------------------
// AliasResolver
//----------------------------------------------------------------------------------------------------------------------

// NewAliasResolver creates a new alias resolver resource.
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

// NewEnvironmentResolver creates a new environment resolver resource.
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

// NewExportedEnvironmentResolver creates a new exported environment resolver resource.
func NewExportedEnvironmentResolver() types.Resource[types.APIVersionKind] {
	return util.Must(
		NewAVKResolver(
			ExportedEnvironmentResolverRef,
			types.VibSystemNamespace,
			ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
				Type: FmtResolverType,
				Fmt: &FmtResolverSpec{
					Template:     "%s=%q\nexport %s",
					FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument, KeyFmtArgument},
				},
			},
		),
	)
}
