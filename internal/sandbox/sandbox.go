package sandbox

import (
	"context"
	"log/slog"
)

type Sandbox struct {
	ID   string
	Spec any

	Runtime Runtime
}

type Runtime interface {
	Exec(ctx context.Context, cmd []string, specs []any, logger *slog.Logger) error
}
