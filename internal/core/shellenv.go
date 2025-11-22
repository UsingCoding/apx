package core

import (
	"os"

	"github.com/pkg/errors"

	"github.com/UsingCoding/apx/internal/app"
	"github.com/UsingCoding/apx/internal/shellenv"
)

type Shellenv struct {
	Shell string
	Reg   app.Registry
}

func (s Shellenv) Do() error {
	env, err := shellenv.Env{}.Generate(s.Shell, s.Reg)
	if err != nil {
		return errors.Wrap(err, "generate shellenv")
	}

	_, _ = os.Stdout.WriteString(env)

	return nil
}
