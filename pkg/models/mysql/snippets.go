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
	return 0, nil
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
