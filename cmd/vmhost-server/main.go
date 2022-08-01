package main

import (
	"fmt"
	"log"

	"github.com/x1unix/desktop-guest-tools/pkg/vmci"
)

func main() {
	family := try(vmci.GetAFValue())
	if family == 0 {
		log.Fatalln("No VMCI found")
	}

	fmt.Println("VMware VMCI address family:", family)
	localCID := try(vmci.GetLocalCID())
	fmt.Println("local ContextID:", localCID)
	dumpCID(localCID)
	getVersion()
}

func dumpCID(cid vmci.ContextID) {
	if !vmci.IsHypervisorContext(cid) {
		fmt.Println("Guest machine detected - CID=", cid)
		return
	}

	fmt.Print("Hypervisor: ")
	switch cid {
	case vmci.VMwareHypervisorCID:
		fmt.Println("VMware Workstation")
	case vmci.VMwareESXIHostCID:
		fmt.Println("VMware ESXi")
	case vmci.VMwarePlayerHostCID:
		fmt.Println("VMware Player")
	default:
		fmt.Println("Not Found")
	}
}

func try[T any](val T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return val
}

func getVersion() {
	ver, err := vmci.Version()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Version:", ver.UInt32())
	fmt.Println("\tMajor:", ver.Major())
	fmt.Println("\tMinor:", ver.Minor())
	fmt.Println("\tEpoch:", ver.Epoch())
}
