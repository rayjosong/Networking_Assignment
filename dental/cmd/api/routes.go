package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", app.homeHandler)
	r.HandleFunc("/login", app.loginHandler)
	r.HandleFunc("/signup", app.signupHandler)
	r.HandleFunc("/logout", app.logoutHandler)

	return r
}
