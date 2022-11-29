package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"dental-clinic/internal/models"
)

const (
	port = ":5221"
)

type application struct {
	users *models.UserModel
}

// var tpl *template.Template

var mapSessions = map[string]string{}

func init() {
	// tpl = template.Must(template.ParseGlob("../../../frontend/web/templates/*"))

	// creating admin user (this should not be the way to do it on a live system)
	// mapUsers["Admin"] = User{"Admin", convertToHash("Password"), "admin", "admin", "admin"} // do not do this. use json file outside of module instead
}

func main() {
	dsn := flag.String("dsn", "web:pass@/dental?parseTime=true", "MySQL data")

	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("mysql running")
	}

	defer db.Close()

	app := &application{
		users: models.UserModel{DB: db},
	}

	r := mux.NewRouter()

	r.HandleFunc("/", app.homeHandler)
	// r.HandleFunc("/login", app.loginHandler)
	// r.HandleFunc("/signup", app.signupHandler)
	// r.HandleFunc("/logout", app.logoutHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  45 * time.Minute,
	}
	fmt.Println("Server running on port ", port)
	log.Fatal(srv.ListenAndServe())

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
