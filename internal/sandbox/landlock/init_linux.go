package landlock

import "github.com/UsingCoding/apx/internal/sandbox"

//nolint:gochecknoinits
func init() {
	sandbox.R.Register("landlock", Landlock{})
}
