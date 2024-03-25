package models

import "database/sql"

type User struct{}

type UserService struct {
	DB *sql.DB
}
