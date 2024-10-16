package main

import (
	"github.com/justinas/alice"
	"net/http"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the static dir
	// Note that the path given to the Dir function is relative
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.handle function to register the file server as the handler
	// for all urls that start with static and strip the prefix
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Swap the route declarations to use the application struct's methods as the
	// handler functions.
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)

}
