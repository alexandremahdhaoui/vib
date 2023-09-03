package main

import (
	"github.com/alexandremahdhaoui/vib"
	"github.com/alexandremahdhaoui/vib/apis/v1alpha1"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"os"
	"path/filepath"
)

const (
	configName = "config"
	configKind = "Config"

	resourcesPath = "resources"
)

// ConfigSpec stores important information to run the vib command line.
// The config is always stored on disk, thus the Operator for managing Config will always be of type
// vib.FilesystemOperator.
type ConfigSpec struct {
	// OperatorStrategy defines which concrete implementation of vib.Operator should be used
	OperatorStrategy api.OperatorStrategy
	// ResourceDir specifies the absolute path to Resource definitions.
	// Defaults to CONFIG_DIR/vib/resources
	ResourceDir string
}

func defaultConfig() (*api.ResourceDefinition, error) {
	resourceDir, err := defaultResourceDir()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return api.NewResourceDefinition(v1alpha1.APIVersion, configKind, configName, ConfigSpec{
		OperatorStrategy: defaultOperatorStrategy(),
		ResourceDir:      resourceDir,
	}), nil
}

// readConfig uses a pointer to a string in order to simplify testing
func readConfig(configDir *string) (*ConfigSpec, error) {
	var err error
	var cfgDir string
	var resource *api.ResourceDefinition

	if configDir != nil {
		cfgDir = *configDir
	} else {
		cfgDir, err = vibConfigDir()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	// Initiate FS strategy for reading the config
	strategy, err := api.NewFilesystemOperator(v1alpha1.APIVersion, configKind, cfgDir, api.YAMLEncoding)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	resources, err := strategy.Get(vib.ToPointer(configName))
	if err != nil {
		return nil, err
	}

	// check if there is no existing resources
	if len(resources) == 0 {
		resource, err = defaultConfig()
		if err != nil {
			return nil, err
		}

		// create a new Config and return it
		err = strategy.Create(resource)
		if err != nil {
			return nil, err
		}
	} else {
		resource = &resources[0]
	}

	config := new(ConfigSpec)
	if err = mapstructure.Decode(resource.Spec, config); err != nil {
		return nil, err
	}

	return config, nil
}

func vibConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, cliName), nil
}

func defaultResourceDir() (string, error) {
	path, err := vibConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, resourcesPath), nil
}

func defaultOperatorStrategy() api.OperatorStrategy {
	return api.FileSystemOperatorStrategy
}
