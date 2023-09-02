package main

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib"
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
func ParseAPIVersionAndKindFromArgs(cctx *cli.Context) (*vib.APIVersion, *vib.Kind) {
	s := strings.ToLower(cctx.Args().Get(0))
	if s == "" {
		return nil, nil
	}

	apiVersionAndKindRegex := regexp.MustCompile(
		// ${domain_name}/v[0-9]+[a-z0-9]*/${kind}
		vib.RegexAPIVersionAndKind(),
	)

	if apiVersionAndKindRegex.MatchString(s) {
		sl := strings.Split(s, "/")
		apiVersion := vib.APIVersion(strings.Join(sl[0:1], sl[1]))
		kind := vib.Kind(sl[2])
		return &apiVersion, &kind
	}

	if regexp.MustCompile(vib.RegexKindLowered()).MatchString(s) {
		kind := vib.Kind(s)
		return nil, &kind
	}
	return nil, nil
}

func ParseResourceNamesFromArgs(cctx *cli.Context) []string {
	if cctx.Args().Len() < 2 {
		return nil
	}

	regex := regexp.MustCompile(vib.RegexResourceName())

	results := make([]string, 0)
	for _, name := range cctx.Args().Slice()[1:] {
		if !regex.MatchString(name) {
			_ = vib.NewErrAndLog(vib.ErrType, fmt.Sprintf("resource name %q is not supported", name))

			return nil
		}

		results = append(results, name)
	}

	return results
}

func resourceFromFileOrStdin(cctx *cli.Context) (*vib.ResourceDefinition, error) {
	var resource *vib.ResourceDefinition
	var err error

	if file := cctx.String(fileFlagName); file != "" {
		// read file
		resource, err = vib.ReadEncodedFile(file)
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
func unmarshalFromStdin() (*vib.ResourceDefinition, error) {
	logger.Debug("reading resource from stdin")
	// Otherwise read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	var resource *vib.ResourceDefinition
	if err = yaml.Unmarshal(data, resource); err != nil {
		logger.Error(err)

		return nil, err
	}

	return resource, nil
}

func pleaseSpecifyAResourceKind() error {
	return vib.NewErr(vib.ErrArgs, "please specify a resource kind") //nolint:wrapcheck
}

func pleaseSpecifyValidResourceNames() error {
	return vib.NewErr(vib.ErrArgs, "please specify valid resource name(s)") //nolint:wrapcheck
}
