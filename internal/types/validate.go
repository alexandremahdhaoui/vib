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
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
)

var (
	KindRegex         = regexp.MustCompile(`^([A-Z][a-z]+)+$`)
	LoweredKindRegex  = regexp.MustCompile(`^[a-z]+$`)
	ResourceNameRegex = regexp.MustCompile(`^[a-z][a-z0-9]+(\-[a-z0-9]+)*$`)
	APIVersionRegex   = regexp.MustCompile(
		`^([a-z0-9]+([a-z0-9]+)*)+(\.([a-z0-9]+([a-z0-9]+)*)+)*(\.[a-z]+)+/v[0-9]+[a-z0-9]*$`,
	)
	APIVersionAndKindRegex = regexp.MustCompile(
		`^([a-z0-9]+([a-z0-9]+)*)+(\.([a-z0-9]+([a-z0-9]+)*)+)*(\.[a-z]+)+/v[0-9]+[a-z0-9]*/([A-Z][a-z]+)+$`,
	)
)

func ValidateResource[T any](resource Resource[T]) error {
	var errs error

	for _, err := range []error{
		ValidateAPIVersion(resource.APIVersion),
		ValidateKind(resource.Kind),
		ValidateMetadata(resource.Metadata),
		validateSpecIfApplicable(resource.Spec),
	} {
		errs = flaterrors.Join(errs, err)
	}

	if errs != nil {
		return flaterrors.Join(
			errs,
			errors.New("invalid resource"),
			fmt.Errorf("resource %q in namespace %q",
				resource.Metadata.Name,
				resource.Metadata.Namespace,
			),
		)
	}

	return nil
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

func ValidateKind(kind Kind) error {
	loweredKind := strings.ToLower(string(kind))
	if !LoweredKindRegex.MatchString(loweredKind) {
		return flaterrors.Join(
			ErrVal,
			fmt.Errorf("invalid Kind %q", kind),
		)
	}
	return nil
}

func ValidateMetadata(md Metadata) error {
	// TODO: validate annotations
	// TODO: validate labels
	if err := ValidateName(md.Name); err != nil {
		return err
	}
	if err := ValidateNamespace(md.Name); err != nil {
		return err
	}
	return nil
}

func ValidateNamespacedName(nsName NamespacedName) error {
	if err := ValidateName(nsName.Name); err != nil {
		return err
	}
	if err := ValidateNamespace(nsName.Namespace); err != nil {
		return err
	}
	return nil
}

func ValidateName(s string) error {
	if !ResourceNameRegex.MatchString(s) {
		return flaterrors.Join(
			ErrVal,
			fmt.Errorf("invalid name %q", s),
		)
	}
	return nil
}

func ValidateNamespace(s string) error {
	if !ResourceNameRegex.MatchString(s) {
		return flaterrors.Join(
			ErrVal,
			fmt.Errorf("invalid namespace %q", s),
		)
	}
	return nil
}

func validateSpecIfApplicable(v any) error {
	valider, ok := v.(Validator)
	if !ok {
		return nil
	}
	if err := valider.Validate(); err != nil {
		return err
	}
	return nil
}
