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

	"github.com/alexandremahdhaoui/vib/internal/types"
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
		fs:         flag.NewFlagSet("edit", flag.ExitOnError),
		outputEnc:  "",
		storage:    storage,
	}

	NewAPIVersionFlag(out.fs, &out.apiVersion)
	NewOutputEncodingFlag(out.fs, &out.outputEnc)

	return out
}

type edit struct {
	apiServer  types.APIServer
	apiVersion types.APIVersion
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
	panic("unimplemented")
}
