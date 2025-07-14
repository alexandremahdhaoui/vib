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

package main

import (
	"fmt"

	"github.com/alexandremahdhaoui/vib/internal/adapter/resource"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/apis/v1alpha1"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

func NewAPIServer(apiKinds []service.APIKind) (service.APIServer, error) {
	server := service.NewAPIServer()

	for _, apiKind := range apiKinds {
		err := server.Register(apiKind)
		if err != nil {
			return nil, err
		}
	}

	// get Default Resources
	resources, err := resolveradapter.DefaultResolver(apis.V1Alpha1)
	if err != nil {
		return nil, err
	}

	// install the default resources
	for _, resource := range resources {
		err = server.Update(&resource.APIVersion, resource.Kind, resource.Metadata.Name, resource)
		if err != nil {
			return nil, err
		}
	}

	return server, nil
}

func APIKinds(config *ConfigSpec) ([]service.APIKind, error) {
	results := make([]service.APIKind, 0)

	factory, err := getStorage(config)
	if err != nil {
		return nil, err
	}

	for _, x := range []struct {
		APIVersion types.APIVersion
		Kind       types.Kind
	}{
		{v1alpha1.APIVersion, v1alpha1.ProfileKind},
		{v1alpha1.APIVersion, v1alpha1.SetKind},
		{v1alpha1.APIVersion, v1alpha1.ExpressionKind},
		{v1alpha1.APIVersion, v1alpha1.ExpressionSetKind},
		{v1alpha1.APIVersion, v1alpha1.ResolverKind},
	} {
		// WARN: this is terrible
		storage, err := factory(x.APIVersion, x.Kind)
		if err != nil {
			return nil, err
		}
		results = append(results, service.NewAPIKind(x.APIVersion, x.Kind, storage))
	}

	return results, nil
}

// WARN: This is really bad
func getStorage(
	config *ConfigSpec,
) (func(service.APIVersion, service.Kind) (service.Operator, error), error) {
	options := make([]interface{}, 0)

	switch config.StorageStrategy {
	case resourceadapter.FileSystemOperatorStrategy:
		options = append(options, config.ResourceDir, service.YAMLEncoding)
	case resourceadapter.GitOperatorStrategy:
		// TODO: implement me
		panic("not implemented yet")
	default:
		return nil, flaterrors.Join(
			fmt.Errorf("storage strategy %q is not supported", config.StorageStrategy),
			types.ErrType,
		)
	}

	return func(apiVersion service.APIVersion, kind service.Kind) (resourceadapter.Storage, error) {
		options := append([]interface{}{apiVersion, kind}, options...)
		return resourceadapter.New(config.StorageStrategy, options...)
	}, nil
}

func fastInit() (service.APIServer, error) {
	config, err := readConfig(nil)
	if err != nil {
		return nil, err
	}

	apiKinds, err := APIKinds(config)
	if err != nil {
		return nil, err
	}

	return APIServer(apiKinds)
}
