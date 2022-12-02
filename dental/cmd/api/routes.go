package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	r := httprouter.New()

	// fileServer := http.FileServer(http.Dir("../../ui/static"))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	fileServer := http.FileServer(http.Dir("../../ui/static"))
	r.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	r.HandlerFunc(http.MethodGet, "/login", app.loginHandler)
	r.HandlerFunc(http.MethodPost, "/login", app.loginHandler)
	r.HandlerFunc(http.MethodGet, "/signup", app.signupHandler)
	r.HandlerFunc(http.MethodGet, "/logout", app.logoutHandler)
	r.HandlerFunc(http.MethodGet, "/appts", app.showAppointmentsHandler)
	r.HandlerFunc(http.MethodGet, "/users", app.showAllUsersHandler)

	r.HandlerFunc(http.MethodDelete, "/appts/delete/{id}", app.delAppointmentsHandler)

	// r.HandleFunc("/appts/delete/{id}", app.delAppointmentsHandler).Methods(http.MethodDelete)

	return r
}
