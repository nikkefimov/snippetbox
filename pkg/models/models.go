package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: ")

// create struct and define top level data types which model will use and return
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
