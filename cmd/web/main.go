package main

import (
	"flag"
	"os"

	//"fmt"
	"log"
	"net/http"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	infoLog := log.New(os.Stdout, "Info \t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stderr, "Error \t", log.Ltime|log.Ldate|log.Lshortfile)
	flag.Parse()

	
	app := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.Routes(),
		ErrorLog: errorLog,
	}

	
	//fmt.Println("Start listen on port :4000")
	infoLog.Printf("Starting App on Server %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
