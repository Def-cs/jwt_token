package mail

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMail(t *testing.T) {

	NewMailConnection()

	err := Connection.SendWarning("0.0.0.0", "vfrolov2004@gmail.com")

	require.NoError(t, err)
}
