package main

import (
	"dental-clinic/internal/models"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) homeHandler(res http.ResponseWriter, req *http.Request) {
	currentUser := app.getUserFromCookie(res, req)

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/page/home.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	// serve index.html
	tpl.ExecuteTemplate(res, "base", currentUser)
}

func (app *application) loginHandler(res http.ResponseWriter, req *http.Request) {
	if app.alreadyLoggedIn(req) {
		// redirect and return
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		// check if user exist with username (draw from "db")
		myUser, err := app.users.Get(username)
		if err != nil {
			app.notFound(res)
			return
		}

		// Matching of password entered
		err = bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			app.notAuthenticated(res)
			return
		}

		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}

		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = username
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/page/login.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}
	tpl.ExecuteTemplate(res, "base", nil)
}

func (app *application) signupHandler(res http.ResponseWriter, req *http.Request) {
	if app.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var myUser models.User
	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")

		if username != "" {
			// check if username exists/is taken
			// if exist/taken, return error and exit
			if u, _ := app.users.Get(username); u != nil {
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}

			// else, create session and user
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)
			mapSessions[myCookie.Value] = username

			// encrypt the username & password then store user
			// bUsername := convertToHash(username) // ignored potential error
			bPassword := convertToHash(password) // ignored potential error

			// myUser = models.User{username, bPassword, firstname, lastname, "patient"}
			_, err := app.users.Insert(username, bPassword, firstname, lastname, "patient")
			if err != nil {
				app.clientError(res, http.StatusBadRequest)
				return
			}

		}
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/page/signup.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	tpl.ExecuteTemplate(res, "base", myUser)
}

func (app *application) logoutHandler(res http.ResponseWriter, req *http.Request) {
	if !app.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther) // redirect to login page
	}

	myCookie, _ := req.Cookie("myCookie")

	// delete the session
	delete(mapSessions, myCookie.Value)

	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}
