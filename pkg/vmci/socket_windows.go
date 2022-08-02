package vmci

import "syscall"

func (l VSocketListener) Close() error {
	return syscall.Closesocket(l.socketFd)
}
