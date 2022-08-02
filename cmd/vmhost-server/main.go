package main

import (
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/x1unix/desktop-guest-tools/pkg/vmci"
)

func main() {
	family := try(vmci.GetAFValue())
	if family < 0 {
		log.Fatalln("No VMCI found")
	}

	log.Printf("VMware VMCI address family: 0x%[1]x (%[1]d)\n", family)
	localCID := try(vmci.GetLocalCID())
	log.Println("Local CID:", localCID)
	dumpCID(localCID)
	getVersion()
	testConnect()
}

func testConnect() {
	log.Println("Trying to listen...")
	const port = 5000
	l, err := vmci.Listen(port)
	if err != nil {
		if vmci.IsUnsupportedProtocolError(err) {
			log.Println("Seems vSockets protocol is not registred in WinSock protocols catalog")
			log.Println("Please check: netsh winsock show catalog")
			log.Println("Or: HKLM\\SYSTEM\\CurrentControlSet\\Services\\WinSock2\\Parameters\\")
		}
		log.Fatalln(err)
	}
	log.Println("Listening on port", port)
	time.Sleep(5 * time.Second)
	log.Println("Closing socket")
	if err := l.Close(); err != nil {
		log.Fatal("failed to close socket:", err)
	}
}

func dumpCID(cid vmci.ContextID) {
	if !vmci.IsHypervisorContext(cid) {
		log.Println("Guest machine detected - CID=", cid)
		return
	}

	switch cid {
	case vmci.VMwareHypervisorCID:
		log.Println("Hypervisor: VMware Workstation or vSphere")
	case vmci.VMwareESXIHostCID:
		log.Println("Hypervisor: VMware ESXi")
	case vmci.VMwarePlayerHostCID:
		log.Println("Hypervisor: VMware Player")
	default:
		log.Println("Hypervisor: Not Found")
	}
}

func try[T any](val T, err error) T {
	if err != nil {
		strace := string(debug.Stack())
		log.Fatalln(err, "\n", strace)
	}
	runtime.Caller(1)
	return val
}

func getVersion() {
	ver, err := vmci.Version()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("vSockets Version: %d (Major: %d; Minor: %d; Epoch: %d)\n",
		ver.UInt32(), ver.Major(), ver.Minor(), ver.Epoch())
}
