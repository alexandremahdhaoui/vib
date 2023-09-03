package v1alpha1

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"os/exec"
	"strings"
)

type FmtArgument string

const (
	ResolverKind           = "Resolver"
	ExecResolverType       = "exec"
	FmtResolverType        = "fmt"
	PlainResolverType      = "plain"
	GotemplateResolverType = "gotemplate"

	FunctionResolverRef            = "function"
	AliasResolverRef               = "alias"
	EnvironmentResolverRef         = "environment"
	ExportedEnvironmentResolverRef = "exportedEnvironment"

	KeyFmtArgument   FmtArgument = "key"
	ValueFmtArgument FmtArgument = "value"
)

type ResolverSpec struct {
	Type       string              `json:"type"                 yaml:"type"`
	Exec       *ExecResolver       `json:"exec,omitempty"       yaml:"exec,omitempty"`
	Fmt        *FmtResolver        `json:"fmt,omitempty"        yaml:"fmt,omitempty"`
	Plain      *PlainResolver      `json:"plain,omitempty"      yaml:"plain,omitempty"`
	GoTemplate *GotemplateResolver `json:"gotemplate,omitempty" yaml:"gotemplate,omitempty"`
}

func (r *ResolverSpec) Render(key, value string) (string, error) {
	// Dispatcher
	switch r.Type {
	case ExecResolverType:
		return r.Exec.Render(key, value)
	case FmtResolverType:
		return r.Fmt.Render(key, value)
	case PlainResolverType:
		return r.Plain.Render(key, value)
	case GotemplateResolverType:
		return r.Plain.Render(key, value)
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("Resolver type %q is not supported", r.Type))
	}
}

type ExecResolver struct {
	Command string   `json:"command" yaml:"command"`
	Args    []string `json:"args,omitempty" yaml:"args,omitempty"`

	// Stdin is a format-able
	Stdin string `json:"stdin,omitempty" yaml:"stdin,omitempty"`
}

// Render
// TODO: Figure out how to use the key and values when rendering an Exec command.
func (r *ExecResolver) Render(key, value string) (string, error) {
	cmd := exec.Command(r.Command, r.Args...)
	cmd.Stdin = strings.NewReader(r.Stdin)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

type FmtResolver struct {
	Template string
	// FmtArguments is a list of FmtArgument, that will be used to format the template
	FmtArguments []FmtArgument
}

func (r *FmtResolver) Render(key, value string) (string, error) {
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

type PlainResolver bool

func (r *PlainResolver) Render(_, value string) (string, error) {
	return value, nil
}

type GotemplateResolver struct {
	Template string
}

func (r *GotemplateResolver) Render(key, value string) (string, error) {
	// TODO: Implement Me!
	panic("not implemented yet")
}

func NewResolver(name string, spec ResolverSpec) (*api.ResourceDefinition, error) {
	if err := validateResolverSpec(&spec); err != nil {
		return nil, err
	}

	return &api.ResourceDefinition{
			APIVersion: APIVersion,
			Kind:       ResolverKind,
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
// FunctionResolver
//----------------------------------------------------------------------------------------------------------------------

func FunctionResolver() (*api.ResourceDefinition, error) {
	return NewResolver(
		FunctionResolverRef,
		ResolverSpec{ //nolint:exhaustruct,exhaustivestruct
			Type: FmtResolverType,
			Fmt: &FmtResolver{
				Template:     "%s() { %s ; }",
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
			Fmt: &FmtResolver{
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
			Fmt: &FmtResolver{
				Template:     "%s=\"%s\"",
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
			Fmt: &FmtResolver{
				Template:     "export %s=\"%s\"",
				FmtArguments: []FmtArgument{KeyFmtArgument, ValueFmtArgument},
			},
		},
	)
}
