package shellenv

import (
	"fmt"
	"strings"

	"github.com/UsingCoding/apx/internal/app"
)

func bash(r app.Registry) (string, error) {
	const alias = `alias %s.apx="apx %s --"`

	aliases := make([]string, 0, len(r.All()))
	for _, apxtoml := range r.All() {
		aliases = append(aliases, fmt.Sprintf(alias, apxtoml.Name, apxtoml.Name))
	}

	return strings.Join(aliases, "\n") + "\n", nil
}
