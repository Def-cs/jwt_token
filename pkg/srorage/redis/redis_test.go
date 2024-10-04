package redisConn

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRedis(t *testing.T) {

	NewRedisConnection("localhost:6380", "", 0)

	token := "test_token"

	err := Connection.SetToken(token)
	require.NoError(t, err)

	res, err := Connection.GetToken(token)

	require.EqualValues(t, true, res)
	require.NoError(t, err)

}
