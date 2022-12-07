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

var counts int64

// contains services involved in this application
type application struct {
	users        *models.UserModel
	appointments *models.AppointmentsModel
	errorLog     *log.Logger
	infoLog      *log.Logger
}

// var tpl *template.Template

var mapSessions = map[string]string{}

func main() {
	dsn := flag.String("dsn", "web:pass@/dental?parseTime=true", "MySQL data")

	infoLog := log.New(os.Stdout, color.New(color.BgHiGreen).Sprintf(" INFO \t"), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, color.New(color.BgHiRed).Sprintf(" ERROR \t"), log.Ldate|log.Ltime|log.Lshortfile)

	// mySQL connection
	fmt.Println("Connecting to mysql...")
	conn := connectToDB(*dsn)
	// conn := connectToDB("user:my-secret@tcp(db:3306)/dental") // docker-compose option
	if conn == nil {
		log.Panic("Can't connect to mySQL!")
	}
	defer conn.Close()

	app := &application{
		users:    &models.UserModel{DB: conn},
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
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

func connectToDB(dsn string) *sql.DB {
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("mySQL not yet ready...")
			counts++
		} else {
			log.Println("Connected to mySQL!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
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
