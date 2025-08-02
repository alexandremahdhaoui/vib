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

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	codecadapter "github.com/alexandremahdhaoui/vib/internal/adapter/codec"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

const createDesc = `
	Usage:
		vib create [flags] KIND NAME 
	Description:
		Create a new empty resource with the provided name.
	Args:
		Kind: the kind of the resource.
		NAME: name of the resource to create.`

func NewCreate(
	apiServer types.APIServer,
	storage types.Storage,
) Command {
	out := &create{
		apiServer:  apiServer,
		apiVersion: "",
		fs:         flag.NewFlagSet("create", flag.ExitOnError),
		namespace:  "",
		outputEnc:  "",
		storage:    storage,
	}

	NewAPIVersionFlag(out.fs, &out.apiVersion)
	NewNamespaceFlag(out.fs, &out.namespace)
	NewOutputEncodingFlag(out.fs, &out.outputEnc)

	return out
}

type create struct {
	apiServer  types.APIServer
	apiVersion types.APIVersion
	fs         *flag.FlagSet
	namespace  string
	outputEnc  string
	storage    types.Storage
}

// FS implements Command.
func (g *create) FS() *flag.FlagSet {
	return g.fs
}

// Description implements Command.
func (g *create) Description() string {
	return createDesc
}

// Run implements Command.
func (g *create) Run() error {
	outputCodec, err := codecadapter.New(types.Encoding(g.outputEnc))
	if err != nil {
		return err
	}

	if g.fs.NArg() < 2 {
		return flaterrors.Join(
			errors.New("\"CREATE\" expects TWO argument"),
			errors.New(createDesc), //nolint staticcheck
		)
	}

	kind := g.fs.Arg(0)
	name := g.fs.Arg(1)
	avk := types.NewAPIVersionKind(g.apiVersion, kind)

	res, err := g.apiServer.Get(avk)
	if err != nil {
		return err
	}

	res.Metadata.Name = name
	res.Metadata.Namespace = g.namespace
	if err := g.storage.Create(res); err != nil {
		return err
	}

	slog.Info(
		"Successfully created resource",
		"name", name,
		"apiVersion", res.APIVersion,
		"kind", res.Kind,
	)

	nameFilter := map[string]struct{}{name: {}}
	list, err := List(g.storage, res.APIVersion, types.Kind(g.fs.Arg(0)), nameFilter, g.namespace)
	if err != nil {
		return err
	}

	if len(list) != 1 {
		return errors.New("fetching newly created resource")
	}

	b, err := outputCodec.Marshal(list[0])
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
