# VMWare VMCI/vSockets bindings for Go

This package contains VMWare [vSockets](https://developer.vmware.com/docs/5521/vmware-vsockets-programming-guide/doc/vsockPreface.html) and [VMCI](https://www.vmware.com/pdf/ws7_esx4_vmci_sockets.pdf) bindings for Go.

## Usage

Library implements host-side logic. Please use `vsock` protocol on Linux guest side.

See [example_test.go](example_test.go) for usage example.

## What is vSockets?

The vSockets API facilitates fast and efficient communication between VMWare guest virtual machines and their host.
VMware vSockets are built on the VMCI device.

vSockets are supported on:

* VMWare Workstation
* VMWare Player
* VMWare ESXi

See [VMWare vSockets Programming Guide](https://developer.vmware.com/docs/5521/vmware-vsockets-programming-guide/doc/vsockPreface.html) for more information about vSockets.
