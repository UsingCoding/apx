package sandbox

import (
	"dario.cat/mergo"
	"github.com/pkg/errors"
)

func MergeSpec[T any](specs []T) (zero T, err error) {
	if len(specs) == 0 {
		return zero, errors.New("no specs to merge")
	}
	s := specs[0]

	for _, spec := range specs[1:] {
		err = mergo.Merge(&s, spec, mergo.WithOverride)
		if err != nil {
			return zero, errors.Wrap(err, "merge specs")
		}
	}

	return s, nil
}

func AssertSpecs[T any](specs []any) ([]T, error) {
	var zero T
	res := make([]T, 0, len(specs))
	for _, spec := range specs {
		s, ok := spec.(T)
		if !ok {
			return nil, errors.Errorf("spec %T is not %T", spec, zero)
		}
		res = append(res, s)
	}
	return res, nil
}
