package main

import "github.com/urfave/cli/v2"

const (
	// Flags
	fileFlagName = "file"
)

func fileFlag() *cli.StringFlag {
	return &cli.StringFlag{ //nolint:exhaustivestruct
		Name:    fileFlagName,
		Usage:   "specify path to a file",
		Aliases: []string{"f"},
	}
}
