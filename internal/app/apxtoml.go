package app

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

// APXTOML format for apps for apx
type APXTOML struct {
	Name string `toml:"name"`
	// Multiple sandboxes allowed for app
	Sandboxes []Sandbox `toml:"sandboxes"`
}

type Sandbox struct {
	Type   string         `toml:"type"`
	Policy sandbox.Policy `toml:"policy"`
}

func decode(src fs.FS, p string) (apx APXTOML, err error) {
	_, err = toml.DecodeFS(src, p, &apx)
	if err != nil {
		return APXTOML{}, errors.Wrap(err, "error decoding APXTOML")
	}
	apx = expandPaths(apx)
	return apx, err
}

func matcher(p string) (matched bool, err error) {
	return filepath.Match("*.apx.toml", p)
}

func expandPaths(apx APXTOML) APXTOML {
	for _, s := range apx.Sandboxes {
		for i, p := range s.Policy.Filesystem.ROPaths {
			s.Policy.Filesystem.ROPaths[i] = os.ExpandEnv(p)
		}
		for i, p := range s.Policy.Filesystem.RWPaths {
			s.Policy.Filesystem.RWPaths[i] = os.ExpandEnv(p)
		}
	}

	return apx
}
