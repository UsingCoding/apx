package container

import (
	"context"
	"io/fs"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/UsingCoding/apx/internal/app"
	"github.com/UsingCoding/apx/registry"
)

type Container struct {
	Logger      *slog.Logger
	ApxRegistry app.Registry
}

func Make(_ context.Context, cmd *cli.Command, logger *slog.Logger) (Container, error) {
	regLocations := []fs.FS{
		registry.RegFS,
	}
	if f := locateLocalRegistry(cmd, logger); f != nil {
		regLocations = append(regLocations, f)
	}
	reg, err := app.LoadRegistry(regLocations)
	if err != nil {
		return Container{}, errors.Wrap(err, "load registry")
	}

	return Container{
		Logger:      logger,
		ApxRegistry: reg,
	}, nil
}

func locateLocalRegistry(cmd *cli.Command, logger *slog.Logger) fs.FS {
	s := cmd.String("base-dir")
	if s == "" {
		return nil
	}

	stat, err := os.Stat(s)

	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		logger.Warn(
			"base dir stat error",
			slog.String("base-dir", s),
			slog.Any("err", err),
		)
		return nil
	case !stat.IsDir():
		logger.Warn(
			"base dir is not a directory",
			slog.String("base-dir", s),
		)
		return nil
	default:
		return os.DirFS(s)
	}
}
