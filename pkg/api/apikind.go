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
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"regexp"
	"strings"
)

type (
	APIVersion string
	Kind       string
	APIKind    interface {
		APIVersion() APIVersion
		Kind() Kind
		Operator() Operator
	}
)

func NewAPIKind(apiVersion APIVersion, kind Kind, operator Operator) APIKind {
	return &concreteAPIKind{
		apiVersion: apiVersion,
		kind:       kind,
		operator:   operator,
	}
}

type concreteAPIKind struct {
	apiVersion APIVersion
	kind       Kind
	operator   Operator
}

func (c *concreteAPIKind) APIVersion() APIVersion {
	return c.apiVersion
}

func (c *concreteAPIKind) Kind() Kind {
	return c.kind
}

func (c *concreteAPIKind) Operator() Operator {
	return c.operator
}

func RegexAPIVersion() string {
	return "([a-z0-9]+([a-z0-9]+)*)+(\\.([a-z0-9]+([a-z0-9]+)*)+)*(\\.[a-z]+)+/v[0-9]+[a-z0-9]*"
}

func RegexKind() string {
	return "([A-Z][a-z]+)+"
}

func RegexKindLowered() string {
	return "[a-z]+"
}

func RegexAPIVersionAndKind() string {
	return fmt.Sprintf("%s/%s", RegexAPIVersion(), RegexKindLowered())
}

func RegexResourceName() string {
	return "[a-z][a-z0-9]+(\\-[a-z0-9]+)*"
}

func (v APIVersion) ToLower() APIVersion {
	return APIVersion(strings.ToLower(string(v)))
}

func (v APIVersion) Validate() (APIVersion, error) {
	if apiVersion := v.ToLower(); regexp.MustCompile(RegexAPIVersion()).MatchString(string(apiVersion)) {
		return apiVersion, nil
	}

	return "", logger.NewErrAndLog(logger.ErrValidation, fmt.Sprintf("couldn't validate APIVersion %q", v))
}

func (k Kind) ToLower() Kind {
	return Kind(strings.ToLower(string(k)))
}

func (k Kind) Validate() (Kind, error) {
	if kind := k.ToLower(); regexp.MustCompile(RegexKindLowered()).MatchString(string(kind)) {
		return kind, nil
	}

	return "", logger.NewErrAndLog(logger.ErrValidation, fmt.Sprintf("couldn't validate Kind %q", k))
}

func ValidateResourceName(s string) error {
	if regexp.MustCompile(RegexResourceName()).MatchString(s) {
		return nil
	}

	return logger.NewErrAndLog(logger.ErrValidation, fmt.Sprintf("couldn't validate resource name %q", s))
}

func ValidateResourceNamePtr(ptr *string) error {
	if ptr == nil {
		return nil
	}

	if regexp.MustCompile(RegexResourceName()).MatchString(*ptr) {
		return nil
	}

	return logger.NewErrAndLog(logger.ErrValidation, fmt.Sprintf("couldn't validate resource name %q", *ptr))
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
