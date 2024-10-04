package api

import (
	"encoding/json"
	"fmt"
	"jwt_auth.com/internal/auth"
	"jwt_auth.com/pkg/dto"
	"net/http"
)

func createUser(r *http.Request) ([]byte, int) {
	var err error
	var data dto.CreateUserRequest

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res := fmt.Sprintf(`{"err": "%s"}`, err.Error())
		return []byte(res), http.StatusNotAcceptable
	}

	err = auth.CreateUser(data.Login, data.Password, data.Email)
	if err != nil {
		res := fmt.Sprintf(`{"err": "%s"}`, err.Error())
		return []byte(res), http.StatusBadRequest
	}
	res := `{"res": "пользователь успешно создан"}`
	return []byte(res), http.StatusCreated
}

func login(r *http.Request) ([]byte, int) {
	var data dto.LoginRequest
	var token, refreshToken string
	userIp := r.RemoteAddr

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res := fmt.Sprintf(`{"err": "%s"}`, err.Error())
		return []byte(res), http.StatusNotAcceptable
	}

	fmt.Println("SESSION STARTED")
	token, refreshToken, err = auth.StartSession(userIp, data.Login, data.Password)
	if err != nil {
		res := fmt.Sprintf(`{"err": "%s"}`, err.Error())
		return []byte(res), http.StatusNotFound
	}

	tokens := dto.TokensObj{
		Token:        token,
		RefreshToken: refreshToken,
	}

	res, _ := json.Marshal(tokens)
	return res, http.StatusOK
}

func refresh(r *http.Request) ([]byte, int) {
	var data dto.RefreshTokenRequest
	authToken := r.Header.Get("Authorization")
	userIp := r.RemoteAddr

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		res := fmt.Sprintf(`{"err": "%s"}`, err.Error())
		return []byte(res), http.StatusNotAcceptable
	}

	accToken, refreshToken, err := auth.RefreshSessionTokens(userIp, authToken, data.RefToken)

	tokens := dto.TokensObj{
		Token:        accToken,
		RefreshToken: refreshToken,
	}

	res, _ := json.Marshal(tokens)
	return res, http.StatusOK
}
