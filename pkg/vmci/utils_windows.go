package vmci

import "syscall"

func closeSocket(fd syscall.Handle) error {
	return syscall.Closesocket(fd)
}
