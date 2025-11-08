package main

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/UsingCoding/apx/internal/container"
	"github.com/UsingCoding/apx/internal/core"
)

func exec(ctx context.Context, cmd *cli.Command) error {
	c, err := container.Make(ctx, cmd, logger(cmd))
	if err != nil {
		return err
	}

	return core.Exec{
		Reg: c.ApxRegistry,
		// interpret all args as cmd to be sandboxed
		CMD:    cmd.Args().Slice(),
		Logger: logger(cmd),
	}.Do(ctx)
}
