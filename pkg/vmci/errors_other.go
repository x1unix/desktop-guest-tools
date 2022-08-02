//go:build !windows

package vmci

// IsUnsupportedProtocolError checks whenever passed error
// was thrown because vSockets protocol is not registered.
//
// This function works only on Windows systems.
//
// See vmci.Listen function.
func IsUnsupportedProtocolError(err error) bool {
	// Unsupported
	return false
}
