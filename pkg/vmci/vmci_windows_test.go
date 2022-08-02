package vmci

import (
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestVMSockAddrSize(t *testing.T) {
	sockAddr := sockAddrVM{}
	require.Equal(t, len(sockAddr.zero), 4)
	require.Equal(t, unsafe.Sizeof(sockAddr), uintptr(16))
}

func TestVersion(t *testing.T) {
	val, err := Version()
	require.NoError(t, err, "Version() returned an error")
	t.Logf("%d (Major: %d; Minor: %d; Epoch: %d)",
		val, val.Major(), val.Minor(), val.Epoch())
}
