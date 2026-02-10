package seatbelt

import "os"

type Policy struct {
	Env        Env        `toml:"env"`
	Filesystem Filesystem `toml:"fs"`
	Network    Network    `toml:"net"`
}

type Env map[string]string

type Filesystem struct {
	FullDiskReadAccess bool `toml:"fullDiskReadAccess"`
	NoCache            bool `toml:"noCache"`

	ROPaths   []string `toml:"roPaths"`
	RWPaths   []string `toml:"rwPaths"`
	DenyPaths []string `toml:"denyPaths"`
}

type Network struct {
	Deny bool `toml:"deny"`
}

func expandPaths(pol Policy) Policy {
	for i, p := range pol.Filesystem.ROPaths {
		pol.Filesystem.ROPaths[i] = os.ExpandEnv(p)
	}
	for i, p := range pol.Filesystem.RWPaths {
		pol.Filesystem.RWPaths[i] = os.ExpandEnv(p)
	}
	for i, p := range pol.Filesystem.DenyPaths {
		pol.Filesystem.DenyPaths[i] = os.ExpandEnv(p)
	}

	return pol
}
