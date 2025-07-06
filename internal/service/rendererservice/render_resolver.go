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

package rendererservice

import (
	"errors"
	"fmt"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/mitchellh/mapstructure"
)

func Resolve(resource *types.Resource, key, value string) (string, error) {
	switch resource.APIVersion {
	case apis.V1Alpha1:
		return dispatchV1Alpha1Resolver(resource, key, value)
	default:
		return "", errors.Join(
			types.ErrType,
			fmt.Errorf("APIVersion %q is not supported", resource.APIVersion),
		)
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
		return "", flaterrors.Join(
			types.ErrType,
			fmt.Errorf("Resolver type %q is not supported", resolver.Type),
		)
	}
}
