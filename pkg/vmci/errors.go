package vmci

import (
	"errors"
	"fmt"
)

var (
	// ErrContextUnavailable error occurs when requested context ID is not available.
	//
	// See GetLocalCID function.
	ErrContextUnavailable = errors.New("vmci: context ID is not available")

	// ErrInvalidHandle occurs when VMCI device file could not be opened.
	ErrInvalidHandle = errors.New("vmci: invalid handle")

	// ErrInvalidVersion occurs when returned vSockets version is invalid.
	ErrInvalidVersion = errors.New("vmci: invalid version")
)

// DeviceIOControlError is deviceIOControl error.
//
// Error occurs when invalid control code was supplied or ioctl operation failed.
type DeviceIOControlError struct {
	ControlCode controlCode
	InnerError  error
}

func newDeviceIOControlError(controlCode controlCode, baseErr error) *DeviceIOControlError {
	return &DeviceIOControlError{InnerError: baseErr, ControlCode: controlCode}
}

func (err DeviceIOControlError) Error() string {
	return fmt.Sprintf("vmci: DeviceIoControl failed: %s (control code: 0x%x)", err.InnerError, err.ControlCode)
}

func (err DeviceIOControlError) Unwrap() error {
	return err.InnerError
}
