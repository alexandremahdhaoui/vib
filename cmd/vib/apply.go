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
	"log/slog"
	"os"

	"github.com/alexandremahdhaoui/vib/internal/types"
)

// NewApply creates a new "apply" command.
func NewApply(
	decoder types.DynamicDecoder[types.APIVersionKind],
	storage types.Storage,
) Command {
	out := &apply{
		decoder:   decoder,
		filePath:  "",
		fs:        flag.NewFlagSet("apply", flag.ExitOnError),
		namespace: "",
		storage:   storage,
	}

	out.fs.StringVar(
		&out.filePath,
		"f",
		"",
		`The name of the file to apply. Users may use "-" to read from Stdin`,
	)

	NewNamespaceFlag(out.fs, &out.namespace)

	return out
}

const applyDesc = `
	Usage:
		vib apply [flags]
	Description:
		Create or edit the the provided resources.`

// apply holds the dependencies and flags for the "apply" command.
type apply struct {
	decoder   types.DynamicDecoder[types.APIVersionKind]
	filePath  string
	fs        *flag.FlagSet
	namespace string
	storage   types.Storage
}

// Description implements the Command interface.
func (a *apply) Description() string {
	return applyDesc
}

// FS implements the Command interface.
func (a *apply) FS() *flag.FlagSet {
	return a.fs
}

// Run implements the Command interface.
func (a *apply) Run() error {
	filePath := a.filePath
	if filePath == "" {
		return errors.New(`a valid file must be provided using the "-f" flag`)
	}

	if filePath == "-" {
		filePath = os.Stdin.Name()
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	list, err := a.decoder.Decode(f)
	if err != nil {
		return err
	}

	for _, res := range list {
		// -- set namespace if namespace is not specified or flag is set.
		if res.Metadata.Namespace == "" || a.namespace != "default" {
			res.Metadata.Namespace = a.namespace
		}

		if err := types.ValidateResource(res); err != nil {
			return err
		}

		verb := "created"
		err := a.storage.Create(res)
		if errors.Is(err, types.ErrExists) {
			err = a.storage.Update(res)
			verb = "updated"
		}
		if err != nil {
			return err
		}

		slog.Info(
			fmt.Sprintf("Successfully %s resource", verb),
			"name", res.Metadata.Name,
			"apiVersion", res.APIVersion,
			"kind", res.Kind,
			"namespace", res.Metadata.Namespace,
		)
	}

	return nil
}
