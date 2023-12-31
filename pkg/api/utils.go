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

package api

import (
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"os"
	"path/filepath"
)

func fileExist(path string) (bool, error) {
	err := mkBaseDir(path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			logger.Error(err)
			return false, err
		}
	}

	return true, nil
}

func mkBaseDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func ToPointer[T any](t T) *T {
	value := t

	return &value
}
