package main

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

func APIServer(apiKinds []vib.APIKind) (vib.APIServer, error) {
	apiServer := vib.NewAPIServer()

	for _, apiKind := range apiKinds {
		err := apiServer.Register(apiKind)
		if err != nil {
			return nil, err
		}
	}

	return apiServer, nil
}

func APIKinds(config *ConfigSpec) ([]vib.APIKind, error) {
	results := make([]vib.APIKind, 0)

	factory, err := operatorFactory(config)
	if err != nil {
		return nil, err
	}

	for _, x := range []struct {
		APIVersion vib.APIVersion
		Kind       vib.Kind
	}{
		{vib.V1Alpha1, vib.ProfileKind},
		{vib.V1Alpha1, vib.SetKind},
		{vib.V1Alpha1, vib.ExpressionKind},
		{vib.V1Alpha1, vib.ResolverKind},
	} {
		operator, err := factory(x.APIVersion, x.Kind)
		if err != nil {
			return nil, err
		}
		results = append(results, vib.NewAPIKind(x.APIVersion, x.Kind, operator))
	}

	return results, nil
}

func operatorFactory(config *ConfigSpec) (func(vib.APIVersion, vib.Kind) (vib.Operator, error), error) {
	options := make([]interface{}, 0)

	switch config.OperatorStrategy {
	case vib.FileSystemOperatorStrategy:
		options = append(options, config.ResourceDir, vib.YAMLEncoding)
	case vib.GitOperatorStrategy:
		// TODO implement me
		panic("not implemented yet")
	default:
		err := fmt.Errorf("%w: operator strategy %q is not supported", vib.ErrType, config.OperatorStrategy)
		logger.Error(err)
		return nil, err
	}

	return func(apiVersion vib.APIVersion, kind vib.Kind) (vib.Operator, error) {
		options := append([]interface{}{apiVersion, kind}, options...)
		return vib.NewOperator(config.OperatorStrategy, options...)
	}, nil
}

func fastInit() (vib.APIServer, error) {
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
