package core

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/app"
	"github.com/UsingCoding/apx/internal/sandbox"
)

type Exec struct {
	CMD []string
	Reg app.Registry

	Logger *slog.Logger
}

func (e Exec) Do(ctx context.Context) error {
	if len(e.CMD) == 0 {
		return errors.New("no command specified")
	}

	argv0 := e.CMD[0]

	apxtoml, err := e.Reg.Find(argv0)
	if err != nil {
		return err
	}

	// for now, just take first sandbox
	s := apxtoml.Sandboxes[0]

	sndbox, ok := sandbox.R.Lookup(s.Type)
	if !ok {
		return errors.Errorf("sandbox %q for %q not found", s.Type, argv0)
	}

	return sndbox.Exec(ctx, e.CMD, s.Policy, e.Logger)
}
