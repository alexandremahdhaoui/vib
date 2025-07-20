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
	"fmt"
	"io"
	"os"
)

const (
	cliName = "vib"
)

func main() {
	cmds := map[*flag.FlagSet]func(){
		applyFlagSet(): apply(),
		createFlagSet(),
		delFlagSet(),
		editFlagSet(),
		getFlagSet(),
		renderFlagSet(),
	}

	if len(os.Args) < 2 {
		usage(os.Stderr, cmds)
		os.Exit(1)
	}

	for fs, f := range cmds {
		if os.Args[1] == cmd.Name() {
			f()
		}
	}
}

func usage(w io.Writer, cmds []*flag.FlagSet) {
	_, _ = fmt.Fprintf(w, "USAGE: %s [command]\n", os.Args[0])
	_, _ = fmt.Fprintf(w, "Available Commands:\n")
	for _, cmd := range cmds {
		_, _ = fmt.Fprintf(w, "\t%s\t%s\n\n", cmd.Name(), cmd.Usage)
	}
}

// ---------------------------------------------------------------------
// - APPLY
// ---------------------------------------------------------------------

func applyFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	return fs
}

func apply() error {}

// ---------------------------------------------------------------------
// - CREATE
// ---------------------------------------------------------------------

func createFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	return fs
}

// ---------------------------------------------------------------------
// - DELETE
// ---------------------------------------------------------------------

func delFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	return fs
}

// ---------------------------------------------------------------------
// - EDIT
// ---------------------------------------------------------------------

func editFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	return fs
}

// ---------------------------------------------------------------------
// - GET
// ---------------------------------------------------------------------

func getFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	return fs
}

// ---------------------------------------------------------------------
// - RENDER
// ---------------------------------------------------------------------

func renderFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	return fs
}
