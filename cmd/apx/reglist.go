package main

import (
	"context"
	"os"

	"github.com/UsingCoding/apx/internal/container"
	"github.com/UsingCoding/apx/internal/core"

	"github.com/urfave/cli/v3"
)

func regList() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List registries",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c, err := container.Make(ctx, cmd, logger(cmd))
			if err != nil {
				return err
			}

			return core.RegList{
				Reg: c.ApxRegistry,
				W:   os.Stdout,
			}.Do()
		},
	}
}
