package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	//return app.recoverPanic(app.logRequeset(secureHeaders(mux)))
	mychanin := alice.New(app.recoverPanic, app.logRequeset, secureHeaders)
	return mychanin.Then(mux)
}
