package seatbelt

import (
	"github.com/UsingCoding/apx/internal/sandbox"
)

//nolint:gochecknoinits
func init() {
	sandbox.R.Register("seatbelt", Seatbelt{})
}
