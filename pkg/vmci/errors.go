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
)

// DeviceIOControlError is DeviceIOControl error.
//
// Error occurs when invalid control code was supplied or ioctl operation failed.
//
// See DeviceIOControl.
type DeviceIOControlError struct {
	ControlCode ControlCode
	InnerError  error
}

func newDeviceIOControlError(controlCode ControlCode, baseErr error) *DeviceIOControlError {
	return &DeviceIOControlError{InnerError: baseErr, ControlCode: controlCode}
}

func (err DeviceIOControlError) Error() string {
	return fmt.Sprintf("vmci: DeviceIoControl failed: %s (control code: 0x%x)", err.InnerError, err.ControlCode)
}

func (err DeviceIOControlError) Unwrap() error {
	return err.InnerError
}
