package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	// Initialize a new time.Time object and pass it to the humanDate function.
	tm := time.Date(2024, 9, 20, 12, 00, 0, 0, time.UTC)
	hd := humanDate(tm)

	// Check that the output from the humanDate function is in the format we expect. If it isn't
	// that we expect, use the t.Errorf() function to
	// indicate that the test has failed and log the expected and actual values.
	if hd != "20 Sep 2024 at 12:00" {
		t.Errorf("got %q; want %q", hd, "20 Sep 2024 at 12:00")
	}

}
