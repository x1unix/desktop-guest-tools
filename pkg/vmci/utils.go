package vmci

func uuidFromString(str string) uuid {
	result := uuid{}
	if len(str) > len(result) {
		str = str[:len(result)]
	}
	copy(result[:], str)
	return result
}

func checkContextID(cid uint32) error {
	if cid == VMAddrCIDAny {
		return ErrContextUnavailable
	}

	return nil
}

// IsHypervisorContext reports whenever a passed context ID is hypervisor context.
func IsHypervisorContext(cid ContextID) bool {
	switch cid {
	case VMWareHypervisorCID, VMWareESXIHostCID, VMWarePlayerHostCID, VMWareInvalidCID:
		return true
	}
	return false
}
