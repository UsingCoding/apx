package sandbox

type Policy struct {
	Filesystem Filesystem `toml:"fs"`
	Network    Network    `toml:"net"`
}

type Network struct {
	Deny bool `toml:"deny"`
}

type Filesystem struct {
	FullDiskReadAccess bool     `toml:"fullDiskReadAccess"`
	NoCache            bool     `toml:"noCache"`
	ROPaths            []string `toml:"roPaths"`
	RWPaths            []string `toml:"rwPaths"`
}
