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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/alexandremahdhaoui/vib/internal/types"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	debug = "debug"
)

// ParseAPIVersionAndKindFromArgs returns both optionally not nil pointer to vib.APIVersion and vib.Kind
func ParseAPIVersionAndKindFromArgs(cctx *cli.Context) (*types.APIVersion, *types.Kind) {
	s := strings.ToLower(cctx.Args().Get(0))
	if s == "" {
		return nil, nil
	}

	// ${domain_name}/v[0-9]+[a-z0-9]*/${kind}
	if types.APIVersionAndKindRegex.MatchString(s) {
		sl := strings.Split(s, "/")
		apiVersion := types.APIVersion(strings.Join(sl[0:1], sl[1]))
		kind := types.Kind(sl[2])
		return &apiVersion, &kind
	}

	if types.LoweredKindRegex.MatchString(s) {
		kind := types.Kind(s)
		return nil, &kind
	}

	// should this return an err
	return nil, nil
}

func ParseResourceNamesFromArgs(cctx *cli.Context) ([]string, error) {
	if cctx.Args().Len() < 2 {
		return nil, errors.New("expected 2 arguments")
	}

	results := make([]string, 0)
	for _, name := range cctx.Args().Slice()[1:] {
		if !types.ResourceNameRegex.MatchString(name) {
			return nil, flaterrors.Join(
				types.ErrType,
				fmt.Errorf("resource with name %q is invalid", name),
			)
		}

		results = append(results, name)
	}

	return results, nil
}

func resourceFromFileOrStdin(cctx *cli.Context) (types.Resource[json.RawMessage], error) {
	var resource types.Resource[json.RawMessage]
	var err error

	if file := cctx.String(fileFlagName); file != "" {
		// read file
		resource, err = codecadapter.ReadEncodedFile(file)
		if err != nil {
			return resource, err
		}
	} else {
		resource, err = unmarshalFromStdin()
		if err != nil {
			return resource, err
		}
	}
	return resource, nil
}

func debugFlag() *cli.BoolFlag {
	return &cli.BoolFlag{ //nolint:exhaustruct,exhaustivestruct
		Name:     debug,
		Category: miscCategory,
		Action: func(_ *cli.Context, b bool) error {
			if b {
				// -- switching to debug mode
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
			return nil
		},
	}
}

// unmarshalFromStdin only supports vib.YAMLEncoding for now
func unmarshalFromStdin() (types.Resource[json.RawMessage], error) {
	slog.Debug("reading resource from stdin")
	// Otherwise read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return types.Resource[json.RawMessage]{}, err
	}

	var resource types.Resource[json.RawMessage]
	if err = yaml.Unmarshal(data, resource); err != nil {
		return types.Resource[json.RawMessage]{}, err
	}

	return resource, nil
}

func errPleaseSpecifyAResourceKind() error {
	return flaterrors.Join(
		types.ErrArgs,
		errors.New("please specify a resource kind"),
	) //nolint:wrapcheck
}

func errPleaseSpecifyValidResourceNames() error {
	return flaterrors.Join(
		types.ErrArgs,
		errors.New("please specify valid resource name(s)"),
	) //nolint:wrapcheck
}
