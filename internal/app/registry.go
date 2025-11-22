package app

import (
	"io/fs"

	"github.com/pkg/errors"
)

type Registry struct {
	apps []APXTOML
}

func (r Registry) Find(id string) (APXTOML, error) {
	for _, app := range r.apps {
		if app.Name == id {
			return app, nil
		}
	}
	return APXTOML{}, errors.Errorf("app %q not found", id)
}

func (r Registry) All() []APXTOML {
	return r.apps
}

func LoadRegistry(sources []fs.FS) (Registry, error) {
	var apps []APXTOML
	for i, src := range sources {
		a, err := loadApps(src)
		if err != nil {
			return Registry{}, errors.Wrapf(err, "load apps from %d source", i)
		}
		apps = append(apps, a...)
	}
	return Registry{apps: apps}, nil
}

func loadApps(src fs.FS) (res []APXTOML, err error) {
	err = fs.WalkDir(src, ".", func(p string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		matched, err := matcher(p)
		if err != nil {
			return errors.Wrapf(err, "match %s", p)
		}
		if !matched {
			return nil
		}

		apx, err := decode(src, p)
		if err != nil {
			return errors.Wrapf(err, "decode %s", p)
		}

		res = append(res, apx)

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "search for apx.toml files")
	}

	return res, nil
}
