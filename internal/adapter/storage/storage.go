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

package storageadapter

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/internal/util"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

type StorageStrategy string

const (
	FileSystemStorageStrategy StorageStrategy = "filesystem"
	GitStorageStrategy        StorageStrategy = "git"
)

type Storage interface {
	Get(name *string) ([]types.Resource, error)
	Create(*types.Resource) error
	Update(name *string, v *types.Resource) error
	Delete(name string) error
}

//----------------------------------------------------------------------------------------------------------------------
// New
//----------------------------------------------------------------------------------------------------------------------

func New(strategy StorageStrategy, options ...any) (Storage, error) {
	switch strategy {
	case FileSystemStorageStrategy:
		if len(options) != 4 {
			fmt.Printf("%#v", options)
			return nil, flaterrors.Join(
				types.ErrType,
				fmt.Errorf(
					"wrong number of argument to construct FilesystemStorage; got: %d",
					len(options),
				),
			)
		}

		apiVersion, ok := options[0].(types.APIVersion)
		if !ok {
			return nil, flaterrors.Join(types.ErrType, "apiVersion must be of type string")
		}

		kind, ok := options[1].(types.Kind)
		if !ok {
			return nil, flaterrors.Join(types.ErrType, "kind must be of type string")
		}

		resourceDir, ok := options[2].(string)
		if !ok {
			return nil, flaterrors.Join(types.ErrType, "resourceDir must be of type string")
		}

		encoding, ok := options[3].(util.Encoding)
		if !ok {
			return nil, flaterrors.Join(types.ErrType, "encoding must be of type Encoding")
		}

		return NewFilesystem(apiVersion, kind, resourceDir, encoding)
	case GitStorageStrategy:
		// TODO implement me
		panic("not implemented yet")
	default:
		err := fmt.Errorf("%w: operator strategy %q is not supported", types.ErrType, strategy)
		return nil, err
	}
}

//----------------------------------------------------------------------------------------------------------------------
// FilesystemStorage
//----------------------------------------------------------------------------------------------------------------------

// FilesystemStorage operates T through the filesystem.
// Resources are stored on the filesystem using the following convention:
// - Filename: {{ T.APIVersion() }}.{{ T.Kind() }}.{{ T.IMetadata().Name }}. {{ s.encoder.Encoding() }}
type FilesystemStorage struct {
	apiVersion  types.APIVersion
	kind        types.Kind
	resourceDir string
	encoder     util.Encoder
}

func (s *FilesystemStorage) Get(name *string) ([]types.Resource, error) {
	res := make([]types.Resource, 0)
	// if name is specified T with specified name exist, then we return a list of length one containing T
	if name != nil {
		// read can return a nil pointer
		v, err := s.readWithResourceName(*name)
		if err != nil {
			return nil, err
		}

		// read can return a nil pointer. If nil pointer, we directly return a nil array
		if v == nil {
			return nil, nil
		}

		// If not nil then we can dereference and add the struct to the list
		res = append(res, *v)

		return res, nil
	}

	// Get all instance of T.
	list, err := s.list()
	if err != nil {
		return nil, err
	}

	for _, filename := range list {
		v, err := s.readWithFilename(filename)
		if err != nil {
			return nil, err
		}

		// read can return a nil pointer. If nil, continue the loop
		if v == nil {
			continue
		}

		// we can safely dereference the pointer
		res = append(res, *v)
	}

	return res, nil
}

// Create should create only if file does not already exist.
func (s *FilesystemStorage) Create(t *types.Resource) error {
	exist, err := resourceExist(s, t.Metadata.Name)
	if err != nil {
		return err
	}

	if exist {
		return flaterrors.Join(
			fmt.Errorf(
				"apiVersion: %q, kind: %q, name: %q",
				types.ErrExist,
				t.APIVersion,
				t.Kind,
				t.Metadata.Name,
			),
			"cannot create resource",
			types.ErrExist,
		)
	}

	return s.write(t)
}

func (s *FilesystemStorage) Update(name *string, v *types.Resource) error {
	// This operation rename the object
	if name != nil {
		name := *name
		if v.Metadata.Name != name {
			// Delete existing object with name `name *string` and rewrite with new name.
			if err := s.Delete(name); err != nil {
				return err
			}
		}
	}

	// Name is the same or former one was deleted, we can write T
	return s.write(v)
}

func (s *FilesystemStorage) Delete(name string) error {
	return os.Remove(s.filepathFromResourceName(name))
}

// list returns a list of reference to T.
func (s *FilesystemStorage) list() ([]string, error) {
	dirEntries, err := os.ReadDir(s.resourceDir)
	if err != nil {
		return nil, err
	}

	r := strings.ToLower(fmt.Sprintf(
		"%s\\.%s\\..*\\.%s",
		cleanAPIVersionForFilesystem(s.apiVersion),
		s.kind,
		s.encoder.Encoding(),
	))

	regex, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}

	refs := make([]string, 0)
	for _, entry := range dirEntries {
		if entry.IsDir() || !regex.MatchString(entry.Name()) {
			continue
		}

		refs = append(refs, entry.Name())
	}

	return refs, nil
}

// read tries to read file corresponding to the specified object's name.
// Returns a pointer to an unmarshalled T and an error.
// Warn: read can return a nil pointer if resource wasn't find
// The reason read returns a possible nil pointer is to avoid raising an error if the file we're trying to read does not
// exist.
func (s *FilesystemStorage) read(path string) (*types.Resource, error) {
	if ok, err := util.FileExist(path); !ok {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	resource, err := util.ReadEncodedFile(path)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// readWithFilename makes a call to read
func (s *FilesystemStorage) readWithFilename(name string) (*types.Resource, error) {
	return s.read(s.filepathFromFilename(name))
}

// readWithResourceName makes a call to read
func (s *FilesystemStorage) readWithResourceName(name string) (*types.Resource, error) {
	return s.read(s.filepathFromResourceName(name))
}

// write tries to write T to filesystem.
func (s *FilesystemStorage) write(resource *types.Resource) error {
	path := s.filepathFromResourceName(resource.Metadata.Name)
	// Ensures the base dir exist before trying to write into it.
	if err := util.MkBaseDir(path); err != nil {
		return err
	}

	return util.WriteEncodedFile(path, resource)
}

// filepathFromResourceName computes the resourceDir to the corresponding resource name, based on the naming convention
func (s *FilesystemStorage) filepathFromFilename(name string) string {
	return filepath.Join(s.resourceDir, name)
}

// filepathFromResourceName computes the resourceDir to the corresponding resource name, based on the naming convention
func (s *FilesystemStorage) filepathFromResourceName(name string) string {
	return s.filepathFromFilename(s.filename(name))
}

// filepathFromResourceName computes the resource filename, based on the naming convention
func (s *FilesystemStorage) filename(name string) string {
	return strings.ToLower(fmt.Sprintf(
		"%s.%s.%s.%s",
		cleanAPIVersionForFilesystem(s.apiVersion),
		s.kind,
		name,
		s.encoder.Encoding(),
	))
}

// NewFilesystem instantiate a new strategy
func NewFilesystem(
	apiVersion types.APIVersion,
	kind types.Kind,
	resourceDir string,
	encoding util.Encoding,
) (*FilesystemStorage, error) { //nolint:lll
	encoder, err := util.NewEncoder(encoding)
	if err != nil {
		return nil, err
	}

	apiVersion, err = apiVersion.Validate()
	if err != nil {
		return nil, err
	}

	kind, err = kind.Validate()
	if err != nil {
		return nil, err
	}

	return &FilesystemStorage{
		apiVersion:  apiVersion,
		kind:        kind,
		resourceDir: resourceDir,
		encoder:     encoder,
	}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// GitStorage
//----------------------------------------------------------------------------------------------------------------------

// GitStorage uses FilesystemStrategy as a backend, and leverages Git for version control.
type GitStorage struct {
	innerStrategy FilesystemStorage
}

//----------------------------------------------------------------------------------------------------------------------
// Storage Utils
//----------------------------------------------------------------------------------------------------------------------

// resourceExist checks if a named resource already exist
func resourceExist(operator Storage, name string) (bool, error) {
	arr, err := operator.Get(util.Ptr(name))
	if err != nil {
		return false, err
	}

	if len(arr) == 0 {
		return false, nil
	}

	return true, nil
}

func defaultStorageStrategy() StorageStrategy {
	return FileSystemStorageStrategy
}

// cleanAPIVersionForFilesystem transforms `vib/v1alpha1` into `vib_v1alpha1`
func cleanAPIVersionForFilesystem(s types.APIVersion) string {
	return strings.ReplaceAll(string(s), "/", "_")
}
