package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *Application) ServeError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
func (app *Application) notfound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
func (app *Application) render(w http.ResponseWriter, status int, page string, data *templatesData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServeError(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServeError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *Application) newTemplateData(r *http.Request) *templatesData {
	return &templatesData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(),"flash"),
	}
}

func (app *Application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invaildDecoderError *form.InvalidDecoderError
		if errors.As(err, &invaildDecoderError) {
			panic(err)
		}
		return err

	}
	return nil
}
