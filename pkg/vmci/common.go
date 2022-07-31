package vmci

import "math"

const (
	// SoVMCIBufferSize is option name for STREAM socket buffer size.
	//
	// Use as the option name in setsockopt or getsockopt(3) to set
	// or get an uint64 that specifies the size of the
	// buffer underlying a vSockets STREAM socket.
	//
	// Value is clamped to the MIN and MAX.
	SoVMCIBufferSize = 0

	// SoVMCIBufferMinSize is option name for STREAM socket maximum buffer size.
	SoVMCIBufferMinSize = 1

	// SoVMCIBufferMaxSize is option name for STREAM socket maximum buffer size.
	SoVMCIBufferMaxSize = 2

	// SoVMCIPeerHostVmId is option name for socket peer's host-specific VM ID.
	SoVMCIPeerHostVmId = 3

	// SoVMCIServiceLabel is option name for socket's service label.
	SoVMCIServiceLabel = 4

	// SoVMCITrusted is option name for determining if a socket is trusted.
	SoVMCITrusted = 5

	// SoVMCIConnectTimeout is option name for STREAM socket connection timeout.
	SoVMCIConnectTimeout = 6

	// SoVMCINonblockTxrx is option name for using non-blocking send/receive.
	SoVMCINonblockTxrx = 7
)

const (
	// VMAddrCIDAny is vSocket equivalent of INADDR_ANY.
	VMAddrCIDAny = uint32(math.MaxUint32)

	// VMAddrPortAny is bind to any available port.
	VMAddrPortAny = uint32(math.MaxUint32)

	// SocketsInvalidVersion is invalid vSockets version.
	SocketsInvalidVersion = uint32(math.MaxUint32)
)

func omitError(_ error) {}

/**
#if defined(_WIN32) || defined(VMKERNEL)
   typedef unsigned short sa_family_t;
#endif // _WIN32

#if defined(VMKERNEL)
   struct sockaddr {
      sa_family_t sa_family;
      char sa_data[14];
   };
#endif
*/

type saFamily = uint16

type uuid = [128]byte
type uuid2cid struct {
	contextID  uint32
	pad        uint32
	uuidString uuid
}
