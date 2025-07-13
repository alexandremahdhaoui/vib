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
	"log/slog"
	"os"

	"github.com/alexandremahdhaoui/vib/internal/util"
	"github.com/urfave/cli/v2"
	"go.yaml.in/yaml/v3"
)

const (
	cliName = "vib"

	// Commands
	get       = "get"
	create    = "create"
	edit      = "edit"
	deleteCmd = "delete"
	apply     = "apply"
	render    = "render"

	// Flag categories
	basicCategory    = "Basic Commands"
	miscCategory     = "Miscellaneous"
	selectorCategory = "Selectors"
)

var appVersion = "1.0.0" //nolint:gochecknoglobals

func main() {
	// -- new main

	// . Get Storage
	// . Get Resolvers
	// . Create api server

	// -- remove below
	app := &cli.App{ //nolint:exhaustruct,exhaustivestruct
		Name:    cliName,
		Usage:   `vib (pronounced \"vibe\") allows users to intuitively manage their bash environment across all their platforms.`, //nolint:lll
		Version: appVersion,
		Commands: cli.Commands{
			// Basic Commands
			Get(),
			Create(),
			Edit(),
			Delete(),
			Apply(),

			// Render
			Render(),
		},
		Flags: []cli.Flag{debugFlag()},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Get
//----------------------------------------------------------------------------------------------------------------------

func Get() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct,exhaustivestruct
		Name:     get,
		Usage:    "Display one or many resources",
		Category: basicCategory,
		Action: func(cctx *cli.Context) error {
			server, err := fastInit()
			if err != nil {
				return err
			}

			resources, err := GetResources(cctx, server)
			if err != nil {
				return err
			}

			b, err := yaml.Marshal(resources)
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}
}

func GetResources(cctx *cli.Context, server api.APIServer) ([]api.ResourceDefinition, error) {
	apiVersion, kind := ParseAPIVersionAndKindFromArgs(cctx)
	if kind == nil {
		return nil, errPleaseSpecifyAResourceKind()
	}

	names := ParseResourceNamesFromArgs(cctx)
	// Condition were user didn't specify any name
	if len(names) == 0 {
		resources, err := server.Get(apiVersion, *kind, nil)
		if err != nil {
			return nil, err
		}

		return resources, nil
	}

	results := make([]api.ResourceDefinition, 0)
	// Condition were user specified name(s)
	for _, name := range names {
		resources, err := server.Get(apiVersion, *kind, &name)
		if err != nil {
			return nil, err
		}

		results = append(results, resources...)
	}

	return results, nil
}

//----------------------------------------------------------------------------------------------------------------------
// Create
//----------------------------------------------------------------------------------------------------------------------

func Create() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct,exhaustivestruct
		Name:     create,
		Usage:    "Create a resource from a file or from stdin",
		Category: basicCategory,
		Flags: []cli.Flag{
			debugFlag(),
			fileFlag(),
		},
		Action: func(cctx *cli.Context) error {
			apiServer, err := fastInit()
			if err != nil {
				return err
			}

			resource, err := resourceFromFileOrStdin(cctx)
			if err != nil {
				return err
			}

			err = apiServer.Create(resource)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Edit
//----------------------------------------------------------------------------------------------------------------------

func Edit() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct,exhaustivestruct
		Name:     edit,
		Usage:    "Edit a resource",
		Category: basicCategory,
		Action: func(cctx *cli.Context) error {
			_, kind := ParseAPIVersionAndKindFromArgs(cctx)
			if kind == nil {
				return errPleaseSpecifyAResourceKind()
			}

			names := ParseResourceNamesFromArgs(cctx)
			if len(names) == 0 {
				return errPleaseSpecifyValidResourceNames()
			}

			// TODO implement me
			panic("not implemented yet")

			return nil
		},
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Delete
//----------------------------------------------------------------------------------------------------------------------

func Delete() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct,exhaustivestruct
		Name:     deleteCmd,
		Usage:    "Delete a resource by name",
		Category: basicCategory,
		Action: func(cctx *cli.Context) error {
			apiVersion, kind := ParseAPIVersionAndKindFromArgs(cctx)
			if kind == nil {
				return errPleaseSpecifyAResourceKind()
			}

			names := ParseResourceNamesFromArgs(cctx)
			if len(names) == 0 {
				return errPleaseSpecifyValidResourceNames()
			}

			apiServer, err := fastInit()
			if err != nil {
				return err
			}

			for _, name := range names {
				if err = apiServer.Delete(apiVersion, *kind, name); err != nil {
					return err //nolint:wrapcheck
				}
			}

			return nil
		},
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Apply
//----------------------------------------------------------------------------------------------------------------------

func Apply() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct,exhaustivestruct
		Name:  apply,
		Usage: "Apply resource from a file or from stdin",
		Flags: []cli.Flag{
			debugFlag(),
			fileFlag(),
		},
		Action: func(cctx *cli.Context) error {
			apiServer, err := fastInit()
			if err != nil {
				return err
			}

			resource, err := resourceFromFileOrStdin(cctx)
			if err != nil {
				return err
			}

			err = apiServer.Update(
				&resource.APIVersion,
				resource.Kind,
				resource.Metadata.Name,
				resource,
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Render
//----------------------------------------------------------------------------------------------------------------------

func Render() *cli.Command {
	return &cli.Command{ //nolint:exhaustruct,exhaustivestruct
		Name:  render,
		Usage: "Render the designated profile",
		Flags: []cli.Flag{
			debugFlag(),
		},
		Action: func(cctx *cli.Context) error {
			buffer := ""

			server, err := fastInit()
			if err != nil {
				return err
			}

			resources, err := GetResources(cctx, server)
			if err != nil {
				return err
			}

			for _, resource := range resources {
				s, err := vib.Render(&resource, server)
				if err != nil {
					return err
				}

				buffer = util.JoinLine(buffer, s)

			}

			fmt.Println(buffer)

			return nil
		},
	}
}
