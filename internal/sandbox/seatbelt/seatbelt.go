//go:build darwin

package seatbelt

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

//nolint:gochecknoinits
func init() {
	sandbox.R.Register("seatbelt", Seatbelt{})
}

type Seatbelt struct{}

func (s Seatbelt) Exec(ctx context.Context, cmd []string, p sandbox.Policy, logger *slog.Logger) error {
	// use absolute path to sandbox-exec to avoid path substitution attack
	const sb = "/usr/bin/sandbox-exec"

	profile, profileParams, err := makeProfile(p)
	if err != nil {
		return err
	}

	if logger.Enabled(ctx, slog.LevelDebug) {
		fmt.Println(profile)
	}

	args := []string{sb, "-p", profile}
	args = append(args, profileParams.flat()...)

	// separator before cmd
	args = append(args, "--")
	args = append(args, cmd...)

	logger.DebugContext(ctx, "", slog.Any("args", args))

	err = syscall.Exec(sb, args, os.Environ())
	return errors.Wrap(err, "exec seatbelt")
}
