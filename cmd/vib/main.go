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
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
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

// Command is the interface that all commands must implement.
type Command interface {
	// Description returns a short description of the command.
	Description() string
	// FS returns the command's FlagSet.
	FS() *flag.FlagSet
	// Run executes the command.
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
	storageCodec, err := NewCodec(defaultStorageEncoding)
	if err != nil {
		logErrAndExit(err)
		return
	}

	// -- vib config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		logErrAndExit(err)
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
		logErrAndExit(err)
		return
	}

	// --------------------
	// - INIT VIB SYSTEM
	// --------------------
	if err := initVibSystemNamespace(storage); err != nil {
		logErrAndExit(err)
		return
	}

	// --------------------
	// - DECLARE CMDS
	// --------------------

	cmds := []Command{
		NewApply(drd, storage), // Read, UpdateOrCreate
		NewCreate(apiServer, storage),
		NewDelete(apiServer, storage),
		NewEdit(apiServer, storage), // List, EditText, UpdateOrCreate
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
			logErrAndExit(err)
			return
		}

		break
	}

	if !found {
		help(os.Stderr, cmds)
		os.Exit(1)
	}
}

func logErrAndExit(err error) {
	slog.Error(err.Error())
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

// ListArgs holds the arguments for the List function.
type ListArgs struct {
	// APIVersion is the API version of the resources to list.
	APIVersion types.APIVersion
	// Kind is the kind of the resources to list.
	Kind types.Kind
	// NameFilter is a map of names to filter the resources by.
	NameFilter map[string]struct{}
}

// List returns a list of resources.
// The caller (e.g., the "get" command) can then choose how to format the result.
// List can be used for other calls such as "render".
func List(
	storage types.Storage,
	apiVersion types.APIVersion,
	kind types.Kind,
	nameFilter map[string]struct{},
	namespace string,
) ([]types.Resource[types.APIVersionKind], error) {
	avk := types.NewAPIVersionKind(apiVersion, kind)
	list, err := storage.List(avk, namespace)
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

// initVibSystemNamespace initialize the vib system namespace.
func initVibSystemNamespace(storage types.Storage) error {
	for _, resolver := range v1alpha1.DefaultAVKResolver() {
		fixedResolver, ok := any(resolver).(types.Resource[types.APIVersionKind])
		if !ok {
			panic("Please fix your commit before submitting a PR")
		}

		fixedResolver.Metadata.Namespace = types.VibSystemNamespace
		if err := storage.Create(fixedResolver); err != nil && !errors.Is(err, types.ErrExists) {
			return err
		}
	}

	return nil
}
