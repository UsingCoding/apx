package sandbox

import (
	"context"
	"log/slog"
)

type Sandbox interface {
	Exec(ctx context.Context, cmd []string, policy Policy, logger *slog.Logger) error
}
