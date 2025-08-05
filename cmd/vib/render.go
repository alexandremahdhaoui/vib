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
		vib render [flags] KIND NAME
	Args:
		KIND: The kind of the resource to render.
		NAME: The name of the resource to render.`

func NewRender(apiServer types.APIServer, storage types.Storage) Command {
	out := &render{
		apiServer:  apiServer,
		apiVersion: "",
		fs:         flag.NewFlagSet("render", flag.ExitOnError),
		namespace:  "",
		storage:    storage,
	}

	NewAPIVersionFlag(out.fs, &out.apiVersion)
	NewNamespaceFlag(out.fs, &out.namespace)

	return out
}

type render struct {
	apiServer  types.APIServer
	apiVersion types.APIVersion
	fs         *flag.FlagSet
	namespace  string
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
			errors.New("\"RENDER\" requires TWO arguments"),
			errors.New(renderDesc), //nolint staticcheck
		)
	}

	// -- get avk with specific apiVersion
	kind := r.fs.Arg(0)
	name := r.fs.Arg(1)

	// The input apiVersion might be an empty string.
	// This ensure the apiVersion is specified
	avk := types.NewAPIVersionKind(r.apiVersion, kind)
	res, err := r.apiServer.Get(avk)
	if err != nil {
		return err
	}

	specificAvk := types.NewAVKFromResource(res)
	nsName := types.NamespacedName{
		Name:      name,
		Namespace: r.namespace,
	}

	resource, err := r.storage.Get(specificAvk, nsName)
	if err != nil {
		return err
	}

	renderer, ok := any(resource.Spec).(types.Renderer)
	if !ok {
		return fmt.Errorf(
			"cannot render resource: apiVersion=%q,kind=%q,name=%q",
			resource.APIVersion,
			resource.Kind,
			resource.Metadata.Name,
		)
	}

	out, err := renderer.Render(r.storage)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
