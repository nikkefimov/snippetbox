package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize the router
	router := httprouter.New()

	// Create a handler function which wramps our notFound() helper
	// and then assign it as the custom handler for 404 Not Found responses
	// Also set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Update the patter for the route for the static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.showSnippet)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.createSnippet)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet", app.showSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)

	// Update the pattern for the route for the static files
	//fs := http.FileServer(nfs.NeuteredFileSystem{Fs: http.Dir(",/ui/static")})
	//mux.Handle("/static", http.StripPrefix("/static", fs))

	// Create the middleware chain as normal
	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the router with the middleware and return it as normal
	return standart.Then(router)
}
