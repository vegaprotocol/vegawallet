package wallet

const (
	Version1 = uint32(1)
	// Version2 identifies HD wallet v2.
	Version2 = uint32(2)
	// LatestVersion is the latest version of Vega's HD wallet. Created wallets
	// are always pointing to the latest version.
	LatestVersion = Version2
)

// SupportedVersions list versions supported by Vega's HD wallet.
var SupportedVersions = []uint32{Version1, Version2}

func IsVersionSupported(v uint32) bool {
	return v == Version1 || v == Version2
}
