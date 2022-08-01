package vmci

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestVersion(t *testing.T) {
	val, err := Version()
	require.NoError(t, err, "Version() returned an error")
	t.Logf("%d (Major: %d; Minor: %d; Epoch: %d)",
		val, val.Major(), val.Minor(), val.Epoch())
}
