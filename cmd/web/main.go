package main

import (
	"database/sql"
	"flag"
	"html/template"
	"os"

	"log"
	"net/http"

	"github.com/Ghada-Emad1/SnippetBox/internal/models"
	"github.com/go-playground/form/v4"

	_ "github.com/lib/pq"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "user=ghada password=pass dbname=snippetbox sslmode=disable", "postgres data source")
	flag.Parse()

	infoLog := log.New(os.Stdout, "Info \t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stderr, "Error \t", log.Ltime|log.Ldate|log.Lshortfile)

	db, err := OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formdecoder:=form.NewDecoder()
	app := &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formdecoder,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting App on Server %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
