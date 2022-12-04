package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	// fileServer := http.FileServer(http.Dir("../../ui/static"))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Middleware
	r.Use(ChangeMethod)
	r.Use(app.logRequest)

	// Load CSS file
	fileServer := http.FileServer(http.Dir("../../ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Handle("/appts/static/*", http.StripPrefix("/appts/static/", fileServer))
	r.Handle("/appts/delete/static/*", http.StripPrefix("/appts/delete/static/", fileServer))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.Get("/", app.homeHandler)

	// auth
	r.Get("/login", app.loginHandler)
	r.Post("/login", app.loginHandler)
	r.Get("/signup", app.signupHandler)
	r.Post("/signup", app.signupHandler)
	r.Get("/logout", app.logoutHandler)

	// appointments
	r.Get("/appts", app.showAppointmentsHandler)
	r.Get("/appts/book", app.bookAppointmentsHandler)
	r.Post("/appts/book", app.bookAppointmentsHandlerPut)

	// Admin-only pages
	r.Get("/users", app.showAllUsersHandler)

	r.Put("/appts/edit/{apptID}", app.editAppointmentHandler)
	r.Delete("/appts/delete/{apptID}", app.delAppointmentsHandler)

	// r.HandleFunc("/appts/delete/{id}", app.delAppointmentsHandler).Methods(http.MethodDelete)

	return r
}
