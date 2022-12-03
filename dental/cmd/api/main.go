package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"

	"dental-clinic/internal/models"
)

const (
	port = ":5221"
)

// contains services involved in this application
type application struct {
	users        *models.UserModel
	appointments *models.AppointmentsModel
	errorLog     *log.Logger
	infoLog      *log.Logger
}

// var tpl *template.Template

var mapSessions = map[string]string{}

// TODO: REMOVE THIS IF NOT NEEDED
func init() {
	// tpl = template.Must(template.ParseGlob("../../../frontend/web/templates/*"))

	// creating admin user (this should not be the way to do it on a live system)
	// mapUsers["Admin"] = User{"Admin", convertToHash("Password"), "admin", "admin", "admin"} // do not do this. use json file outside of module instead
}

func main() {
	dsn := flag.String("dsn", "web:pass@/dental?parseTime=true", "MySQL data")

	infoLog := log.New(os.Stdout, color.New(color.BgHiGreen).Sprintf(" INFO \t"), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, color.New(color.BgHiRed).Sprintf(" ERROR \t"), log.Ldate|log.Ltime|log.Lshortfile)

	// mySQL connection
	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("mysql running")
	}
	defer db.Close()

	app := &application{
		users:    &models.UserModel{DB: db},
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Handler:      app.routes(),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  45 * time.Minute,
	}

	infoLog.Printf("Starting server on %s", port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
