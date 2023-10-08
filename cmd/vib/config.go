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
	"github.com/alexandremahdhaoui/vib"
	"github.com/alexandremahdhaoui/vib/apis"
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

	return api.NewResourceDefinition(apis.V1Alpha1, configKind, configName, ConfigSpec{
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
	strategy, err := api.NewFilesystemOperator(apis.V1Alpha1, configKind, cfgDir, api.YAMLEncoding)
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
