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
	Sandboxes []sandbox.Sandbox `toml:"sandboxes"`

	AdditionalProperties toml.MetaData `toml:"-"`
}

type Sandbox struct {
	ID   string `toml:"id"`
	Spec any    `toml:"spec"`
}

func decode(src fs.FS, p string) (apx APXTOML, err error) {
	md, err := toml.DecodeFS(src, p, &apx)
	if err != nil {
		return APXTOML{}, errors.Wrap(err, "error decoding APXTOML")
	}
	apx.AdditionalProperties = md
	return apx, err
}

func matcher(p string) (matched bool, err error) {
	return filepath.Match("*.apx.toml", p)
}
