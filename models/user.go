package models

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `db:"id"`
	Email        string `db:"email" form:"email"`
	FirstName    string `db:"first_name" form:"first_name"`
	LastName     string `db:"last_name" form:"last_name"`
	Password     string `db:"password" form:"password"`
	PasswordHash string `db:"password_hash"`
}

type UserService struct {
	DB *sqlx.DB
}

func (s *UserService) CreateUser(data User) (*User, error) {
	data.PasswordHash, _ = s.hashPassword(data.Password)
	user := User{
		Email:        strings.ToLower(data.Email),
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
	return &user, nil
}

func (s *UserService) Authenticate(email, password string) (*User, error) {
	var user User
	err := s.DB.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	ok := s.checkPasswordHash(password, user.PasswordHash)
	if !ok {
		return nil, fmt.Errorf("incorrect password")
	}
	return &user, nil
}

func (*UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

func (*UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
