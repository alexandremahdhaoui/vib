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
	APIServer interface {
		// Registers a new APIVersion to the APIServer
		Register(avkFactory []AVKFunc)

		// Get will return a zero valued instance of a Resource corresponding
		// to the return AVK
		Get(avk APIVersionKind) (Resource[APIVersionKind], error)
	}

	APIVersionKind interface {
		APIVersion() APIVersion
		Kind() Kind
	}

	Codec interface {
		Marshal(v any) ([]byte, error)
		Unmarshal(b []byte, v any) error
		Encoding() Encoding
	}

	DynamicDecoder[T any] interface {
		Decode(io.Reader) ([]Resource[T], error)
	}

	Reader[T APIVersionKind] interface {
		// List resources.
		Read() ([]Resource[T], error)
	}

	Renderer interface {
		Render(storage Storage) (string, error)
	}

	Storage interface {
		// List resources.
		List(avk APIVersionKind, namespace string) ([]Resource[APIVersionKind], error)

		// Get a resource by name. It returns types.ErrNotFound if the resource
		// cannot be found.
		Get(
			avk APIVersionKind,
			namespacedName NamespacedName,
		) (Resource[APIVersionKind], error)

		// Creates a resource if it does not exist in the store.
		Create(Resource[APIVersionKind]) error

		// Update returns types.ErrNotFound if named resource cannot be found
		Update(v Resource[APIVersionKind]) error

		// Delete a resource in the store. Delete is idempotent.
		Delete(avk APIVersionKind, namespacedName NamespacedName) error
	}

	Validator interface {
		Validate() error
	}
)

const (
	DefaultNamespace   = "default"
	VibSystemNamespace = "vib-system"
)

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

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

type AVKFunc func() APIVersionKind

type Encoding string

const (
	JSONEncoding Encoding = "json"
	YAMLEncoding Encoding = "yaml"
)

type (
	APIVersion = string
	Kind       = string
)

func NewAPIVersion(s string) APIVersion {
	return APIVersion(strings.ToLower(s))
}

type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
}

func NewMetadata(name, namespace string) Metadata {
	return Metadata{
		Name:      name,
		Namespace: name,
	} //nolint:exhaustruct,exhaustivestruct
}

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

// APIVersion implements APIVersionKind.
func (a avk) APIVersion() APIVersion {
	return a.apiVersion
}

// Kind implements APIVersionKind.
func (a avk) Kind() Kind {
	return a.kind
}

func NewAPIVersionKind(apiVersion APIVersion, kind Kind) APIVersionKind {
	return avk{
		apiVersion: apiVersion,
		kind:       kind,
	}
}

func NewAVKFromResource[T any](res Resource[T]) APIVersionKind {
	return avk{
		apiVersion: res.APIVersion,
		kind:       res.Kind,
	}
}

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
