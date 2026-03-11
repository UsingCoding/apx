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

	s, sndbox, err := e.findSandbox(apxtoml.Sandboxes)
	if err != nil {
		return errors.Wrapf(err, "for %q", argv0)
	}

	var p projectapx.Project
	p, err = projectapx.Decode()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "load project")
	}
	// project found
	if !errors.Is(err, os.ErrNotExist) {
		e.Logger.Debug("load project", slog.Any("project", p))

		s.Policy, err = sandbox.MergePolicies(s.Policy, p.Policy)
		if err != nil {
			return err
		}
	}

	e.Logger.Debug("policy", slog.Any("policy", s.Policy))
	return sndbox.Exec(ctx, e.CMD, s.Policy, e.Logger)
}

func (e Exec) findSandbox(sandboxes []app.Sandbox) (app.Sandbox, sandbox.Sandbox, error) {
	for _, s := range sandboxes {
		// for now, just take first supported sandbox
		if b, ok := sandbox.R.Lookup(s.Type); ok {
			return s, b, nil
		}
	}

	return app.Sandbox{}, nil, errors.New("supported sandboxes not found")
}
