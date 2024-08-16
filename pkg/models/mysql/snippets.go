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
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt) // used Query method for SQL request, we get sql.Rows with result of our request
	if err != nil {
		return nil, err
	}

	defer rows.Close() // put on hold rows.Close() for be sure that results from sql.Rows will close correct before call method Latest(). This delay operator must execute after check errors in method Query()
	// Another case if Query() returns error it will lead to panic, because Query() will try to close results with value nil.

	var snippets []*models.Snippet // initialise empty slice for keeping objects models.Snippets

	for rows.Next() { // use rows.Next() for enumerating result, this method will call first before rows.Scan()
		s := &models.Snippet{} // create pointer for sctuct Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires) // use rows.Scan() for copying fields values into the sctruct
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s) // adding struct in slice
	}

	if err = rows.Err(); err != nil { // when the rows.Next finishes we call the rows.Err() method to find out, if we have encountered any error during the work
		return nil, err
	}

	return snippets, nil // if everything is okay, return slice with data
}

// method returns the most frequently used notes (10)
