package mysql

import (
	"database/sql"
	"errors"
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
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// SQL request for returning data for one snippet
	row := m.DB.QueryRow(stmt, id)
	// used method QueryRow() for SQL request, return pointer to object sql.Row, which contains snippet's data

	s := &models.Snippet{}
	// initializes pointer for new structure Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires) // uses row.Scan() for copy value from every sql.Row's field in structure Snippet's fields
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // error check for fields data, if request has any error and error was detected, than error returns from models.ErrNoRecord
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil // if everything is ok, returns object Snippet
}

// method for returning note's data by note ID

func (m *SnippeModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

// method returns the most frequently used notes (10)
