package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"os"
	"time"

	"log"
	"net/http"

	"github.com/Ghada-Emad1/SnippetBox/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	_ "github.com/lib/pq"
)

type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	sessionManager *scs.SessionManager
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
	sessionManager:=scs.New()
	sessionManager.Store=postgresstore.New(db)
	sessionManager.Lifetime=12*time.Hour
	sessionManager.Cookie.Secure=true

	app := &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formdecoder,
		sessionManager: sessionManager,
	}
	tlsConfig:=&tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519,tls.CurveP256},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 10*time.Second,
	}

	infoLog.Printf("Starting App on Server %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")
	errorLog.Fatal(err)
}
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err=db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
