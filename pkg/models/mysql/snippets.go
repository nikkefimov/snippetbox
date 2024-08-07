package mysql

import (
	"database/sql"
	"snippetbox/pkg/models"
)

// here import new created models package which was created in file models.go before and write path to

type SnippeModel struct {
	DB *sql.DB
}

// define the type which wraps the connection pool

func (m *SnippeModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)           
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))` // SQL request, used `` here because request has two strings
	result, err := m.DB.Exec(stmt, title, content, expires) // used Exec() for making request, first its SQL request and second header of note, body of note and lifetime of note, this method returns sql.Result, which contains data what happened after request
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId() // used method LastInsertId() for get last ID of created note from table snippets
	if err != nil {
		return 0, err
	}
	return int(id), nil // returned ID has tipe 'int64', so here we convert it into tipe 'int'
}

// method for creating a new note in database

func (m *SnippeModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// method for returning note's data by note ID

func (m *SnippeModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

// method returns the most frequently used notes (10)
