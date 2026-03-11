package landlock

import (
	"context"
	"log/slog"
	"os"

	ll "github.com/landlock-lsm/go-landlock/landlock"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/UsingCoding/apx/internal/sandbox"
)

var (
	// systemROPaths are paths that must be readable for basic process execution.
	// These are always granted read-only access.
	systemROPaths = []string{
		"/usr",
		"/bin",
		"/sbin",
		"/lib",
		"/lib64",
		"/etc",
		"/proc",
		"/run",
		"/sys",
	}
	// systemRWPaths are paths that require read-write access for basic process
	// execution
	systemRWPaths = []string{
		"/var",
		"/tmp",
		"/dev",
	}
)

func applyLandlock(ctx context.Context, p sandbox.Policy, logger *slog.Logger) error {
	rules, err := fsRules(ctx, p, logger)
	if err != nil {
		return errors.Wrap(err, "make rules for FS")
	}

	cfg := ll.V5.BestEffort()

	err = cfg.RestrictPaths(lo.Map(
		rules,
		func(r ll.FSRule, _ int) ll.Rule {
			return r
		},
	)...)
	if err != nil {
		return errors.Wrap(err, "restrict fs")
	}

	if p.Network.Deny {
		// Landlock (V4+) can only restrict TCP bind/connect. UDP and other
		// protocols are not covered by Landlock's network restriction.
		// Calling RestrictNet with no rules blocks all TCP bind and connect.
		err = cfg.RestrictNet()
		if err != nil {
			return errors.Wrap(err, "restrict net")
		}
	}

	logger.DebugContext(ctx, "landlock applied", slog.String("config", cfg.String()))

	return nil
}

func fsRules(ctx context.Context, p sandbox.Policy, logger *slog.Logger) ([]ll.FSRule, error) {
	p.Filesystem.ROPaths = append(p.Filesystem.ROPaths, systemROPaths...)
	p.Filesystem.RWPaths = append(p.Filesystem.RWPaths, systemRWPaths...)

	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "current working directory")
	}

	p.Filesystem.RWPaths = append(p.Filesystem.RWPaths, wd)

	var rules []ll.FSRule

	lo.ForEach(p.Filesystem.ROPaths, func(p string, _ int) {
		// ignore err, if path not exists or incorrect - skip it
		stat, _ := os.Stat(p)
		if stat != nil {
			switch stat.IsDir() {
			case true:
				rules = append(rules, ll.RODirs(p))
			default:
				rules = append(rules, ll.ROFiles(p))
			}
		}
	})
	lo.ForEach(p.Filesystem.RWPaths, func(p string, _ int) {
		// ignore err, if path not exists or incorrect - skip it
		stat, _ := os.Stat(p)
		if stat != nil {
			switch stat.IsDir() {
			case true:
				rules = append(rules, ll.RWDirs(p))
			default:
				rules = append(rules, ll.RWFiles(p))
			}
		}
	})
	if len(p.Filesystem.DenyPaths) > 0 {
		logger.DebugContext(
			ctx,
			"landlock: DenyPaths are naturally excluded from the allowlist",
			slog.Any("denyPaths", p.Filesystem.DenyPaths),
		)
	}

	if p.Filesystem.FullDiskReadAccess {
		rules = append(rules, ll.RODirs("/"))
	}

	return rules, nil
}
