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
func NewFilesystem(
	apiServer types.APIServer,
	codec types.Codec,
	resourceDir string,
) (types.Storage, error) {
	// ensure the resourceDir exists
	if err := os.MkdirAll(resourceDir, 0777); err != nil {
		return nil, err
	}

	return &filesystem{
		resourceDir: resourceDir,
		codec:       codec,
		apiServer:   apiServer,
	}, nil
}

// filesystem operates T through the filesystem.
// Resources are stored on the filesystem using the following convention:
// - Filename: {{ T.APIVersion() }}.{{ T.Kind() }}.{{ T.IMetadata().Name }}. {{ s.encoder.Encoding() }}
type filesystem struct {
	resourceDir string
	codec       types.Codec
	apiServer   types.APIServer
}

var errAPIVersionMustBeSpecified = errors.New("apiVersion must be specified")

func (fs *filesystem) List(
	avk types.APIVersionKind,
) ([]types.Resource[types.APIVersionKind], error) {
	if avk.APIVersion() == "" {
		return nil, errAPIVersionMustBeSpecified
	}

	res := make([]types.Resource[types.APIVersionKind], 0)
	// Get all instance of T.
	dentries, err := os.ReadDir(fs.resourceDir)
	if err != nil {
		return nil, err
	}

	r := strings.ToLower(fmt.Sprintf(
		`%s\.%s\..*\.%s`,
		cleanAPIVersionForFilesystem(avk.APIVersion()),
		avk.Kind(),
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

		v, err := fs.read(fs.joinResourceDirAndBasename(dentry.Name()))
		if err != nil {
			return nil, err
		}

		res = append(res, v)

	}

	return res, nil
}

func (fs *filesystem) Get(
	avk types.APIVersionKind,
	name string,
) (types.Resource[types.APIVersionKind], error) {
	if avk.APIVersion() == "" {
		return types.Resource[types.APIVersionKind]{}, errAPIVersionMustBeSpecified
	}

	v, err := fs.read(fs.filepathFromResourceName(avk, name))
	if os.IsNotExist(err) {
		return types.Resource[types.APIVersionKind]{}, flaterrors.Join(err, types.ErrNotFound)
	} else if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	return v, nil
}

// Create should create only if file does not already exist.
func (fs *filesystem) Create(res types.Resource[types.APIVersionKind]) error {
	exist, err := resourceExist(fs, res.Spec, res.Metadata.Name)
	if err != nil {
		return err
	}

	if exist {
		return flaterrors.Join(
			types.ErrExist,
			fmt.Errorf(
				"apiVersion: %q, kind: %q, name: %q",
				res.APIVersion,
				res.Kind,
				res.Metadata.Name,
			),
			errors.New("cannot create resource"),
		)
	}

	return fs.writeAtomic(res)
}

func (fs *filesystem) writeAtomic(v types.Resource[types.APIVersionKind]) error {
	tmpDest := fs.tmpFilepathFromResourceName(v.Spec, v.Metadata.Name)
	dest := fs.filepathFromResourceName(v.Spec, v.Metadata.Name)

	b, err := fs.codec.Marshal(v)
	if err != nil {
		return err
	}

	if err = os.WriteFile(tmpDest, b, 0640); err != nil {
		return err
	}

	if err := os.Rename(tmpDest, dest); err != nil {
		return err
	}

	return nil
}

func (fs *filesystem) Update(oldName string, v types.Resource[types.APIVersionKind]) error {
	if err := fs.writeAtomic(v); err != nil {
		return err
	}

	// Remove old resource once new one has been created
	if oldName != v.Metadata.Name {
		oldPath := fs.filepathFromResourceName(v.Spec, oldName)
		return os.Remove(oldPath)
	}

	return nil
}

func (fs *filesystem) Delete(avk types.APIVersionKind, name string) error {
	if avk.APIVersion() == "" {
		return errAPIVersionMustBeSpecified
	}

	return os.Remove(fs.filepathFromResourceName(avk, name))
}

// read tries to read file corresponding to the specified object's name.
// It returns os.ErrNotExist if file does not exist.
func (fs *filesystem) read(path string) (types.Resource[types.APIVersionKind], error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return types.Resource[types.APIVersionKind]{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	rawRes := new(types.Resource[any])
	if err := fs.codec.Unmarshal(b, rawRes); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	out, err := fs.apiServer.Get(types.NewAPIVersionKind(rawRes.APIVersion, rawRes.Kind))
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	// NOTE: we must copy metadata over
	out.Metadata = rawRes.Metadata

	b, err = fs.codec.Marshal(rawRes.Spec)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	if err = fs.codec.Unmarshal(b, &out.Spec); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	return out, nil
}

// joinResourceDirAndBasename joins the resourceDir to the basename
func (fs *filesystem) joinResourceDirAndBasename(basename string) string {
	return filepath.Join(fs.resourceDir, basename)
}

// filepathFromResourceName computes the path to the corresponding resource name.
func (fs *filesystem) filepathFromResourceName(
	avk types.APIVersionKind,
	resourceName string,
) string {
	basename := fs.basename(avk, resourceName)
	return fs.joinResourceDirAndBasename(basename)
}

// filepathFromResourceName computes the resourceDir to the corresponding resource name, based on the naming convention
func (fs *filesystem) tmpFilepathFromResourceName(
	avk types.APIVersionKind,
	resourceName string,
) string {
	basename := fs.basename(avk, resourceName)
	tmpBasename := fmt.Sprintf(".tmp.%s", basename)
	return fs.joinResourceDirAndBasename(tmpBasename)
}

// filepathFromResourceName computes the resource filename, based on the naming convention
func (fs *filesystem) basename(avk types.APIVersionKind, resourceName string) string {
	return strings.ToLower(fmt.Sprintf(
		"%s.%s.%s.%s",
		cleanAPIVersionForFilesystem(avk.APIVersion()),
		avk.Kind(),
		resourceName,
		fs.codec.Encoding(),
	))
}

//----------------------------------------------------------------------------------------------------------------------
// Storage Utils
//----------------------------------------------------------------------------------------------------------------------

// resourceExist checks if a named resource already exist
func resourceExist(storage types.Storage, avk types.APIVersionKind, name string) (bool, error) {
	_, err := storage.Get(avk, name)
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
