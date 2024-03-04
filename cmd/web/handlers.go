package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	Files:=[]string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts,err:=template.ParseFiles(Files...)
	if err!=nil{
		log.Println(err.Error())
		http.Error(w,"Internal Server Error",http.StatusInternalServerError)
		return
	}
	err=ts.ExecuteTemplate(w,"base",nil)
	if err!=nil{
		log.Println(err.Error())
		http.Error(w,"Internal Server Error",http.StatusInternalServerError)
		return
	}
	//w.Write([]byte("Hello From snippet Application"))
}
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Displaying a specific snippet with ID %d", id)
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "Post")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create A New Snippet"))
}
