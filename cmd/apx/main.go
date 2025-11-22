package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	zeroslog "github.com/samber/slog-zerolog"
	"github.com/urfave/cli/v3"
)

const (
	appID = "apx"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	ctx := context.Background()

	ctx = subscribeForKillSignals(ctx)

	err := runApp(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = ".apx"
	}

	c := &cli.Command{
		Name:    appID,
		Version: version,
		// do not use built-in version flag
		HideVersion:           true,
		Usage:                 "CLI wrapper for platform-specific sandboxes",
		EnableShellCompletion: true,
		Action:                exec,
		Commands: []*cli.Command{
			shellenv(),
			versionCMD(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
			},
			&cli.StringFlag{
				Name:    "base-dir",
				Aliases: []string{"b"},
				Usage:   "Dir for config and local registry",
				Value:   path.Join(userConfigDir, "apx"),
			},
		},
	}

	return c.Run(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
		}
	}()

	return ctx
}

func logger(cmd *cli.Command) *slog.Logger {
	level := zerolog.InfoLevel
	leveler := slog.LevelInfo
	if cmd.Bool("verbose") {
		level = zerolog.DebugLevel
		leveler = slog.LevelDebug
	}

	w := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}

	zerologL := zerolog.New(w).Level(level)

	opts := zeroslog.Option{
		Logger: &zerologL,
		Level:  leveler,
	}
	handler := opts.NewZerologHandler()
	return slog.New(handler)
}
