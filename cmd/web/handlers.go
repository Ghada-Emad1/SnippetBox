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
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	app.notfound(w)
	// 	return
	// }
	//panic("OHH SOMETHING WENT WRONG")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.ServeError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
	// for _, snipet := range snippets {
	// 	fmt.Fprintf(w, "%v\n", snipet)
	// }
	// Files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }
	// ts, err := template.ParseFiles(Files...)
	// if err != nil {
	// 	app.ServeError(w, err)
	// 	return
	// }
	// data:=&templatesData{
	// 	Snippets:snippets,
	// }
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.ServeError(w, err)
	// 	return
	// }
	// w.Write([]byte("Hello From snippet Application"))
}
func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// if err != nil || id < 0 {
	// 	app.notfound(w)
	// 	return
	// }
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
	// Files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// ts, err := template.ParseFiles(Files...)
	// if err != nil {
	// 	app.ServeError(w, err)
	// 	return
	// }
	// data := &templatesData{
	// 	Snippet: snippet,
	// }
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.ServeError(w, err)
	// }

	// fmt.Fprintf(w, "Displaying a specific snippet with ID %v", snippet)
}
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This is Field can't be blank")
	form.CheckField(validator.NotBlank(form.Content), "content", "This Field can't be blank")
	form.CheckField(validator.MaxChar(form.Title, 100), "title", "This is field can't be more than 100 char")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This Field must equal 1,7 or 365")
	// if strings.TrimSpace(form.Title) == "" {
	// 	form.FieldsErrors["title"] = "This Field can't be blank"
	// } else if utf8.RuneCountInString("title") > 100 {
	// 	form.FieldsErrors["title"] = "This Field can't be more than 100 characters long"
	// }

	// if strings.TrimSpace(form.Content) == "" {
	// 	form.FieldsErrors["content"] = "This is Field can't be blank"
	// }

	// if form.Expires != 7 && form.Expires != 1 && form.Expires != 356 {
	// 	form.FieldsErrors["expires"] = "This Field must equal to 1,7 or 365"
	// }

	// if len(form.FieldsErrors) > 0 {
	// 	data := app.newTemplateData(r)
	// 	data.Form = form
	// 	app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
	// 	return
	// }

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

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
