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
	"github.com/alexandremahdhaoui/vib/internal/types"
)

const renderDesc = `
	Description:
		Render a resource of kind "KIND".
	Usage:
		vib render KIND NAME {flags}
	Args:
		KIND: The kind of the resource to render.
		NAME: The name of the resource to render.`

func NewRender(storage types.Storage) Command {
	out := &render{
		apiVersion: "",
		fs:         flag.NewFlagSet("render", flag.ExitOnError),
		storage:    storage,
	}

	out.fs.StringVar(
		(*string)(&out.apiVersion),
		"apiVersion",
		"",
		"The APIVersion of the resource to render",
	)

	return out
}

type render struct {
	apiVersion types.APIVersion
	fs         *flag.FlagSet
	storage    types.Storage
}

// FS implements Command.
func (r *render) FS() *flag.FlagSet {
	return r.fs
}

// Description implements Command.
func (r *render) Description() string {
	return renderDesc
}

// Run implements Command.
func (r *render) Run() error {
	if r.fs.NArg() < 2 {
		return flaterrors.Join(
			errors.New("[ERROR] \"\" requires at least two arguments"),
			errors.New(renderDesc), //nolint staticcheck
		)
	}

	kind := types.NewKind(r.fs.Arg(0))
	nameFilter := map[string]struct{}{
		r.fs.Arg(1): {},
	}

	list, err := List(r.storage, r.apiVersion, kind, nameFilter)
	if err != nil {
		return err
	}

	res := list[0]
	renderer, ok := any(res.Spec).(types.Renderer)
	if !ok {
		return fmt.Errorf(
			"cannot render resource: apiVersion=%q,kind=%q,name=%q",
			res.APIVersion,
			res.Kind,
			res.Metadata.Name,
		)
	}

	out, err := renderer.Render(r.storage)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
