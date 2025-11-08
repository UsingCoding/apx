package app

import (
	"io/fs"
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
	return apx, errors.Wrap(err, "error decoding APXTOML")
}

func matcher(p string) (matched bool, err error) {
	return filepath.Match("*.apx.toml", p)
}
