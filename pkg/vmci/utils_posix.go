//go:build linux || darwin || freebsd || openbsd

package vmci

import (
	"errors"
	"syscall"
)

func closeSocket(fd syscall.Handle) error {
	if fd == syscall.Stdin {
		return errors.New("close: invalid socket handle")
	}

	return syscall.Close(fd)
}
