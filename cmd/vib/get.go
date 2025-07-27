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

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	codecadapter "github.com/alexandremahdhaoui/vib/internal/adapter/codec"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

const getDesc = `
	Description:
		Get resources of kind "KIND". Resources can be optionally filtered by
		name.
	Usage:
		vib get KIND [name0] [name1] [flags]
	Args:
		KIND: the kind of the resource.
		name: name(s) of resources that must be returned. (optional)`

func NewGet(storage types.Storage) Command {
	out := &get{
		apiVersion: "",
		fs:         flag.NewFlagSet("get", flag.ExitOnError),
		outputEnc:  "",
		storage:    storage,
	}

	out.fs.StringVar(
		(*string)(&out.apiVersion),
		"apiVersion",
		"",
		"The APIVersion of the resource to get",
	)

	out.fs.StringVar(
		(*string)(&out.outputEnc),
		"o",
		string(defaultOutputEncoding),
		"The output encoding must be one of [json,yaml]; default is \"yaml\"",
	)

	return out
}

type get struct {
	apiVersion types.APIVersion
	fs         *flag.FlagSet
	outputEnc  string
	storage    types.Storage
}

// FS implements Command.
func (g *get) FS() *flag.FlagSet {
	return g.fs
}

// Description implements Command.
func (g *get) Description() string {
	return getDesc
}

// Run implements Command.
func (g *get) Run() error {
	if g.fs.NArg() < 1 {
		return flaterrors.Join(
			errors.New("[ERROR] \"Get\" expects at least one argument"),
			errors.New(getDesc), //nolint staticcheck
		)
	}

	outputCodec, err := codecadapter.New(types.Encoding(g.outputEnc))
	if err != nil {
		return err
	}

	nameFilter := make(map[string]struct{})
	for i := 1; i < g.fs.NArg(); i++ {
		nameFilter[g.fs.Arg(i)] = struct{}{}
	}

	list, err := List(g.storage, g.apiVersion, types.NewKind(g.fs.Arg(0)), nameFilter)
	if err != nil {
		return err
	}

	// hack to better print a single resource
	var v any = list
	if len(list) == 1 {
		v = list[0]
	}

	b, err := outputCodec.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
