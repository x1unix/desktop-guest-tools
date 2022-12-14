package vmci

import (
	"fmt"
	"syscall"
	"unsafe"
)

const listenBacklogSize = 100

// Workaround to call syscall.Bind with a custom syscall.Sockaddr implementation.
// Sockaddr interface methods are unexported, so it's impossible to implement it.
//
// Under the hood, syscall.Bind just calls Sockaddr.sockaddr to obtain a raw sockaddr
// pointer and pass it to bind() implementation.

//go:linkname syscallBind syscall.bind
func syscallBind(s syscall.Handle, name unsafe.Pointer, namelen int32) (err error)

type VSocketListener struct {
	socketFd syscall.Handle
	sockAddr *sockAddrVM
}

func (l VSocketListener) Close() error {
	return closeSocket(l.socketFd)
}

// Listen acts like net.Listen for vSocket sockets.
func Listen(port int) (*VSocketListener, error) {
	afVmci, err := GetAFValue()
	if err != nil {
		return nil, err
	}

	sockFd, err := syscall.Socket(afVmci, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open socket: %w", err)
	}

	addr := newSockAddr(saFamily(afVmci), uint32(port), VMAddrCIDAny)
	addrPtr, ptrSize := addr.sockaddr()
	if err := syscallBind(sockFd, addrPtr, ptrSize); err != nil {
		_ = closeSocket(sockFd)
		return nil, fmt.Errorf("socket bind failed: %w", err)
	}

	if err := syscall.Listen(sockFd, listenBacklogSize); err != nil {
		_ = closeSocket(sockFd)
		return nil, fmt.Errorf("socket listen failed: %w", err)
	}

	return &VSocketListener{
		socketFd: sockFd,
		sockAddr: addr,
	}, nil
}
