package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) ServeError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2,trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (app *Application) notfound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
