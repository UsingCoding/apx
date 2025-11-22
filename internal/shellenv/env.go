package shellenv

import (
	"fmt"

	"github.com/UsingCoding/apx/internal/app"
)

type Env struct{}

func (e Env) Generate(shell string, r app.Registry) (string, error) {
	switch shell {
	case "bash",
		// in shellenv - zsh fully compatible with bash, so fall to bash implementation
		"zsh":
		return bash(r)
	default:
		return "", fmt.Errorf("unsupported shell: %q", shell)
	}
}
