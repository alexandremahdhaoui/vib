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
	"regexp"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

type APIServer interface {
	Register(APIVersion, map[Kind]func() any)
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

type (
	APIVersion string
	Kind       string
)

func (v APIVersion) ToLower() APIVersion {
	return APIVersion(strings.ToLower(string(v)))
}

func (v APIVersion) Validate() (APIVersion, error) {
	if apiVersion := v.ToLower(); !APIVersionRegex.
		MatchString(string(apiVersion)) {
		return apiVersion, nil
	}

	return "", flaterrors.Join(
		ErrVal,
		fmt.Errorf("couldn't validate APIVersion %q", v),
	)
}

func (k Kind) ToLower() Kind {
	return Kind(strings.ToLower(string(k)))
}

func (k Kind) Validate() (Kind, error) {
	if kind := k.ToLower(); LoweredKindRegex.MatchString(string(kind)) {
		return kind, nil
	}

	return "", flaterrors.Join(
		ErrVal,
		fmt.Errorf("couldn't validate Kind %q", k),
	)
}

type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"      yaml:"labels,omitempty"`
	Name        string            `json:"name"                  yaml:"name"`
}

func NewMetadata(name string) Metadata {
	return Metadata{Name: name} //nolint:exhaustruct,exhaustivestruct
}

type Resource struct {
	APIVersion APIVersion `json:"apiVersion" yaml:"apiVersion"`
	Kind       Kind       `json:"kind"       yaml:"kind"`
	Metadata   Metadata   `json:"metadata"   yaml:"metadata"`
	Spec       any        `json:"spec"       yaml:"spec"`
}

func NewResource(
	apiVersion APIVersion,
	kind Kind,
	name string,
	spec any,
) *Resource {
	return &Resource{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata:   NewMetadata(name),
		Spec:       spec,
	}
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

func ValidateAPIVersionPtr(ptr *APIVersion) error {
	if ptr == nil {
		return nil
	}

	apiVersion, err := ptr.Validate()
	if err != nil {
		return err
	}

	*ptr = apiVersion

	return nil
}
