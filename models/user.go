package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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
	data.PasswordHash, _ = hashPassword(data.Password)
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
	return nil, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
