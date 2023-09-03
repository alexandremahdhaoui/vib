package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"os/exec"
	"strings"
)

// Render
// TODO: Figure out how to use the key and values when rendering an Exec command.
func ResolveExec(resolver *v1alpha1.ResolverSpec, key, value string) (string, error) {
	cmd := exec.Command(resolver.Exec.Command, resolver.Exec.Args...)
	cmd.Stdin = strings.NewReader(resolver.Exec.Stdin)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	panic("use of keys & values is not implemented yet")
	return string(out), nil
}

func ResolveFmt(resolver *v1alpha1.ResolverSpec, key, value string) (string, error) {
	args := make([]any, 0)
	for _, fmtArg := range resolver.Fmt.FmtArguments {
		var arg string
		if fmtArg == v1alpha1.KeyFmtArgument {
			arg = key
		} else {
			arg = value
		}
		args = append(args, arg)
	}
	return fmt.Sprintf(resolver.Fmt.Template, args...), nil
}

// ResolvePlain returns the key
func ResolvePlain(_ *v1alpha1.ResolverSpec, key, _ string) (string, error) {
	return key, nil
}

func ResolveGotemplate(resolver *v1alpha1.ResolverSpec, key, value string) (string, error) {
	// TODO: Implement Me!
	panic("not implemented yet")
}

//----------------------------------------------------------------------------------------------------------------------
// DefaultResolver
//----------------------------------------------------------------------------------------------------------------------

func DefaultResolver(apiVersion api.APIVersion) ([]*api.ResourceDefinition, error) {
	switch apiVersion {
	case apis.V1Alpha1:
	default:

	}

	results := make([]*api.ResourceDefinition, 0)
	for _, f := range []func() (*api.ResourceDefinition, error){
		v1alpha1.PlainResolver,
		v1alpha1.FunctionResolver,
		v1alpha1.AliasResolver,
		v1alpha1.EnvironmentResolver,
		v1alpha1.ExportedEnvironmentResolver,
	} {
		resource, err := f()
		if err != nil {
			return nil, err
		}

		results = append(results, resource)
	}
	return results, nil
}
