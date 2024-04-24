package models

import (
	"github.com/jmoiron/sqlx"
)

type Contact struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type ContactService struct {
	DB *sqlx.DB
}

func (s *ContactService) CreateContact(email string, userId int) (*Contact, error) {
	var contact Contact
	targetUser, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	row, err := s.DB.Query("INSERT INTO contacts (user_id, contact_id) VALUES ($1, $2) returning id;", userId, targetUser.ID)
	if err != nil {
		return nil, err
	}
	if row.Next() {
		row.Scan(&contact.ID)
	}
	contact.Email = email
	contact.FirstName = targetUser.FirstName
	contact.LastName = targetUser.LastName

	return &contact, nil
}

func (s *ContactService) GetContacts(userId int) ([]Contact, error) {
	var contacts []Contact
	err := s.DB.Select(&contacts, `
		SELECT c.id, u.email, u.first_name, u.last_name
		FROM contacts c
		INNER JOIN users u ON c.contact_id = u.id
		WHERE c.user_id = $1;
	`, userId)
	return contacts, err
}

func (s *ContactService) GetContactByIdAndValidate(contactId, userId int) (*Contact, error) {
	var contact Contact
	err := s.DB.Get(&contact, `
		SELECT c.id, u.email, u.first_name, u.last_name
			FROM contacts c
		INNER JOIN users u ON c.contact_id = u.id
			WHERE c.id = $1 AND c.user_id = $2;
	`, contactId, userId)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *ContactService) GetContactUserId(contactId int) (int, error) {
	var userId int
	row := s.DB.QueryRow(`
		SELECT u.id
			FROM contacts c
		INNER JOIN users u ON c.contact_id = u.id
			WHERE c.id = $1;
	`, contactId)
	err := row.Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *ContactService) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.DB.Get(&user, "SELECT id, email, first_name, last_name FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
