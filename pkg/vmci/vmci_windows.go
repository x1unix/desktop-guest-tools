package vmci

import (
	"fmt"
	"math"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Source code adapted from vmci-socket.h
// Source: https://github.com/vmware/open-vm-tools/blob/master/open-vm-tools/lib/include/vmci_sockets.h

type controlCode = uint32

const (
	ioctlVMCISocketsVersion     = controlCode(0x81032058)
	ioctlVMCISocketsGetAFValue  = controlCode(0x81032068)
	ioctlVMCISocketsGetLocalCID = controlCode(0x8103206c)
	ioctlVMCISocketsUUID2CID    = controlCode(0x810320a4)
)

const (
	socketsDevicePath = `\\.\VMCI`

	// svmZeroSize value is based on substraction of sizeof(s):
	// sizeof(sockaddr - sa_family_t - uint16 - uint32 - uint32)
	svmZeroSize = 4
)

// sockAddrVM is sockaddr equivalent for VMCI.
type sockAddrVM struct {
	family    saFamily
	reserved1 uint16
	port      uint32
	cid       uint32
	zero      [svmZeroSize]byte
}

func (sa *sockAddrVM) sockaddr() (unsafe.Pointer, int32) {
	size := unsafe.Sizeof(*sa)
	return unsafe.Pointer(sa), int32(size)
}

func newSockAddr(family saFamily, port, cid uint32) *sockAddrVM {
	return &sockAddrVM{
		family: family,
		cid:    cid,
		port:   port,
	}
}

func openSocketDevice() (windows.Handle, error) {
	lpSocketAddr, err := windows.UTF16PtrFromString(socketsDevicePath)
	if err != nil {
		return 0, err
	}

	hDevice, err := windows.CreateFile(lpSocketAddr, windows.GENERIC_READ, 0, nil,
		windows.OPEN_EXISTING, windows.FILE_FLAG_OVERLAPPED, 0)
	if err != nil {
		return 0, fmt.Errorf("vmci: failed to open VMCI device: %w", err)
	}

	if hDevice == windows.InvalidHandle {
		return 0, ErrInvalidHandle
	}

	return hDevice, nil
}

func deviceIOControl(cmd controlCode) (uint32, error) {
	hDevice, err := openSocketDevice()
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = windows.CloseHandle(hDevice)
	}()

	// overflow trick is used very often in original source code as:
	// 	unsigned int val = (unsigned int)-1;
	val := uint32(math.MaxUint32)

	// dirty trick to pass val as *byte and mitigate type signature mismatch
	valPtr := (*byte)(unsafe.Pointer(&val))
	valSize := uint32(unsafe.Sizeof(val))

	var ioReturn uint32
	err = windows.DeviceIoControl(hDevice, cmd, valPtr, valSize, valPtr, valSize, &ioReturn, nil)
	if err != nil {
		return 0, newDeviceIOControlError(cmd, err)
	}

	return val, nil
}

func resultAsInt(val uint32, err error) (int, error) {
	i32 := (int32)(val)
	return int(i32), err
}

// Version retrieves the vSockets version.
//
// See Version structure.
func Version() (VersionNumber, error) {
	val, err := deviceIOControl(ioctlVMCISocketsVersion)
	if err != nil {
		return 0, err
	}

	if val == socketsInvalidVersion {
		return 0, ErrInvalidVersion
	}

	return VersionNumber(val), err
}

// GetAFValue retrieves the address family value for vSockets.
//
// Returns the value to be used for the vSockets address family.
// This value should be used as the domain argument to socket() (when
// you might otherwise use AF_INET).  For vSocket-specific options,
// this value should also be used for the level argument to
// setsockopt (when you might otherwise use SOL_TCP).
//
// This function leaves its descriptor to the vsock device open so
// that the socket implementation knows that the socket family is still in
// use.  This is done because the address family is registered with the
// kernel on-demand and a notification is needed to unregister the address
// family.
//
// Use of this function is thus discouraged; please use GetAFValueFd() instead.
func GetAFValue() (int, error) {
	return resultAsInt(deviceIOControl(ioctlVMCISocketsGetAFValue))
}

// GetAFValueFd retrieves the address family value for vSockets.
//
// Returns the value to be used for the vSockets address family.
// This value should be used as the domain argument to socket() (when
// you might otherwise use AF_INET).  For vSocket-specific options,
// this value should also be used for the level argument to
// setsockopt (when you might otherwise use SOL_TCP).
//
// Receives a file descriptor to the VMCI device. The address family
// value is valid until this descriptor is closed. This parameter is
// only valid if the return value is not -1.
//
// Call ReleaseAFValueFd() to close this descriptor.
func GetAFValueFd(_ *int) (int, error) {
	// Unused, see: vmci_sockets.h:444
	return GetAFValue()
}

// ReleaseAFValueFd releases file descriptor obtained when retrieving
// the address family value.
//
// Accepts file descriptor to the VMCI device.
//
// Use this to release the file descriptor obtained by calling GetAFValueFd().
func ReleaseAFValueFd(_ int) {
	// Unused, see: vmci_sockets.h:447
}

// GetLocalCID retrieves the current context ID.
//
// Returns ErrContextUnavailable when current context ID is not available.
func GetLocalCID() (ContextID, error) {
	result, err := deviceIOControl(ioctlVMCISocketsGetLocalCID)
	if err != nil {
		return 0, err
	}

	return result, checkContextID(result)
}

// UUID2ContextID retrieves the context ID of a running VM, given a VM's UUID.
//
// Retrieves the context ID of a running virtual machine given that virtual
// machine's unique identifier.  The identifier is local to the host and
// its meaning is platform-specific.  On ESX, which is currently the only
// supported platform, it is the "bios.uuid" field as specified in the VM's
// VMX file.
//
// Returns ErrContextUnavailable when context is not available.
//
func UUID2ContextID(uuid string) (ContextID, error) {
	hDevice, err := openSocketDevice()
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = windows.CloseHandle(hDevice)
	}()

	io := uuid2cid{
		contextID:  VMAddrCIDAny,
		uuidString: uuidFromString(uuid),
	}
	// dirty trick to pass val as *byte and mitigate type signature mismatch
	ioPtr := (*byte)(unsafe.Pointer(&io))
	ioSize := uint32(unsafe.Sizeof(io))

	var ioReturn uint32
	err = windows.DeviceIoControl(hDevice, ioctlVMCISocketsUUID2CID, ioPtr, ioSize, ioPtr, ioSize, &ioReturn, nil)
	if err != nil {
		return 0, newDeviceIOControlError(ioctlVMCISocketsUUID2CID, err)
	}

	return io.contextID, checkContextID(io.contextID)
}
