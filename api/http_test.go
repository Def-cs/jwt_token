package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"jwt_auth.com/pkg/dto"
	"jwt_auth.com/pkg/srorage/db/postgres"
	redisConn "jwt_auth.com/pkg/srorage/redis"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	mux := RegisterMux()

	login := "test_user1"
	password := "test_password1"
	email := "test@mail.com"

	ts := httptest.NewServer(mux)
	defer ts.Close()

	time.Sleep(2 * time.Second)
	postgres.InitConn(5433, "localhost", "vladislav", "examplePass", "test_db")
	defer postgres.Connection.Close()
	postgres.SetupDB()

	redisConn.NewRedisConnection("localhost:6380", "", 0)

	client := &http.Client{}

	createUserRequest := dto.CreateUserRequest{
		Login:    login,
		Password: password,
		Email:    email,
	}
	createUserRequestBody, _ := json.Marshal(createUserRequest)

	req, err := http.NewRequest("GET", ts.URL+"/register", bytes.NewBuffer(createUserRequestBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	loginRequest := dto.LoginRequest{
		Login:    login,
		Password: password,
	}
	loginRequestBody, _ := json.Marshal(loginRequest)

	req, err = http.NewRequest("GET", ts.URL+"/login", bytes.NewBuffer(loginRequestBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err = client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var loginResponseVal *dto.TokensObj

	err = json.NewDecoder(res.Body).Decode(&loginResponseVal)
	assert.NoError(t, err)

	refreshRequest := dto.RefreshTokenRequest{
		RefToken: loginResponseVal.RefreshToken,
	}
	refreshRequestBody, _ := json.Marshal(refreshRequest)

	req, err = http.NewRequest("GET", ts.URL+"/refresh", bytes.NewBuffer(refreshRequestBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", loginResponseVal.Token)

	res, err = client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

}
