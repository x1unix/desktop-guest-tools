package vmci

import (
	"errors"
	"syscall"

	"golang.org/x/sys/windows"
)

// IsUnsupportedProtocolError checks whenever passed error
// was thrown because vSockets protocol is not registered.
//
// See vmci.Listen function.
func IsUnsupportedProtocolError(err error) bool {
	if err == nil {
		return false
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		unwrapped = err
	}

	errno, ok := unwrapped.(syscall.Errno)
	if !ok {
		return false
	}

	return errno == windows.WSAEAFNOSUPPORT
}
