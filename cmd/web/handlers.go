package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ghada-Emad1/SnippetBox/internal/models"
	"github.com/Ghada-Emad1/SnippetBox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int     `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.ServeError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)

}
func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notfound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notfound(w)
		} else {
			app.ServeError(w, err)
		}
		return
	}
	
    data := app.newTemplateData(r)
	data.Snippet = snippet


	app.render(w, http.StatusOK, "view.tmpl", data)

}
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    var form snippetCreateForm
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.ClientError(w, http.StatusBadRequest)
        return
    }
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxChar(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }
    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.ServeError(w, err)
        return
    }
	app.sessionManager.Put(r.Context(),"flash","snippet successfully created!")
    
	
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
 }
