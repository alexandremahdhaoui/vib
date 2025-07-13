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

package resourceadapter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

//----------------------------------------------------------------------------------------------------------------------
// FilesystemStorage
//----------------------------------------------------------------------------------------------------------------------

// NewFilesystem instantiate a new strategy
func NewFilesystem[T types.APIVersionKind](
	resourceDir string,
	codec types.Codec,
) (types.Storage[T], error) {
	// ensure the resourceDir exists
	if err := os.MkdirAll(resourceDir, 0777); err != nil {
		return nil, err
	}

	v := *new(T)
	return &filesystem[T]{
		apiVersion:  v.APIVersion(),
		kind:        v.Kind(),
		resourceDir: resourceDir,
		codec:       codec,
	}, nil
}

// filesystem operates T through the filesystem.
// Resources are stored on the filesystem using the following convention:
// - Filename: {{ T.APIVersion() }}.{{ T.Kind() }}.{{ T.IMetadata().Name }}. {{ s.encoder.Encoding() }}
type filesystem[T types.APIVersionKind] struct {
	apiVersion  types.APIVersion
	kind        types.Kind
	resourceDir string
	codec       types.Codec
}

func (fs *filesystem[T]) List() ([]types.Resource[T], error) {
	res := make([]types.Resource[T], 0)
	// Get all instance of T.
	dentries, err := os.ReadDir(fs.resourceDir)
	if err != nil {
		return nil, err
	}

	r := strings.ToLower(fmt.Sprintf(
		`%s\.%s\..*\.%s`,
		cleanAPIVersionForFilesystem(fs.apiVersion),
		fs.kind,
		fs.codec.Encoding(),
	))

	regex, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}

	for _, dentry := range dentries {
		if dentry.IsDir() || !regex.MatchString(dentry.Name()) {
			continue
		}

		v, err := fs.read(fs.joinDirAndBasename(dentry.Name()))
		if err != nil {
			return nil, err
		}

		res = append(res, v)

	}

	return res, nil
}

func (fs *filesystem[T]) Get(name string) (types.Resource[T], error) {
	v, err := fs.read(fs.filepathFromResourceName(name))
	if os.IsNotExist(err) {
		return types.Resource[T]{}, flaterrors.Join(err, types.ErrNotFound)
	} else if err != nil {
		return types.Resource[T]{}, err
	}

	return v, nil
}

// Create should create only if file does not already exist.
func (fs *filesystem[T]) Create(res types.Resource[T]) error {
	exist, err := resourceExist[T](fs, res.Metadata.Name)
	if err != nil {
		return err
	}

	if exist {
		return flaterrors.Join(
			fmt.Errorf(
				"apiVersion: %q, kind: %q, name: %q",
				types.ErrExist,
				res.APIVersion,
				res.Kind,
				res.Metadata.Name,
			),
			errors.New("cannot create resource"),
			types.ErrExist,
		)
	}

	return fs.writeAtomic(res)
}

func (fs *filesystem[T]) writeAtomic(v types.Resource[T]) error {
	f, err := os.CreateTemp("", "vibtmp_*")
	if err != nil {
		return err
	}

	b, err := fs.codec.Marshal(v)
	if err != nil {
		return err
	}

	err = os.WriteFile(f.Name(), b, 0640)
	if err != nil {
		return err
	}

	dest := fs.filepathFromResourceName(v.Metadata.Name)
	if err := os.Rename(f.Name(), dest); err != nil {
		return err
	}

	return nil
}

func (fs *filesystem[T]) Update(oldName string, v types.Resource[T]) error {
	// Rename if required
	if oldName != v.Metadata.Name {
		oldPath := fs.filepathFromResourceName(oldName)
		newPath := fs.filepathFromResourceName(v.Metadata.Name)
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return fs.writeAtomic(v)
}

func (fs *filesystem[T]) Delete(name string) error {
	return os.Remove(fs.filepathFromResourceName(name))
}

// read tries to read file corresponding to the specified object's name.
// It returns os.ErrNotExist if file does not exist.
func (fs *filesystem[T]) read(path string) (types.Resource[T], error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return types.Resource[T]{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return types.Resource[T]{}, err
	}

	v := new(types.Resource[T])
	if err := fs.codec.Unmarshal(b, v); err != nil {
		return types.Resource[T]{}, err
	}
	return *v, nil
}

// joinDirAndBasename joins the resourceDir to the basename
func (fs *filesystem[T]) joinDirAndBasename(basename string) string {
	return filepath.Join(fs.resourceDir, basename)
}

// filepathFromResourceName computes the resourceDir to the corresponding resource name, based on the naming convention
func (fs *filesystem[T]) filepathFromResourceName(resourceName string) string {
	return fs.joinDirAndBasename(fs.basename(resourceName))
}

// filepathFromResourceName computes the resource filename, based on the naming convention
func (fs *filesystem[T]) basename(resourceName string) string {
	return strings.ToLower(fmt.Sprintf(
		"%s.%s.%s.%s",
		cleanAPIVersionForFilesystem(fs.apiVersion),
		fs.kind,
		resourceName,
		fs.codec.Encoding(),
	))
}

//----------------------------------------------------------------------------------------------------------------------
// GitStorage
//----------------------------------------------------------------------------------------------------------------------

// GitStorage uses FilesystemStrategy as a backend, and leverages Git for version control.
type GitStorage[T types.APIVersionKind] struct {
	innerStrategy filesystem[T]
}

//----------------------------------------------------------------------------------------------------------------------
// Storage Utils
//----------------------------------------------------------------------------------------------------------------------

// resourceExist checks if a named resource already exist
func resourceExist[T types.APIVersionKind](storage types.Storage[T], name string) (bool, error) {
	_, err := storage.Get(name)
	if errors.Is(err, types.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// cleanAPIVersionForFilesystem transforms `vib/v1alpha1` into `vib_v1alpha1`
func cleanAPIVersionForFilesystem(s types.APIVersion) string {
	return strings.ReplaceAll(string(s), "/", "_")
}
