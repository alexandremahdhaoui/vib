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
	"log/slog"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

const deleteDesc = `
	Usage:
		vib delete [flags] KIND NAME [NAME0] [NAME1]
	Description:
		Delete a new empty resource with the provided name.
	Args:
		KIND: the kind of the resource to delete.
		NAME [NAME{X}]: name(s) of resources to delete.`

func NewDelete(
	apiServer types.APIServer,
	storage types.Storage,
) Command {
	out := &del{
		apiServer:  apiServer,
		apiVersion: "",
		fs:         flag.NewFlagSet("delete", flag.ExitOnError),
		storage:    storage,
	}

	NewAPIVersionFlag(out.fs, &out.apiVersion)

	return out
}

type del struct {
	apiServer  types.APIServer
	apiVersion types.APIVersion
	fs         *flag.FlagSet
	storage    types.Storage
}

// Description implements Command.
func (d *del) Description() string {
	return deleteDesc
}

// FS implements Command.
func (d *del) FS() *flag.FlagSet {
	return d.fs
}

// Run implements Command.
func (d *del) Run() error {
	if d.fs.NArg() < 2 {
		return flaterrors.Join(
			errors.New("[ERROR] \"Delete\" expects at least two argument"),
			errors.New(createDesc), //nolint staticcheck
		)
	}

	names := make([]string, d.fs.NArg()-1)
	for i := 1; i < d.fs.NArg(); i++ {
		names[i-1] = d.fs.Arg(i)
	}

	kind := d.fs.Arg(0)
	avk := types.NewAPIVersionKind(d.apiVersion, kind)
	// The input apiVersion might be an empty string.
	// This ensure the apiVersion is specified
	res, err := d.apiServer.Get(avk)
	if err != nil {
		return err
	}

	specificAvk := types.NewAVKFromResource(res)
	for _, name := range names {
		if err := d.storage.Delete(specificAvk, name); err != nil {
			return err
		}

		slog.Info(
			"Successfully deleted resource",
			"name", name,
			"apiVersion", specificAvk.APIVersion(),
			"kind", specificAvk.Kind(),
		)
	}

	return nil
}
