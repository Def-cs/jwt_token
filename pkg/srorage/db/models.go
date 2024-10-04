package db

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Login    string
	Password string
	Email    string
}

type Token struct {
	UserId string
	Token  string
}
