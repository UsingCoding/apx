package osconfig

import (
	"os"

	"github.com/pkg/errors"
)

// UserConfigDir follows unix path from os.UserConfigDir
// That achieves:
// * Unix compatible config location
// * Do not use native darwin config path
func UserConfigDir() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
		}
		dir += "/.config"
	}
	return dir, nil
}
