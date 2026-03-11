package landlock

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

type Landlock struct{}

func (l Landlock) Exec(ctx context.Context, cmd []string, p sandbox.Policy, logger *slog.Logger) error {
	err := applyLandlock(ctx, p, logger)
	if err != nil {
		return err
	}

	for k, v := range p.Env {
		_ = os.Setenv(k, v)
	}

	bin, err := exec.LookPath(cmd[0])
	if err != nil {
		return errors.Wrap(err, "resolve command path")
	}

	logger.DebugContext(
		ctx,
		"exec with landlock",
		slog.Any("cmd", cmd),
		slog.Any("env", p.Env),
	)

	err = syscall.Exec(bin, cmd, os.Environ())
	return errors.Wrap(err, "exec with landlock")
}
