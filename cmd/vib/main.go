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
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	codecadapter "github.com/alexandremahdhaoui/vib/internal/adapter/codec"
	storageadapter "github.com/alexandremahdhaoui/vib/internal/adapter/storage"
	"github.com/alexandremahdhaoui/vib/internal/service"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/pkg/apis/v1alpha1"
)

// // ConfigSpec stores important information to run the vib command line.
// // The config is always stored on disk, thus the Operator for managing Config will always be of type
// // vib.FilesystemOperator.
// type ConfigSpec struct {
// 	// StorageStrategy defines which storage strategy must be used (only filesystem is supported).
// 	StorageStrategy storageadapter.StorageStrategy
// 	// ResourceDir specifies the absolute path to Resource definitions.
// 	// Defaults to CONFIG_DIR/vib/resources
// 	ResourceDir string
// }

const (
	defaultStorageEncoding = types.YAMLEncoding
	defaultOutputEncoding  = types.YAMLEncoding
)

type Command interface {
	Description() string
	FS() *flag.FlagSet
	Run() error
}

func main() {
	// --------------------
	// - INIT
	// --------------------

	// -- apiServer
	apiServer := service.NewAPIServer()
	v1alpha1.RegisterWithManager(apiServer)

	// -- dynamic resource decoder
	drd := codecadapter.NewDynamicResourceDecoder(apiServer)

	// -- storage encoding
	storageCodec, err := codecadapter.New(defaultStorageEncoding)
	if err != nil {
		errAndExit(err)
		return
	}

	// -- vib config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		errAndExit(err)
		return
	}
	vibConfigDir := filepath.Join(userConfigDir, "vib")

	// -- storage
	storage, err := storageadapter.NewFilesystem(
		apiServer,
		storageCodec,
		vibConfigDir,
	)
	if err != nil {
		errAndExit(err)
		return
	}

	// --------------------
	// - DECLARE CMDS
	// --------------------

	cmds := []Command{
		// NewApply(), // Read, UpdateOrCreate
		NewCreate(apiServer, storage),
		NewDelete(apiServer, storage),
		NewEdit(apiServer, drd, storage), // List, EditText, UpdateOrCreate
		NewGet(apiServer, storage),
		// NewGrep(TODO), // List, regexp.Match, Print
		NewRender(apiServer, storage),
	}

	if len(os.Args) < 2 {
		help(os.Stderr, cmds)
		os.Exit(1)
	}

	// --------------------
	// - RUN CMDS
	// --------------------

	found := false
	for _, cmd := range cmds {
		if os.Args[1] != cmd.FS().Name() {
			continue
		}

		found = true
		if len(os.Args) > 1 {
			if err := cmd.FS().Parse(os.Args[2:]); err != nil {
				help(os.Stderr, cmds) // actually not called
				os.Exit(1)            // not called
			}
		}

		if err := cmd.Run(); err != nil {
			errAndExit(err)
			return
		}

		break
	}

	if !found {
		help(os.Stderr, cmds)
		os.Exit(1)
	}
}

func errAndExit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

const usageFmt = `USAGE: %s [command]

Available Commands:

`

func help(w io.Writer, cmds []Command) {
	fmt.Fprintf(w, usageFmt, os.Args[0]) //nolint: errcheck

	for _, cmd := range cmds {
		fmt.Fprintf(w, "%s\n", cmd.FS().Name()) //nolint: errcheck

		// TODO: ensure that description does not exceed (80-indent) charachters per line
		indent := "\t"
		fmt.Fprintf(w, "%s%s\n", indent, cmd.Description()) //nolint: errcheck
		fmt.Fprintf(w, "%sFlags:\n", indent)                //nolint: errcheck

		indent = "\t\t"
		cmd.FS().VisitAll(func(fl *flag.Flag) {
			var longFlag string
			if len(fl.Name) > 1 {
				longFlag = "-"
			}

			fmt.Fprintf(w, "%s-%s%s\t%s\n", indent, longFlag, fl.Name, fl.Usage) //nolint: errcheck
		})

		fmt.Fprintf(w, "\n") //nolint: errcheck
	}

	fmt.Fprintf(w, "\n") //nolint: errcheck
}

// ---------------------------------------------------------------------
// - LIST RESOURCES
// ---------------------------------------------------------------------

type ListArgs struct {
	APIVersion types.APIVersion
	Kind       types.Kind
	NameFilter map[string]struct{}
}

// List must return a list of resources.
// The caller (i.e. "get") can then choose how to format the result.
// List can be used for other calls such as render
func List(
	storage types.Storage,
	apiVersion types.APIVersion,
	kind types.Kind,
	nameFilter map[string]struct{},
) ([]types.Resource[types.APIVersionKind], error) {
	list, err := storage.List(types.NewAPIVersionKind(apiVersion, kind))
	if err != nil {
		return nil, err
	}

	if len(nameFilter) == 0 {
		return list, nil
	}

	out := make([]types.Resource[types.APIVersionKind], 0)
	found := false
	for _, res := range list {
		if _, ok := nameFilter[res.Metadata.Name]; !ok {
			continue
		}
		out = append(out, res)
		found = true
	}

	if !found {
		return nil, errResource(
			"cannot find resource",
			apiVersion,
			kind,
			fmtSet(nameFilter),
		)
	}

	return out, nil
}

// ---------------------------------------------------------------------
// - HELPERS
// ---------------------------------------------------------------------

func errResource(errStr string, apiVersion types.APIVersion, kind types.Kind, name string) error {
	return fmt.Errorf(
		"%s: apiVersion=%q,kind=%q,nameFilters=%q",
		errStr,
		apiVersion,
		kind,
		name,
	)
}

func fmtSet(set map[string]struct{}) string {
	if len(set) == 0 {
		return ""
	}
	var out string
	for s := range set {
		out = fmt.Sprintf("%s, %s", out, s)
	}
	return fmt.Sprintf("[%s]", out[2:])
}
