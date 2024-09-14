package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// The patter for the route for the static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Unprotected application routes using the 'dynamic' middleware chain.
	// Use the nosurf middleware on all 'dynamic' routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// Unprotected "dynamic" routes.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.showSnippet))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-onlu) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	// The 'protected' middleware chain appends to the 'dynaic' chain
	// the noSurf middleware will also be used on the tree routes below too.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create the middleware chain as normal.
	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the router with the middleware and return it as normal.
	return standart.Then(router)
}
