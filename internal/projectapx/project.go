package projectapx

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/sandbox"
)

type Project struct {
	Policy sandbox.Policy `toml:"policy"`
}

const (
	defPath = "project.apx.toml"
)

func Decode() (p Project, err error) {
	_, err = toml.DecodeFile(defPath, &p)
	if err != nil {
		return Project{}, errors.Wrap(err, "decoding project.apx.toml")
	}
	p, err = expandPaths(p)
	if err != nil {
		return Project{}, errors.Wrap(err, "expanding paths in project.apx.toml")
	}

	return p, nil
}

func expandPaths(project Project) (Project, error) {
	expand := func(p string) (string, error) {
		p = os.ExpandEnv(p)
		p, err := filepath.Abs(p)
		if err != nil {
			return "", err
		}

		return p, nil
	}

	for i, p := range project.Policy.Filesystem.ROPaths {
		var err error
		p, err = expand(p)
		if err != nil {
			return Project{}, errors.Wrap(err, "roPaths")
		}
		project.Policy.Filesystem.ROPaths[i] = p
	}
	for i, p := range project.Policy.Filesystem.RWPaths {
		var err error
		p, err = expand(p)
		if err != nil {
			return Project{}, errors.Wrap(err, "rwPaths")
		}
		project.Policy.Filesystem.RWPaths[i] = p
	}
	for i, p := range project.Policy.Filesystem.DenyPaths {
		var err error
		p, err = expand(p)
		if err != nil {
			return Project{}, errors.Wrap(err, "denyPaths")
		}
		project.Policy.Filesystem.DenyPaths[i] = p
	}

	return project, nil
}
