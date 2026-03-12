package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/UsingCoding/apx/internal/app"
)

type RegList struct {
	Reg app.Registry
	W   io.Writer
}

func (s RegList) Do() error {
	for _, a := range s.Reg.AllWraps() {
		data, err := app.Encode(a.APXTOML)
		if err != nil {
			return err
		}

		sep := strings.Repeat("=", 30)
		_, _ = fmt.Fprintln(s.W, sep)

		msg := fmt.Sprintf("%s: %s", a.APXTOML.Name, a.Source)
		_, _ = fmt.Fprintln(s.W, msg)
		// content of apx.toml
		// marshaled toml already contains new line
		_, _ = fmt.Fprint(s.W, string(data))

		_, _ = fmt.Fprintln(s.W, sep)
	}

	return nil
}
