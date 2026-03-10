package sandbox

import (
	"os"
	"path/filepath"
	"slices"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Policy struct {
	Env        Env        `toml:"env"`
	Filesystem Filesystem `toml:"fs"`
	Network    Network    `toml:"net"`
}

func (p *Policy) UnmarshalTOML(a any) error {
	data, err := toml.Marshal(a)
	if err != nil {
		return err
	}

	// HACK: to avoid recursive self unmarshall
	type poli Policy
	var pol poli
	err = toml.Unmarshal(data, &pol)
	if err != nil {
		return err
	}

	res, err := decode(Policy(pol))
	if err != nil {
		return err
	}

	*p = res
	return nil
}

type Env map[string]string

type Filesystem struct {
	FullDiskReadAccess bool `toml:"fullDiskReadAccess"`
	NoCache            bool `toml:"noCache"`

	// Home control access to $HOME dir
	// if set - home available
	Home *Home `toml:"home"`

	ROPaths   []string `toml:"roPaths"`
	RWPaths   []string `toml:"rwPaths"`
	DenyPaths []string `toml:"denyPaths"`
}

type Home struct {
	// AllPaths allows access to any subpath under Home dir
	AllPaths bool `toml:"allPaths"`
	// SkipDefaultDenyList do not use default deny list
	SkipDefaultDenyList bool `toml:"skipDefaultDenyList"`
	// DenyList under home dir where access is forbidden
	DenyList []string `toml:"denyList"`
	// RW makes access to $HOME rw
	RW bool `toml:"rw"`
}

type Network struct {
	Deny bool `toml:"deny"`
}

func decode(p Policy) (Policy, error) {
	var err error
	p.Filesystem.ROPaths, err = expandPaths(p.Filesystem.ROPaths)
	if err != nil {
		return Policy{}, errors.Wrap(err, "ro paths")
	}
	p.Filesystem.RWPaths, err = expandPaths(p.Filesystem.RWPaths)
	if err != nil {
		return Policy{}, errors.Wrap(err, "rw paths")
	}
	p.Filesystem.DenyPaths, err = expandPaths(p.Filesystem.DenyPaths)
	if err != nil {
		return Policy{}, errors.Wrap(err, "deny paths")
	}

	if p.Filesystem.Home != nil {
		p, err = homeDirPolicy(p, *p.Filesystem.Home)
		if err != nil {
			return Policy{}, errors.Wrap(err, "home policy")
		}
	}

	return p, nil
}

func MergePolicies(p1, p2 Policy) (Policy, error) {
	err := mergo.Merge(&p1, p2, mergo.WithOverride)
	if err != nil {
		return Policy{}, errors.Wrap(err, "merge policies")
	}
	return p1, nil
}

func expandPaths(paths []string) ([]string, error) {
	expand := func(p string) (string, error) {
		p = os.ExpandEnv(p)
		p, err := filepath.Abs(p)
		if err != nil {
			return "", err
		}

		return p, nil
	}

	for i, p := range paths {
		var err error
		p, err = expand(p)
		if err != nil {
			return nil, errors.Wrapf(err, "expanding path %d", i)
		}
		paths[i] = p
	}
	return paths, nil
}

func homeDirPolicy(p Policy, h Home) (Policy, error) {
	homeDir, err2 := os.UserHomeDir()
	if err2 != nil {
		return Policy{}, errors.Wrap(err2, "get home dir")
	}

	entries, err2 := os.ReadDir(homeDir)
	if err2 != nil {
		return Policy{}, errors.Wrap(err2, "read home dir")
	}

	// full access to $HOME
	if h.AllPaths {
		switch h.RW {
		case true:
			p.Filesystem.RWPaths = append(p.Filesystem.RWPaths, homeDir)
		default:
			p.Filesystem.ROPaths = append(p.Filesystem.ROPaths, homeDir)
		}
		return p, nil
	}

	var denyList []string
	if !h.SkipDefaultDenyList {
		def := []string{
			".ssh",
			".kube",
		}
		denyList = append(denyList, def...)
	}
	if h.DenyList != nil {
		denyList = append(denyList, h.DenyList...)
	}

	for _, e := range entries {
		if slices.Contains(denyList, e.Name()) {
			continue
		}

		switch h.RW {
		case true:
			p.Filesystem.RWPaths = append(p.Filesystem.RWPaths, filepath.Join(homeDir, e.Name()))
		default:
			p.Filesystem.ROPaths = append(p.Filesystem.ROPaths, filepath.Join(homeDir, e.Name()))
		}
	}

	return p, nil
}
