package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           int    `db:"id"`
	Email        string `db:"email"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	PasswordHash string `db:"password_hash"`
}

type CreateUserDTO struct {
	Email        string `db:"email" form:"email"`
	FirstName    string `db:"first_name" form:"first_name"`
	LastName     string `db:"last_name" form:"last_name"`
	Password     string `db:"password" form:"password"`
	PasswordHash string `db:"password_hash"`
}

type UserService struct {
	DB *sqlx.DB
}

func (s *UserService) CreateUser(data CreateUserDTO) (*User, error) {
	user := User{
		Email:        data.Email,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		PasswordHash: data.PasswordHash,
	}
	rows, err := s.DB.NamedQuery("INSERT INTO users (email, first_name, last_name, password_hash) VALUES (:email, :first_name, :last_name, :password_hash) returning id;", data)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if rows.Next() {
		rows.Scan(&user.ID)
	}
	fmt.Println(user)
	return nil, nil
}
