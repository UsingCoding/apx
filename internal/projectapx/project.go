package projectapx

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

type Project struct {
	Sandboxes []sandbox.Sandbox `toml:"sandboxes"`
}

const (
	defPath = "project.apx.toml"
)

func Decode() (p Project, err error) {
	_, err = toml.DecodeFile(defPath, &p)
	if err != nil {
		return Project{}, errors.Wrap(err, "decoding project.apx.toml")
	}

	return p, nil
}
