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
	"os"
	"path/filepath"

	storageadapter "github.com/alexandremahdhaoui/vib/internal/adapter/storage"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/internal/util"
)

const (
	configName = "config"
	configKind = "Config"

	resourcesPath = "resources"
)

const DefaultStorageStrategy = storageadapter.FileSystemStorageStrategy

// ConfigSpec stores important information to run the vib command line.
// The config is always stored on disk, thus the Operator for managing Config will always be of type
// vib.FilesystemOperator.
type ConfigSpec struct {
	// StorageStrategy defines which storage strategy must be used (only filesystem is supported).
	StorageStrategy storageadapter.StorageStrategy
	// ResourceDir specifies the absolute path to Resource definitions.
	// Defaults to CONFIG_DIR/vib/resources
	ResourceDir string
}

func defaultConfig() (*types.Resource, error) {
	resourceDir, err := defaultResourceDir()
	if err != nil {
		return nil, err
	}

	return types.NewResource(apis.V1Alpha1, configKind, configName, ConfigSpec{
		StorageStrategy: DefaultStorageStrategy,
		ResourceDir:     resourceDir,
	}), nil
}

// readConfig uses a pointer to a string in order to simplify testing
func readConfig(configDir *string) (*ConfigSpec, error) {
	var err error
	var cfgDir string
	var resource *types.Resource

	if configDir != nil {
		cfgDir = *configDir
	} else {
		cfgDir, err = vibConfigDir()
		if err != nil {
			return nil, err
		}
	}

	// Initiate FS strategy for reading the config
	strategy, err := storageadapter.NewFilesystem(
		apis.V1Alpha1,
		configKind,
		cfgDir,
		util.YAMLEncoding,
	)
	if err != nil {
		return nil, err
	}

	resources, err := strategy.Get(util.Ptr(configName))
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
