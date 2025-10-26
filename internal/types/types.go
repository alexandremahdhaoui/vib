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

package types

import (
	"io"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

type (
	// APIServer is the interface that defines the methods for an API server.
	APIServer interface {
		// Register registers a new APIVersion to the APIServer.
		Register(avkFactory []AVKFunc)

		// Get will return a zero valued instance of a Resource corresponding
		// to the return AVK.
		Get(avk APIVersionKind) (Resource[APIVersionKind], error)
	}

	// APIVersionKind is the interface that defines the methods for an API version and kind.
	APIVersionKind interface {
		APIVersion() APIVersion
		Kind() Kind
	}

	// Codec is the interface that defines the methods for a codec.
	Codec interface {
		Marshal(v any) ([]byte, error)
		Unmarshal(b []byte, v any) error
		Encoding() Encoding
	}

	// DynamicDecoder is the interface that defines the methods for a dynamic decoder.
	DynamicDecoder[T any] interface {
		Decode(io.Reader) ([]Resource[T], error)
	}

	// Reader is the interface that defines the methods for a reader.
	Reader[T APIVersionKind] interface {
		// Read reads and returns a list of resources.
		Read() ([]Resource[T], error)
	}

	// Renderer is the interface that defines the methods for a renderer.
	Renderer interface {
		Render(storage Storage) (string, error)
	}

	// Storage is the interface that defines the methods for a storage.
	Storage interface {
		// List lists resources.
		List(avk APIVersionKind, namespace string) ([]Resource[APIVersionKind], error)

		// Get gets a resource by name. It returns types.ErrNotFound if the resource
		// cannot be found.
		Get(
			avk APIVersionKind,
			namespacedName NamespacedName,
		) (Resource[APIVersionKind], error)

		// Create creates a resource if it does not exist in the store.
		Create(Resource[APIVersionKind]) error

		// Update updates a resource. It returns types.ErrNotFound if named resource cannot be found.
		Update(v Resource[APIVersionKind]) error

		// Delete deletes a resource in the store. Delete is idempotent.
		Delete(avk APIVersionKind, namespacedName NamespacedName) error
	}

	// Validator is the interface that defines the methods for a validator.
	Validator interface {
		Validate() error
	}
)

const (
	// DefaultNamespace is the default namespace.
	DefaultNamespace = "default"
	// VibSystemNamespace is the namespace for vib system resources.
	VibSystemNamespace = "vib-system"
)

// NamespacedName is a namespaced name.
type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// NewNamespacedNameFromMetadata returns a new NamespacedName from the given metadata.
func NewNamespacedNameFromMetadata(metadata Metadata) NamespacedName {
	namespace := metadata.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}

	return NamespacedName{
		Name:      metadata.Name,
		Namespace: namespace,
	}
}

// AVKFunc is a function that returns an APIVersionKind.
type AVKFunc func() APIVersionKind

// Encoding is the encoding of a resource.
type Encoding string

const (
	// JSONEncoding is the JSON encoding.
	JSONEncoding Encoding = "json"
	// YAMLEncoding is the YAML encoding.
	YAMLEncoding Encoding = "yaml"
)

type (
	// APIVersion is the API version of a resource.
	APIVersion = string
	// Kind is the kind of a resource.
	Kind = string
)

// NewAPIVersion returns a new APIVersion.
func NewAPIVersion(s string) APIVersion {
	return APIVersion(strings.ToLower(s))
}

// Metadata is the metadata of a resource.
type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
}

// NewMetadata returns a new Metadata.
func NewMetadata(name, namespace string) Metadata {
	return Metadata{
		Name:      name,
		Namespace: name,
	} //nolint:exhaustruct,exhaustivestruct
}

// Resource is a generic resource.
type Resource[T any] struct {
	APIVersion APIVersion `json:"apiVersion"`
	Kind       Kind       `json:"kind"`
	Metadata   Metadata   `json:"metadata"`
	Spec       T          `json:"spec"`
}

type avk struct {
	apiVersion APIVersion
	kind       Kind
}

// APIVersion implements the APIVersionKind interface.
func (a avk) APIVersion() APIVersion {
	return a.apiVersion
}

// Kind implements the APIVersionKind interface.
func (a avk) Kind() Kind {
	return a.kind
}

// NewAPIVersionKind returns a new APIVersionKind.
func NewAPIVersionKind(apiVersion APIVersion, kind Kind) APIVersionKind {
	return avk{
		apiVersion: apiVersion,
		kind:       kind,
	}
}

// NewAVKFromResource returns a new APIVersionKind from the given resource.
func NewAVKFromResource[T any](res Resource[T]) APIVersionKind {
	return avk{
		apiVersion: res.APIVersion,
		kind:       res.Kind,
	}
}

// GetTypedResourceFromStorage returns a typed resource from the storage.
func GetTypedResourceFromStorage[T APIVersionKind](
	storage Storage,
	namespacedName NamespacedName,
	v T,
) (Resource[T], error) {
	res, err := storage.Get(v, namespacedName)
	if err != nil {
		return Resource[T]{}, err
	}

	out := Resource[T]{
		APIVersion: res.APIVersion,
		Kind:       res.Kind,
		Metadata:   res.Metadata,
	}

	var ok bool
	out.Spec, ok = res.Spec.(T)
	if !ok {
		return Resource[T]{}, ErrType
	}

	return out, nil
}

// ListTypedResourceFromStorage returns a list of typed resources from the storage.
func ListTypedResourceFromStorage[T APIVersionKind](
	storage Storage,
	namespace string,
	v T,
) ([]Resource[T], error) {
	list, err := storage.List(v, namespace)
	if err != nil {
		return nil, err
	}

	out := make([]Resource[T], 0, len(list))
	for i, res := range list {
		typed := Resource[T]{
			APIVersion: res.APIVersion,
			Kind:       res.Kind,
			Metadata:   res.Metadata,
		}

		var ok bool
		typed.Spec, ok = res.Spec.(T)
		if !ok {
			return nil, flaterrors.Join(ErrType, ErrAtIndex(i))
		}

		out = append(out, typed)
	}

	return out, nil
}
