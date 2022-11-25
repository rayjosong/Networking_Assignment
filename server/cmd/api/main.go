package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	port = ":5221"
)

type User struct {
	Username  string // might need to be stored as hash value
	Password  []byte
	FirstName string
	LastName  string
	Role      string
}

var tpl *template.Template
var mapUsers = map[string]User{}
var mapSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("../../../frontend/web/templates/*"))

	// creating admin user (this should not be the way to do it on a live system)
	mapUsers["Admin"] = User{"Admin", convertToHash("Password"), "admin", "admin", "admin"} // do not do this. use json file outside of module instead
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/signup", signupHandler)
	r.HandleFunc("/logout", logoutHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  45 * time.Minute,
	}
	log.Fatal(srv.ListenAndServe())

}
