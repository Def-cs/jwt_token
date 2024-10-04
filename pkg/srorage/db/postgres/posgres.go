package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"jwt_auth.com/pkg/srorage/db"
)

var Connection PostgresConn

type PostgresConn struct {
	conn *sql.DB
}

func InitConn(port int, host, user, password, dbname string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	databBase, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = databBase.Ping()
	if err != nil {
		panic(err)
	}

	Connection.conn = databBase
}

func (p *PostgresConn) Close() {
	p.conn.Close()
}

func (p *PostgresConn) AddUser(login, password, email string) error {
	_, err := p.conn.Exec(`INSERT INTO users (login, password, email) VALUES ($1, $2, $3);`, login, password, email)
	return err
}

func (p *PostgresConn) GetUserForAuth(login, password string) (db.User, error) {
	var user db.User

	row := p.conn.QueryRow(`SELECT * FROM users WHERE users.login = $1 AND users.password = $2`, login, password)

	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Email)

	return user, err
}

func (p *PostgresConn) GetUserByUid(uid uuid.UUID) (db.User, error) {
	var user db.User

	row := p.conn.QueryRow(`SELECT * FROM users WHERE users.id = $1`, uid)

	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Email)

	return user, err
}

// только для тестовой бд
func (p *PostgresConn) deleteUserByUid(uid uuid.UUID) error {
	_, err := p.conn.Exec(`DELETE FROM users WHERE id = $1`, uid)
	return err
}

func (p *PostgresConn) AddHashToken(uid uuid.UUID, token string) error {
	_, err := p.conn.Exec(`INSERT INTO tokens (user_id, token) VALUES ((SELECT id FROM users WHERE id = $1), $2);`, uid, token)
	return err
}

func (p *PostgresConn) GetHashToken(uid uuid.UUID) (db.Token, error) {
	var token db.Token

	row := p.conn.QueryRow(`SELECT * FROM tokens WHERE user_id = $1`, uid)
	err := row.Scan(&token.UserId, &token.Token)

	return token, err
}

func (p *PostgresConn) DelHashToken(uid uuid.UUID) error {
	_, err := p.conn.Exec(`DELETE FROM tokens WHERE user_id = $1;`, uid)
	return err
}

func (p *PostgresConn) UpdateHashToken(uid uuid.UUID, token string) error {
	_, err := p.conn.Exec(`UPDATE tokens SET token = $1 WHERE user_id = $2;`, token, uid.String())
	return err
}
