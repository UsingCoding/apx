package registry

import (
	"embed"
)

//go:embed *.toml
var RegFS embed.FS
