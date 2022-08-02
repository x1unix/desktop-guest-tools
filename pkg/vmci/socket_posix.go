//go:build linux || darwin || freebsd || openbsd

package vmci

import (
	"errors"
	"syscall"
)

func (l VSocketListener) Close() error {
	if l.socketFd == syscall.Stdin {
		return errors.New("close: invalid socket handle")
	}

	return syscall.Close(l.socketFd)
}
