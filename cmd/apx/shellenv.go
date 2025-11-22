package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/UsingCoding/apx/internal/container"
	"github.com/UsingCoding/apx/internal/core"
)

func shellenv() *cli.Command {
	return &cli.Command{
		Name:  "shellenv",
		Usage: "apx shellenv <SHELL> - generates apx helpers for bash,zsh",
		UsageText: `apx shellenv <SHELL>
Supported <SHELL>: bash, zsh

Usage bash:
echo 'eval "$(apx shellenv bash)"' >> ~/.bashrc

Usage zsh:
echo 'eval "$(apx shellenv zsh)"' >> ~/.zshrc

Helpers:
For each app in registry it generates shortener like
codex.apx for codex
opencode.apx for opencode
`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errors.New("shellenv requires 1 argument")
			}

			shell := cmd.Args().First()

			c, err := container.Make(ctx, cmd, logger(cmd))
			if err != nil {
				return err
			}

			return core.Shellenv{
				Shell: shell,
				Reg:   c.ApxRegistry,
			}.Do()
		},
	}
}
