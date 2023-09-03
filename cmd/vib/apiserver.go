package main

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib"
	"github.com/alexandremahdhaoui/vib/apis"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

func APIServer(apiKinds []api.APIKind) (api.APIServer, error) {
	server := api.NewAPIServer()

	for _, apiKind := range apiKinds {
		err := server.Register(apiKind)
		if err != nil {
			return nil, err
		}
	}

	// get Default Resources
	resources, err := vib.DefaultResolver(apis.V1Alpha1)
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

func APIKinds(config *ConfigSpec) ([]api.APIKind, error) {
	results := make([]api.APIKind, 0)

	factory, err := operatorFactory(config)
	if err != nil {
		return nil, err
	}

	for _, x := range []struct {
		APIVersion api.APIVersion
		Kind       api.Kind
	}{
		{apis.V1Alpha1, apis.ProfileKind},
		{apis.V1Alpha1, apis.SetKind},
		{apis.V1Alpha1, apis.ExpressionKind},
		{apis.V1Alpha1, apis.ExpressionSetKind},
		{apis.V1Alpha1, apis.ResolverKind},
	} {
		operator, err := factory(x.APIVersion, x.Kind)
		if err != nil {
			return nil, err
		}
		results = append(results, api.NewAPIKind(x.APIVersion, x.Kind, operator))
	}

	return results, nil
}

func operatorFactory(config *ConfigSpec) (func(api.APIVersion, api.Kind) (api.Operator, error), error) {
	options := make([]interface{}, 0)

	switch config.OperatorStrategy {
	case api.FileSystemOperatorStrategy:
		options = append(options, config.ResourceDir, api.YAMLEncoding)
	case api.GitOperatorStrategy:
		// TODO implement me
		panic("not implemented yet")
	default:
		err := fmt.Errorf("%w: operator strategy %q is not supported", logger.ErrType, config.OperatorStrategy)
		logger.Error(err)
		return nil, err
	}

	return func(apiVersion api.APIVersion, kind api.Kind) (api.Operator, error) {
		options := append([]interface{}{apiVersion, kind}, options...)
		return api.NewOperator(config.OperatorStrategy, options...)
	}, nil
}

func fastInit() (api.APIServer, error) {
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
