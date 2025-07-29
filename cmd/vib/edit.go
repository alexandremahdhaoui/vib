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

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	codecadapter "github.com/alexandremahdhaoui/vib/internal/adapter/codec"
	"github.com/alexandremahdhaoui/vib/internal/types"
	"github.com/alexandremahdhaoui/vib/internal/util"
)

const editDesc = `
	Usage:
		vib edit [flags] KIND NAME [NAME0] [NAME1]
	Description:
		Edit resources interactively.
	Args:
		KIND: the kind of the resource.
		NAME [NAME{X}]: name of resource(s) to edit.`

func NewEdit(
	apiServer types.APIServer,
	drd types.DynamicDecoder[types.APIVersionKind],
	storage types.Storage,
) Command {
	out := &edit{
		apiServer:  apiServer,
		apiVersion: "",
		editor:     "",
		fs:         flag.NewFlagSet("edit", flag.ExitOnError),
		outputEnc:  "",
		storage:    storage,
	}

	NewAPIVersionFlag(out.fs, &out.apiVersion)
	NewOutputEncodingFlag(out.fs, &out.outputEnc)

	out.fs.StringVar(
		&out.editor,
		"editor",
		getDefaultEditor(),
		"Executable to edit resources with. Defaults to $EDITOR, then $VISUAL and finally vi",
	)

	return out
}

type edit struct {
	apiServer  types.APIServer
	apiVersion types.APIVersion
	editor     string
	fs         *flag.FlagSet
	outputEnc  string
	storage    types.Storage
}

// Description implements Command.
func (e *edit) Description() string {
	return editDesc
}

// FS implements Command.
func (e *edit) FS() *flag.FlagSet {
	return e.fs
}

// Run implements Command.
func (e *edit) Run() error {
	// -- 1. List resources
	// -- 2. For each resource: Edit
	// -- 3. For each resource: Update
	// -- 4. List resources
	// -- 5. Print resources

	if e.fs.NArg() < 2 {
		return flaterrors.Join(
			errors.New("[ERROR] \"EDIT\" expects at least TWO argument"),
			errors.New(getDesc), //nolint staticcheck
		)
	}

	outputCodec, err := codecadapter.New(types.Encoding(e.outputEnc))
	if err != nil {
		return err
	}

	// -- get avk with specific apiVersion
	kind := e.fs.Arg(0)
	avk := types.NewAPIVersionKind(e.apiVersion, kind)

	// The input apiVersion might be an empty strine.
	// This ensure the apiVersion is specified
	res, err := e.apiServer.Get(avk)
	if err != nil {
		return err
	}

	nameFilter := make(map[string]struct{})
	for i := 1; i < e.fs.NArg(); i++ {
		nameFilter[e.fs.Arg(i)] = struct{}{}
	}

	// -- 1. List resources
	list, err := List(e.storage, res.APIVersion, kind, nameFilter)
	if err != nil {
		return err
	}

	for _, res := range list {
		// -- 2. For each resource: Edit
		// save resource unmutable properties
		apiVersion := res.APIVersion
		kind := res.Kind
		name := res.Metadata.Name

		// marshal content to edit
		bIn, err := outputCodec.Marshal(res)
		if err != nil {
			return err
		}

		bOut, err := util.EditFile(e.editor, bIn, outputCodec.Encoding())
		if err != nil {
			return err
		}

		// unmarshal edited content
		if err := outputCodec.Unmarshal(bOut, &res); err != nil {
			return err
		}

		if res.APIVersion != apiVersion ||
			res.Kind != kind ||
			res.Metadata.Name != name {
			return errors.New(`[ERROR] "apiVersion", "kind" and "name" are unmutable`)
		}

		// -- 3. For each resource: Update
		if err := e.storage.Update(name, res); err != nil {
			return err
		}

		slog.Info(
			"Successfully updated resource",
			"apiVersion", apiVersion,
			"kind", kind,
			"name", name,
		)
	}

	// -- 4. List resources
	out, err := List(e.storage, res.APIVersion, kind, nameFilter)
	if err != nil {
		return err
	}

	// hack to better print a single resource
	var v any = out
	if len(list) == 1 {
		v = out[0]
	}

	b, err := outputCodec.Marshal(v)
	if err != nil {
		return err
	}

	// -- 5. Print resources
	fmt.Println(string(b))

	return nil
}

func getDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Fallback to vi if no editor environment variable is set
		editor = "vi"
	}
	return editor
}
