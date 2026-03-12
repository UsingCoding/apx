package main

import (
	"github.com/urfave/cli/v3"
)

func reg() *cli.Command {
	return &cli.Command{
		Name:  "reg",
		Usage: "Manage registries",
		Commands: []*cli.Command{
			regList(),
		},
	}
}
