package models

import (
	"database/sql"
	"errors"
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
	stmt := "SELECT id, email, subject, author, recipient, content, salutation, created from letters where id = ?"
	row := model.DB.QueryRow(stmt, id)

	letter := &Letter{}

	err := row.Scan(&letter.ID, &letter.Email, &letter.Subject, &letter.Author, &letter.Recipient, &letter.Content, &letter.Salutation, &letter.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return letter, nil
}

func (model *LetterModel) Latest(limit int) ([]*Letter, error) {
	stmt := "SELECT id, email, subject, author, recipient, content, salutation, created from letters ORDER BY id DESC LIMIT ?"
	rows, err := model.DB.Query(stmt, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	letters := []*Letter{}

	for rows.Next() {
		letter := &Letter{}
		err = rows.Scan(&letter.ID, &letter.Email, &letter.Subject, &letter.Author, &letter.Recipient, &letter.Content, &letter.Salutation, &letter.Created)
		if err != nil {
			return nil, err
		}
		letters = append(letters, letter)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return letters, nil
}
