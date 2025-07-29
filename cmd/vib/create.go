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
		vib create KIND NAME [flags]
	Description:
		Create a new empty resource with the provided name.
	Args:
		Kind: the kind of the resource.
		resourceName: name(s) of resources that must be returned. (optional)`

func NewCreate(
	apiServer types.APIServer,
	storage types.Storage,
) Command {
	out := &create{
		apiServer:  apiServer,
		apiVersion: "",
		fs:         flag.NewFlagSet("create", flag.ExitOnError),
		outputEnc:  "",
		storage:    storage,
	}

	out.fs.StringVar(
		(*string)(&out.apiVersion),
		"apiVersion",
		"",
		"The APIVersion of the resource to create",
	)

	out.fs.StringVar(
		(*string)(&out.outputEnc),
		"o",
		string(defaultOutputEncoding),
		"The output encoding must be one of [json,yaml]; default is \"yaml\"",
	)

	return out
}

type create struct {
	apiServer  types.APIServer
	apiVersion types.APIVersion
	fs         *flag.FlagSet
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

	if g.fs.NArg() < 1 {
		return flaterrors.Join(
			errors.New("[ERROR] \"Create\" expects at least two argument"),
			errors.New(createDesc), //nolint staticcheck
		)
	}

	apiVersion := types.APIVersion(g.apiVersion)
	kind := types.Kind(g.fs.Arg(0))
	name := g.fs.Arg(1)
	avk := types.NewAPIVersionKind(apiVersion, kind)

	res, err := g.apiServer.Get(avk)
	if err != nil {
		return err
	}

	res.Metadata.Name = name
	if err := g.storage.Create(res); err != nil {
		return err
	}

	slog.Info("successfully created resource", "name", res.Metadata.Name)

	list, err := List(g.storage, res.APIVersion, types.Kind(g.fs.Arg(0)), map[string]struct{}{})
	if err != nil {
		return err
	}

	if len(list) != 1 {
		return errors.New("[ERROR] fetching newly created resource")
	}

	b, err := outputCodec.Marshal(list[0])
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
