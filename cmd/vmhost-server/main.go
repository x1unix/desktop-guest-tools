package main

import (
	"fmt"
	"log"

	"github.com/x1unix/desktop-guest-tools/pkg/vmci"
)

func main() {
	val, err := vmci.GetAFValue()
	if err != nil {
		log.Fatalln(err)
	}

	if val == 0 {
		log.Fatalln("VMIC not available")
	}

	getVersion()
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
