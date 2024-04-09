package main

import (
	"database/sql"
	"flag"
	"os"

	"log"
	"net/http"

	"github.com/Ghada-Emad1/SnippetBox/internal/models"

	_ "github.com/lib/pq"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	snippets *models.SnippetModel
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
	app := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.Routes(),
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
