package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ghada-Emad1/SnippetBox/internal/models"
	"github.com/julienschmidt/httprouter"
)

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

	data:=app.newTemplateData(r)
	data.Snippets=snippets
	app.render(w, http.StatusOK, "home.tmpl",data)
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
	params:=httprouter.ParamsFromContext(r.Context())
	id,err:=strconv.Atoi(params.ByName("id"))
	if err!=nil{
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
	data:=app.newTemplateData(r)
	data.Snippet=snippet
	app.render(w,http.StatusOK,"view.tmpl",data)
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
func(app *Application)snippetCreate(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Display the form for creating new snippet"))
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", "Post")
	// 	app.ClientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"

	// expires,err := strconv.Atoi("7") // current local time
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServeError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
