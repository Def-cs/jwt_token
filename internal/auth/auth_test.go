package auth

import (
	"github.com/stretchr/testify/assert"
	"jwt_auth.com/pkg/srorage/db/postgres"
	redisConn "jwt_auth.com/pkg/srorage/redis"
	"sync"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {

	time.Sleep(2 * time.Second)
	postgres.InitConn(5433, "localhost", "vladislav", "examplePass", "test_db")
	defer postgres.Connection.Close()
	postgres.SetupDB()

	redisConn.NewRedisConnection("localhost:6380", "", 0)

	deleteLiveTimeRefToken = &authDeleteRefToken{
		mu: &sync.Mutex{},
	}

	login := "test_user1"
	password := "test_password1"
	email := "test@mail.com"
	ip := "0.0.0.0"
	ipAnother := "0.0.0.1"

	err := CreateUser(login, password, email)
	assert.NoError(t, err)

	authToken, refToken, err := StartSession(ip, login, password)
	assert.NoError(t, err)

	res, err := AuthCheck(authToken)
	assert.NoError(t, err)
	assert.EqualValues(t, true, res)

	authToken, refToken, err = RefreshSessionTokens(ip, authToken, refToken)
	assert.NoError(t, err)

	_, _, err = RefreshSessionTokens(ipAnother, authToken, refToken)
	assert.NoError(t, err)
}
