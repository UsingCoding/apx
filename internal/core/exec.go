package core

import (
	"context"
	"log/slog"
	"os"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/app"
	"github.com/UsingCoding/apx/internal/projectapx"
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

	sndbox, ok := sandbox.R.Lookup(s.ID)
	if !ok {
		return errors.Errorf("sandbox %q for %q not found", s.ID, argv0)
	}

	var p projectapx.Project
	p, err = projectapx.Decode()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "load project")
	}

	specs := []any{}

	// project found
	if !errors.Is(err, os.ErrNotExist) {
		e.Logger.Debug("load project", slog.Any("project", p))

	}

	sndbox.Runtime.Exec(ctx, e.CMD)
	return sndbox.Exec(ctx, e.CMD, s.Policy, e.Logger)
}
