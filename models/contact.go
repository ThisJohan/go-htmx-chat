package models

import "github.com/jmoiron/sqlx"

type Contact struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type ContactService struct {
	DB *sqlx.DB
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
