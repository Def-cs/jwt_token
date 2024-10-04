package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"jwt_auth.com/internal/mail"
	"jwt_auth.com/pkg/srorage/db/postgres"
	redisConn "jwt_auth.com/pkg/srorage/redis"
	"sync"
	"time"
)

var secretJWT = []byte("GIHRu4hg489ehHh44hHEFHROI484UYW")
var deleteLiveTimeRefToken *authDeleteRefToken

type authDeleteRefToken struct {
	mu *sync.Mutex
}

func init() {
	deleteLiveTimeRefToken = &authDeleteRefToken{
		mu: &sync.Mutex{},
	}
}

func CreateUser(login, password, email string) error {
	err := postgres.Connection.AddUser(login, password, email)
	return err
}

func StartSession(ip, login, password string) (string, string, error) {
	user, err := postgres.Connection.GetUserForAuth(login, password)
	if err != nil {
		return "", "", err
	}

	accToken, err := createAccToken(user.Id, ip)
	if err != nil {
		return "", "", err
	}

	err = redisConn.Connection.SetToken(accToken)
	if err != nil {
		return "", "", err
	}

	refToken, err := createRefToken(user.Id)
	if err != nil {
		return "", "", err
	}

	err = saveToken(refToken, user.Id)
	if err != nil {
		return "", "", err
	}

	go deleteLiveTimeRefToken.delete(user.Id)

	return accToken, refToken, nil
}

func createAccToken(uid uuid.UUID, ip string) (string, error) {
	claims := newClaims(uid, ip)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(secretJWT)
	return tokenStr, err
}

func newClaims(uid uuid.UUID, ip string) jwt.Claims {
	return jwt.MapClaims{
		"uid": uid,
		"ip":  ip,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
}

func createRefToken(uid uuid.UUID) (string, error) {
	bytes := make([]byte, 54)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)

	return token, nil
}

func saveToken(token string, uid uuid.UUID) error {
	hashToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = postgres.Connection.AddHashToken(uid, string(hashToken))
	return err
}

func updateToken(token string, uid uuid.UUID) error {
	hashToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = postgres.Connection.UpdateHashToken(uid, string(hashToken))
	return err
}

func AuthCheck(token string) (bool, error) {
	res, err := redisConn.Connection.GetToken(token)
	if err != nil {
		return res, err
	}

	_, err = verifyToken(token)
	return res, err
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretJWT, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("токен не валиден")
	}

	return token, nil
}

func RefreshSessionTokens(ip, authToken, refToken string) (string, string, error) {

	deleteLiveTimeRefToken.mu.Lock()

	token, _ := verifyToken(authToken)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("ошибка получения сведений токена")
	}

	uidStr, okUid := claims["uid"].(string)
	ipOld, okIp := claims["ip"].(string)

	if !okUid || !okIp {
		return "", "", errors.New("ошибка получения сведений uid или ip")
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return "", "", err
	}

	if ipOld != ip {
		user, err := postgres.Connection.GetUserByUid(uid)
		if err != nil {
			return "", "", err
		}

		mail.Connection.SendWarning(ip, user.Email)
	}

	hashToken, err := verityRefToken(uid, refToken)
	if err != nil {
		return "", "", err
	}
	if !hashToken {
		return "", "", errors.New("рефреш токен невалиден")
	}

	accToken, err := createAccToken(uid, ip)
	if err != nil {
		return "", "", err
	}

	err = redisConn.Connection.SetToken(accToken)
	if err != nil {
		return "", "", err
	}

	newRefToken, err := createRefToken(uid)
	if err != nil {
		return "", "", err
	}

	err = updateToken(newRefToken, uid)
	if err != nil {
		return "", "", err
	}

	deleteLiveTimeRefToken.mu.Unlock()

	go deleteLiveTimeRefToken.delete(uid)

	return accToken, newRefToken, nil
}

func verityRefToken(uid uuid.UUID, refTocken string) (bool, error) {

	token, err := postgres.Connection.GetHashToken(uid)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(token.Token), []byte(refTocken))
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (d *authDeleteRefToken) delete(uid uuid.UUID) {
	time.Sleep(1 * time.Hour)

	d.mu.Lock()
	defer d.mu.Unlock()

	err := postgres.Connection.DelHashToken(uid)

	fmt.Println(err.Error()) //думаю, что нет смысла его возвращать, можно, наверн, кидать в какую-нибудь систему логов
}
