package main

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

func APIServer(apiKinds []api.APIKind) (api.APIServer, error) {
	apiServer := api.NewAPIServer()

	for _, apiKind := range apiKinds {
		err := apiServer.Register(apiKind)
		if err != nil {
			return nil, err
		}
	}

	return apiServer, nil
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
		{v1alpha1.APIVersion, v1alpha1.ProfileKind},
		{v1alpha1.APIVersion, v1alpha1.SetKind},
		{v1alpha1.APIVersion, v1alpha1.ExpressionKind},
		{v1alpha1.APIVersion, v1alpha1.ResolverKind},
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
