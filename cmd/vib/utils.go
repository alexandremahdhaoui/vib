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
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	debug = "debug"
)

// ParseAPIVersionAndKindFromArgs returns both optionally not nil pointer to vib.APIVersion and vib.Kind
func ParseAPIVersionAndKindFromArgs(cctx *cli.Context) (*api.APIVersion, *api.Kind) {
	s := strings.ToLower(cctx.Args().Get(0))
	if s == "" {
		return nil, nil
	}

	apiVersionAndKindRegex := regexp.MustCompile(
		// ${domain_name}/v[0-9]+[a-z0-9]*/${kind}
		api.RegexAPIVersionAndKind(),
	)

	if apiVersionAndKindRegex.MatchString(s) {
		sl := strings.Split(s, "/")
		apiVersion := api.APIVersion(strings.Join(sl[0:1], sl[1]))
		kind := api.Kind(sl[2])
		return &apiVersion, &kind
	}

	if regexp.MustCompile(api.RegexKindLowered()).MatchString(s) {
		kind := api.Kind(s)
		return nil, &kind
	}
	return nil, nil
}

func ParseResourceNamesFromArgs(cctx *cli.Context) []string {
	if cctx.Args().Len() < 2 {
		return nil
	}

	regex := regexp.MustCompile(api.RegexResourceName())

	results := make([]string, 0)
	for _, name := range cctx.Args().Slice()[1:] {
		if !regex.MatchString(name) {
			_ = logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("resource name %q is not supported", name))

			return nil
		}

		results = append(results, name)
	}

	return results
}

func resourceFromFileOrStdin(cctx *cli.Context) (*api.ResourceDefinition, error) {
	var resource *api.ResourceDefinition
	var err error

	if file := cctx.String(fileFlagName); file != "" {
		// read file
		resource, err = api.ReadEncodedFile(file)
		if err != nil {
			return nil, err
		}
	} else {
		resource, err = unmarshalFromStdin()
		if err != nil {
			return nil, err
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
				logger.Info("switching to debug mode")
				logger.New(true)
			}
			return nil
		},
	}
}

// unmarshalFromStdin only supports vib.YAMLEncoding for now
func unmarshalFromStdin() (*api.ResourceDefinition, error) {
	logger.Debug("reading resource from stdin")
	// Otherwise read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	var resource *api.ResourceDefinition
	if err = yaml.Unmarshal(data, resource); err != nil {
		logger.Error(err)

		return nil, err
	}

	return resource, nil
}

func pleaseSpecifyAResourceKind() error {
	return logger.NewErr(logger.ErrArgs, "please specify a resource kind") //nolint:wrapcheck
}

func pleaseSpecifyValidResourceNames() error {
	return logger.NewErr(logger.ErrArgs, "please specify valid resource name(s)") //nolint:wrapcheck
}
