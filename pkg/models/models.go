package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: ")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// create struct and define top level daya types which model will use and return
