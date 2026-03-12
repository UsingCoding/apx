package container

import (
	"context"
	"fmt"
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

func Make(ctx context.Context, cmd *cli.Command, logger *slog.Logger) (Container, error) {
	regLocations := []fs.FS{
		registry.RegFS,
	}
	if f := locateLocalRegistry(ctx, cmd, logger); f != nil {
		regLocations = append(regLocations, f)
	}
	if f := locateLegacyLocalRegistry(ctx, logger); f != nil {
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

func locateLocalRegistry(ctx context.Context, cmd *cli.Command, logger *slog.Logger) fs.FS {
	s := cmd.String("base-dir")
	if s == "" {
		return nil
	}

	stat, err := os.Stat(s)

	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		logger.WarnContext(
			ctx,
			"base dir stat error",
			slog.String("base-dir", s),
			slog.Any("err", err),
		)
		return nil
	case !stat.IsDir():
		logger.WarnContext(
			ctx,
			"base dir is not a directory",
			slog.String("base-dir", s),
		)
		return nil
	default:
		return os.DirFS(s)
	}
}

func locateLegacyLocalRegistry(ctx context.Context, logger *slog.Logger) fs.FS {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}

	_, err = os.Stat(dir)
	if err != nil {
		// ignore any error
		return nil
	}

	msg := fmt.Sprintf("use legacy local registry at %s; please use new one", dir)
	logger.DebugContext(ctx, msg)

	return os.DirFS(dir)
}
