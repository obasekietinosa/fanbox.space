package models

import (
	"database/sql"
	"time"
)

type Letter struct {
	ID         int
	Email      string
	Subject    string
	Author     string
	Recipient  string
	Content    string
	Salutation string
	Created    time.Time
}

type LetterModel struct {
	DB *sql.DB
}

func (model *LetterModel) Insert(
	email string, subject string, author string, recipient string, content string, salutation string) (int, error) {
	stmt := `INSERT INTO letters (email, subject, author, recipient, content, salutation, created) VALUES (
			?, ?, ?, ?, ?, ?, UTC_TIMESTAMP()
		)`
	result, err := model.DB.Exec(stmt, email, subject, author, recipient, content, salutation)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (model *LetterModel) Get(id int) (*Letter, error) {
	return nil, nil
}

func (model *LetterModel) Latest(limit int) (*[]Letter, error) {
	return nil, nil
}
