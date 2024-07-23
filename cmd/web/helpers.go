package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// helper serverError writes error msg in errorLog and send to user answer 500 "Internal server error"

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// helper clientError sends exact status code and exact discription to user, in next steps, later it will look like 400 "Bad request"
// it happens when we have problem with user's request

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// helper notFound its a something like convenient shell around clientError, which sends to user answer "404 error"
