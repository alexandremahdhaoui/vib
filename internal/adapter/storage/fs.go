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
	defaultNSDir := filepath.Join(resourceDir, types.DefaultNamespace)
	if err := os.MkdirAll(defaultNSDir, 0777); err != nil {
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
	namespace string,
) ([]types.Resource[types.APIVersionKind], error) {
	if err := types.ValidateAPIVersion(avk.APIVersion()); err != nil {
		return nil, flaterrors.Join(err, errAPIVersionMustBeSpecified)
	}

	if err := types.ValidateNamespace(namespace); err != nil {
		return nil, err
	}

	out := make([]types.Resource[types.APIVersionKind], 0)

	namespaceAbsPath := fs.computeNamespaceAbsPath(namespace)

	// Get all instance of T.
	dentries, err := os.ReadDir(namespaceAbsPath)
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
		filename := dentry.Name()
		if dentry.IsDir() || !regex.MatchString(filename) {
			continue
		}

		resourceAbsPath := filepath.Join(namespaceAbsPath, filename)
		v, err := fs.read(resourceAbsPath)
		if err != nil {
			return nil, err
		}

		out = append(out, v)
	}

	return out, nil
}

func (fs *filesystem) Get(
	avk types.APIVersionKind,
	nsName types.NamespacedName,
) (types.Resource[types.APIVersionKind], error) {
	if err := types.ValidateAPIVersion(avk.APIVersion()); err != nil {
		return types.Resource[types.APIVersionKind]{}, flaterrors.Join(
			err,
			errAPIVersionMustBeSpecified,
		)
	}

	if err := types.ValidateNamespacedName(nsName); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	v, err := fs.read(fs.computeResourceAbsPath(avk, nsName, false))
	if os.IsNotExist(err) {
		return types.Resource[types.APIVersionKind]{}, flaterrors.Join(err, types.ErrNotFound)
	} else if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	return v, nil
}

// Create should create only if file does not already exist.
func (fs *filesystem) Create(res types.Resource[types.APIVersionKind]) error {
	if err := types.ValidateResource(res); err != nil {
		return err
	}

	nsName := types.NamespacedName{
		Name:      res.Metadata.Name,
		Namespace: res.Metadata.Namespace,
	}

	exist, err := resourceExist(fs, res.Spec, nsName)
	if err != nil {
		return err
	}

	if exist {
		return flaterrors.Join(
			types.ErrExists,
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

func (fs *filesystem) Update(v types.Resource[types.APIVersionKind]) error {
	if err := types.ValidateResource(v); err != nil {
		return err
	}

	if err := fs.writeAtomic(v); err != nil {
		return err
	}

	return nil
}

func (fs *filesystem) Delete(avk types.APIVersionKind, nsName types.NamespacedName) error {
	if err := types.ValidateAPIVersion(avk.APIVersion()); err != nil {
		return flaterrors.Join(err, errAPIVersionMustBeSpecified)
	}

	if err := types.ValidateNamespacedName(nsName); err != nil {
		return err
	}

	return os.Remove(fs.computeResourceAbsPath(avk, nsName, false))
}

func (fs *filesystem) writeAtomic(v types.Resource[types.APIVersionKind]) error {
	nsName := types.NamespacedName{
		Name:      v.Metadata.Name,
		Namespace: v.Metadata.Namespace,
	}

	tmpDest := fs.computeResourceAbsPath(v.Spec, nsName, true)
	dest := fs.computeResourceAbsPath(v.Spec, nsName, false)

	b, err := fs.codec.Marshal(v)
	if err != nil {
		return err
	}

	// ensure namespace directory exists
	destDir := fs.computeNamespaceAbsPath(nsName.Namespace)
	if err = os.MkdirAll(destDir, 0777); err != nil {
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

	raw := new(types.Resource[any])
	if err := fs.codec.Unmarshal(b, raw); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	out, err := fs.apiServer.Get(types.NewAPIVersionKind(raw.APIVersion, raw.Kind))
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	out.Metadata = raw.Metadata // NOTE: metadata must be copied

	b, err = fs.codec.Marshal(raw.Spec)
	if err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	if err = fs.codec.Unmarshal(b, &out.Spec); err != nil {
		return types.Resource[types.APIVersionKind]{}, err
	}

	return out, nil
}

// filepathByNamespacedName computes the resource filename, based on the naming convention
func (fs *filesystem) basename(avk types.APIVersionKind, resourceName string) string {
	return strings.ToLower(fmt.Sprintf(
		"%s.%s.%s.%s",
		cleanAPIVersionForFilesystem(avk.APIVersion()),
		avk.Kind(),
		resourceName,
		fs.codec.Encoding(),
	))
}

// computeNamespaceAbsPath joins the resourceDir to the basename
func (fs *filesystem) computeNamespaceAbsPath(
	namespace string,
) string {
	if namespace == "" {
		namespace = types.DefaultNamespace
	}
	return filepath.Join(fs.resourceDir, namespace)
}

// computeResourceAbsPath joins the resourceDir to the basename
func (fs *filesystem) computeResourceAbsPath(
	avk types.APIVersionKind,
	nsName types.NamespacedName,
	isTmp bool,
) string {
	basename := fs.basename(avk, nsName.Name)
	if isTmp {
		basename = fmt.Sprintf(".tmp.%s", basename)
	}
	return filepath.Join(
		fs.computeNamespaceAbsPath(nsName.Namespace),
		basename,
	)
}

//----------------------------------------------------------------------------------------------------------------------
// Storage Utils
//----------------------------------------------------------------------------------------------------------------------

// resourceExist checks if a named resource already exist
func resourceExist(
	storage types.Storage,
	avk types.APIVersionKind,
	nsName types.NamespacedName,
) (bool, error) {
	_, err := storage.Get(avk, nsName)
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
