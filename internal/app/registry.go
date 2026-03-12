package app

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Registry struct {
	apps []RegWrapper
}

// RegWrapper wrapper for registry specific metadata
type RegWrapper struct {
	APXTOML APXTOML

	Source string
}

func (r Registry) Find(id string) (APXTOML, error) {
	for _, app := range r.apps {
		if app.APXTOML.Name == id {
			return app.APXTOML, nil
		}
	}
	return APXTOML{}, errors.Errorf("app %q not found", id)
}

func (r Registry) All() []APXTOML {
	return lo.Map(r.apps, func(a RegWrapper, _ int) APXTOML {
		return a.APXTOML
	})
}

func (r Registry) AllWraps() []RegWrapper {
	return r.apps
}
