package sandbox

import (
	"dario.cat/mergo"
	"github.com/pkg/errors"
)

type Policy struct {
	Filesystem Filesystem `toml:"fs"`
	Network    Network    `toml:"net"`
}

type Network struct {
	Deny bool `toml:"deny"`
}

type Filesystem struct {
	FullDiskReadAccess bool `toml:"fullDiskReadAccess"`
	NoCache            bool `toml:"noCache"`

	ROPaths   []string `toml:"roPaths"`
	RWPaths   []string `toml:"rwPaths"`
	DenyPaths []string `toml:"denyPaths"`
}

func MergePolicies(p1, p2 Policy) (Policy, error) {
	err := mergo.Merge(&p1, p2, mergo.WithOverride)
	if err != nil {
		return Policy{}, errors.Wrap(err, "merge policies")
	}
	return p1, nil
}
