package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
)

func Resolve(resource *api.ResourceDefinition, key, value string) (string, error) {
	switch resource.APIVersion {
	case apis.V1Alpha1:
		return dispatchV1Alpha1Resolver(resource, key, value)
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("APIVersion %q is not supported", resource.APIVersion))
	}
}

func dispatchV1Alpha1Resolver(resource *api.ResourceDefinition, key, value string) (string, error) {
	resolver := new(v1alpha1.ResolverSpec)
	err := mapstructure.Decode(resource.Spec, resolver)
	if err != nil {
		return "", err
	}

	switch resolver.Type {
	case v1alpha1.ExecResolverType:
		return ResolveExec(resolver, key, value)
	case v1alpha1.FmtResolverType:
		return ResolveFmt(resolver, key, value)
	case v1alpha1.PlainResolverType:
		return ResolvePlain(resolver, key, value)
	case v1alpha1.GotemplateResolverType:
		return ResolveGotemplate(resolver, key, value)
	default:
		return "", logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("Resolver type %q is not supported", resolver.Type))
	}
}
