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
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

type (
	APIServer interface {
		// Registers a new APIVersion to the APIServer
		Register(APIVersion, map[Kind]func() APIVersionKind)

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
		List() ([]Resource[T], error)

		// Get a resource by name. It returns types.ErrNotFound if the resource
		// cannot be found.
		Get(name string) (Resource[T], error)
	}

	Renderer interface {
		Render(APIServer) (string, error)
	}

	Storage[T APIVersionKind] interface {
		Reader[T]

		// Creates a resource if it does not exist in the store.
		Create(Resource[T]) error

		// Update returns types.ErrNotFound if named resource cannot be found
		Update(oldName string, v Resource[T]) error

		// Delete a resource in the store. Delete is idempotent.
		Delete(name string) error
	}

	StorageServer interface {
		// Returns the Storage corresponding to the input APIVersionKind
		Get(avk APIVersionKind) (Storage[APIVersionKind], error)
	}
)

type AVKFactory func() APIVersionKind

func GetTypedStorage[T APIVersionKind](
	storageServer StorageServer,
) (Storage[T], error) {
	avk := *new(T)

	st, err := storageServer.Get(avk)
	if err != nil {
		return nil, err
	}

	typed, ok := st.(Storage[T])
	if !ok {
		return nil, flaterrors.Join(
			ErrType,
			ErrNotFound,
			fmt.Errorf(
				"storage cannot be found for: apiVersion=%q, kind=%q",
				avk.APIVersion(),
				avk.Kind(),
			),
		)
	}

	return typed, nil
}

var (
	KindRegex         = regexp.MustCompile(`([A-Z][a-z]+)+`)
	LoweredKindRegex  = regexp.MustCompile(`[a-z]+`)
	ResourceNameRegex = regexp.MustCompile(`[a-z][a-z0-9]+(\-[a-z0-9]+)*`)
	APIVersionRegex   = regexp.MustCompile(
		`([a-z0-9]+([a-z0-9]+)*)+(\.([a-z0-9]+([a-z0-9]+)*)+)*(\.[a-z]+)+/v[0-9]+[a-z0-9]*`,
	)
	APIVersionAndKindRegex = regexp.MustCompile(
		fmt.Sprintf("%s/%s", APIVersionRegex.String(), LoweredKindRegex.String()),
	)
)

type Encoding string

const (
	JSONEncoding Encoding = "json"
	YAMLEncoding Encoding = "yaml"
)

type (
	APIVersion string
	Kind       string
)

func NewAPIVersion(s string) APIVersion {
	return APIVersion(strings.ToLower(s))
}

func ValidateAPIVersion(apiVersion APIVersion) error {
	if !APIVersionRegex.MatchString(string(apiVersion)) {
		return flaterrors.Join(
			ErrVal,
			fmt.Errorf("invalid APIVersion %q", apiVersion),
		)
	}
	return nil
}

func NewKind(s string) Kind {
	return Kind(strings.ToLower(s))
}

func ValidateKind(kind Kind) error {
	if !LoweredKindRegex.MatchString(string(kind)) {
		return flaterrors.Join(
			ErrVal,
			fmt.Errorf("couldn't validate Kind %q", kind),
		)
	}
	return nil
}

type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"      yaml:"labels,omitempty"`
	Name        string            `json:"name"                  yaml:"name"`
}

func NewMetadata(name string) Metadata {
	return Metadata{Name: name} //nolint:exhaustruct,exhaustivestruct
}

type Resource[T any] struct {
	APIVersion APIVersion `json:"apiVersion" yaml:"apiVersion"`
	Kind       Kind       `json:"kind"       yaml:"kind"`
	Metadata   Metadata   `json:"metadata"   yaml:"metadata"`
	Spec       T          `json:"spec"       yaml:"spec"`
}

func NewResource[T any](
	apiVersion APIVersion,
	kind Kind,
	name string,
	spec T,
) (Resource[T], error) {
	if err := ValidateAPIVersion(apiVersion); err != nil {
		return Resource[T]{}, err
	}
	if err := ValidateKind(kind); err != nil {
		return Resource[T]{}, err
	}
	if err := ValidateResourceName(name); err != nil {
		return Resource[T]{}, err
	}

	return Resource[T]{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata:   NewMetadata(name),
		Spec:       spec,
	}, nil
}

func ValidateResourceName(s string) error {
	if ResourceNameRegex.MatchString(s) {
		return nil
	}

	return flaterrors.Join(
		ErrVal,
		fmt.Errorf("couldn't validate resource name %q", s),
	)
}

func ValidateResourceNamePtr(ptr *string) error {
	if ptr == nil {
		return nil
	}

	if ResourceNameRegex.MatchString(*ptr) {
		return nil
	}

	return flaterrors.Join(
		ErrVal,
		fmt.Errorf("couldn't validate resource name %q", *ptr),
	)
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

func NewAVKFromResource[T any](res Resource[T]) APIVersionKind {
	return avk{
		apiVersion: res.APIVersion,
		kind:       res.Kind,
	}
}
