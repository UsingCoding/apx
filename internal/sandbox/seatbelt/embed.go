package seatbelt

import _ "embed"

// Provide embedded policy blobs for non-darwin builds too so profile generation
// can be tested without macOS-specific files being excluded by build tags.

var (
	//go:embed seatbelt_base_policy.sbpl
	seatbeltBasePolicy []byte
	//go:embed seatbelt_network_policy.sbpl
	seatbeltNetworkPolicy []byte
)
