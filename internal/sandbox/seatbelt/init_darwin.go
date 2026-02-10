package seatbelt

import (
	"github.com/UsingCoding/apx/internal/sandbox"
)

//nolint:gochecknoinits
func init() {
	sandbox.R.Register(sandbox.Sandbox{
		ID:      "seatbelt",
		Spec:    Policy{},
		Runtime: Seatbelt{},
	})
}
