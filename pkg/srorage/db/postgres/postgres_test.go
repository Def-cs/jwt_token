package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostgres(t *testing.T) {

	login := "test_user1"
	password := "test_password1"
	email := "test@,ail.com"
	token := "test_token"
	tokenNew := "test_token_new"

	time.Sleep(2 * time.Second)
	InitConn(5433, "localhost", "vladislav", "examplePass", "test_db")
	defer Connection.Close()
	SetupDB()

	err := Connection.AddUser(login, password, email)
	assert.NoError(t, err)

	user, err := Connection.GetUserForAuth(login, password)
	assert.NoError(t, err)

	err = Connection.AddHashToken(user.Id, token)
	assert.NoError(t, err)

	err = Connection.UpdateHashToken(user.Id, tokenNew)
	assert.NoError(t, err)

	tokenGot, err := Connection.GetHashToken(user.Id)
	assert.NoError(t, err)
	assert.EqualValues(t, tokenNew, tokenGot.Token)

	err = Connection.DelHashToken(user.Id)
	assert.NoError(t, err)

	userByID, err := Connection.GetUserByUid(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, user, userByID)

	err = Connection.DelHashToken(user.Id)
	assert.NoError(t, err)

	err = Connection.deleteUserByUid(user.Id)
	assert.NoError(t, err)

}
