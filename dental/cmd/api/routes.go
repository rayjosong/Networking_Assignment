package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()

	// fileServer := http.FileServer(http.Dir("../../ui/static"))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	fileServer := http.FileServer(http.Dir("../../ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	// r.PathPrefix("/css").Handler(twhandler.New(http.Dir("../../ui/static/css"), "/css", twembed.New()))

	r.HandleFunc("/", app.homeHandler).Methods("GET")
	r.HandleFunc("/login", app.loginHandler).Methods("GET")
	r.HandleFunc("/signup", app.signupHandler)
	r.HandleFunc("/logout", app.logoutHandler)
	r.HandleFunc("/appts", app.showAppointmentsHandler).Methods("GET")

	r.HandleFunc("")

	return r
}
