package seatbelt

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"syscall"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

type Seatbelt struct{}

func (s Seatbelt) Exec(ctx context.Context, cmd []string, p sandbox.Policy, logger *slog.Logger) error {
	// use absolute path to sandbox-exec to avoid path substitution attack
	const sb = "/usr/bin/sandbox-exec"

	p = s.appendROPaths(p)

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

	logger.DebugContext(
		ctx,
		"",
		slog.Any("args", args),
		slog.Any("env", p.Env),
	)

	for k, v := range p.Env {
		_ = os.Setenv(k, v)
	}

	err = syscall.Exec(sb, args, os.Environ())
	return errors.Wrap(err, "exec seatbelt")
}

// Append RW to RO paths, so rw paths can be read
// At darwin world - file-write* only gives write rule
func (s Seatbelt) appendROPaths(p sandbox.Policy) sandbox.Policy {
	for _, rwPath := range p.Filesystem.RWPaths {
		if !slices.Contains(p.Filesystem.ROPaths, rwPath) {
			p.Filesystem.ROPaths = append(p.Filesystem.ROPaths, rwPath)
		}
	}
	return p
}
