package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippeModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippeModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute. Split it over two lines for readability.
	// SQL request, used `` here because request has two strings.
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Used Exec() for making request, first its SQL request and second header of note body and lifetime of note.
	// This method returns sql.Result, which contains data what happened after request.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Used method LastInsertId() for get last ID of created note from table snippets.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Returned ID has tipe 'int64', so here we convert it into tipe 'int'.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippeModel) Get(id int) (*Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Used method QueryRow() for SQL request, return pointer to object sql.Row, which contains snippet's data.
	row := m.DB.QueryRow(stmt, id)

	// Initializes pointer for new structure Snippet.
	s := &Snippet{}

	// Uses row.Scan() to copy value from every sql.Row's field in structure Snippet's fields.
	// The arguments to row.Scan are *pointers* to he place you want,
	// to copy the data into and the number of arguments must be exactly the same
	// as the number of columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// Error check for fields data, if request has any error and error was detected,
		// than error returns from models.ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything is ok, returns object Snippet.
	return s, nil
}

// This will return the most recently created snippets.
func (m *SnippeModel) Latest() ([]*Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	// Used Query method for SQL request, we get sql.Rows with result of our request.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Defer rows.Close() for be sure that results from sql.Rows
	// will close correct before call method Latest() returns.
	// This delay operator must execute after check errors in method Query(),
	// if Query() return an error, you will get a panic trying to close a nil resultset.
	defer rows.Close()

	// Initialise empty slice for keeping objects models.Snippets.
	var snippets []*Snippet

	// Use rows.Next() for enumerating result,
	// this method prepares the first before rows.Scan().
	// If iteration over all the rows completes then the
	// resultset automatically closes itself and
	// frees-up the underlying database connection.
	for rows.Next() {
		// Create a pointer for sctuct Snippet.
		s := &Snippet{}

		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want, to copy the data into and the
		// number of arguments must be exactly the same, as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err()
	// to retrieve any error that was encountered during the iteration.
	// It's important to calls this, dont assume that a succesful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything is okay, then return the Snippets slice.
	return snippets, nil
}
