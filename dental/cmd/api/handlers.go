package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"
)

func (app *application) homeHandler(res http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData()

	currentUser, err := app.getUserFromCookie(res, req)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data.CurrentUser = currentUser

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/page/home.gohtml",
		"../../ui/partials/navbar.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	// serve index.html
	err = tpl.ExecuteTemplate(res, "base", data)
	if err != nil {
		log.Fatal(err)
	}
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

		data := app.newTemplateData()
		data.CurrentUser = myUser
		app.infoLog.Printf("Current User: %s\tRole: %s", myUser.Username, myUser.Role)

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
		"../../ui/partials/navbar.gohtml",
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

	data := app.newTemplateData()

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")

		if username != "" {
			// check if username exists/is taken
			// if exist/taken, return error and exit
			if u, _ := app.users.Get(username); u.Username == username {
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
			bPassword := convertToHash(password)

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
		"../../ui/partials/navbar.gohtml",
		"../../ui/page/signup.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	tpl.ExecuteTemplate(res, "base", data)
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

func (app *application) updateCredsHandler(res http.ResponseWriter, req *http.Request) {
	if !app.alreadyLoggedIn(req) {
		app.errorLog.Println("You are not logged in")
		http.Redirect(res, req, "/", http.StatusSeeOther) // redirect to login page
	}

	data := app.newTemplateData()
	currentUser, err := app.getUserFromCookie(res, req)
	if err != nil {
		app.errorLog.Println(err)
	}
	data.CurrentUser = currentUser

	if req.Method == http.MethodPost {
		newUsername := req.FormValue("username")
		newPassword := req.FormValue("password")

		if newUsername != "" {
			// check if username exists/is taken
			// if exist/taken, return error and exit
			if u, _ := app.users.Get(newUsername); u.Username == newUsername {
				http.Error(res, "Username already taken", http.StatusForbidden)
				return
			}

			// get current cookie and change the user
			myCookie, err := req.Cookie("myCookie")
			if err != nil {
				app.errorLog.Println(err)
			}

			mapSessions[myCookie.Value] = newUsername

			// encrypt the username & password then store user
			bNewPassword := convertToHash(newPassword)

			_, err = app.users.Update(currentUser.Uid, newUsername, bNewPassword)
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
		"../../ui/partials/navbar.gohtml",
		"../../ui/page/updateCreds.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	tpl.ExecuteTemplate(res, "base", data)
}

func (app *application) showAppointmentsHandler(res http.ResponseWriter, req *http.Request) {

	data := app.newTemplateData()
	currentUser, err := app.getUserFromCookie(res, req)
	if err != nil {
		app.errorLog.Println(err)
	}
	data.CurrentUser = currentUser
	fmt.Println("Current user role: ", data.CurrentUser.Role)

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/partials/navbar.gohtml",
		"../../ui/page/restrictAppts.gohtml",
	}

	fmt.Println(data.CurrentUser.Role)
	switch data.CurrentUser.Role {
	case "admin":
		appts, err := app.appointments.GetAll()
		if err != nil {
			app.errorLog.Println(err)
		}

		data.Appointments = appts

		files = []string{
			"../../ui/base.gohtml",
			"../../ui/partials/navbar.gohtml",
			"../../ui/page/adminAppts.gohtml",
		}

	case "patient":
		myAppts, err := app.appointments.Get(data.CurrentUser)
		if err != nil {
			app.errorLog.Println(err)
		}

		data.Appointments = myAppts

		files = []string{
			"../../ui/base.gohtml",
			"../../ui/partials/navbar.gohtml",
			"../../ui/page/userAppts.gohtml",
		}

	}

	funcMap := template.FuncMap{
		"formatDateTime": app.appointments.FormatDateTime,
	}

	// tpl, err := template.ParseFiles(files...)
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
	// if err != nil {
	// 	app.serverError(res, err)
	// }

	err = tpl.ExecuteTemplate(res, "base", data)
	if err != nil {
		app.errorLog.Fatalln(err)
	}
}

func (app *application) delAppointmentsHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		res.Header().Set("Allow", "DELETE")
		res.WriteHeader(405)
		res.Write([]byte("Method not allowed"))
		return
	}

	param := chi.URLParam(req, "apptID")
	apptID, err := strconv.Atoi(param)
	if err != nil {
		app.errorLog.Println(err)
		res.Write([]byte(err.Error()))
	}

	// Remove the appt
	err = app.appointments.Delete(apptID)
	if err != nil {
		app.errorLog.Println(err)
		res.Write([]byte(err.Error()))
	} else {
		files := []string{
			"../../ui/base.gohtml",
			"../../ui/partials/navbar.gohtml",
			"../../ui/page/notifications/notifyApptsDeleted.gohtml",
		}

		tpl := template.Must(template.New("").ParseFiles(files...))
		err = tpl.ExecuteTemplate(res, "base", nil)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}

	}
	// TODO: MAKE SURE THIS WORKS
	time.Sleep(time.Second * 1)
	http.Redirect(res, req, "/appts", http.StatusSeeOther)
}

func (app *application) bookAppointmentsHandler(res http.ResponseWriter, req *http.Request) {

	data := app.newTemplateData()
	currentUser, err := app.getUserFromCookie(res, req)
	if err != nil {
		app.errorLog.Println(err)
	}

	data.CurrentUser = currentUser

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/partials/navbar.gohtml",
		"../../ui/page/bookAppts.gohtml",
	}

	appts, err := app.appointments.GetAvailable()
	if err != nil {
		app.errorLog.Println(err)
	}

	data.Appointments = appts

	funcMap := template.FuncMap{
		"formatDateTime": app.appointments.FormatDateTime,
	}

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
	err = tpl.ExecuteTemplate(res, "base", data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (app *application) bookAppointmentsHandlerPut(res http.ResponseWriter, req *http.Request) {

	data := app.newTemplateData()
	currentUser, err := app.getUserFromCookie(res, req)
	if err != nil {
		app.errorLog.Println(err)
	}

	data.CurrentUser = currentUser

	req.ParseForm()
	selected := req.Form["selectedAppts"]

	fmt.Println(selected)

	for _, num := range selected {
		aID, err := strconv.Atoi(num)
		if err != nil {
			app.errorLog.Println(err)
		}

		updated, err := app.appointments.Update(aID, data.CurrentUser, false)
		if err != nil {
			app.errorLog.Println(err)
		}

		app.infoLog.Println("Updated payload: ", updated)
	}

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/partials/navbar.gohtml",
		"../../ui/page/notifications/notifyApptsUpdated.gohtml",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}

	err = tpl.ExecuteTemplate(res, "base", data)
	if err != nil {
		log.Fatal(err)
	}
}

// admin-only
func (app *application) showAllUsersHandler(res http.ResponseWriter, req *http.Request) {
	// obtain all users, pass to view

	allUsers, err := app.users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
	}

	data := app.newTemplateData()
	data.Users = allUsers

	// get current user to perform authorisation to view user data
	currentUser, err := app.getUserFromCookie(res, req)
	if err != nil {
		app.errorLog.Println(err)
	}

	data.CurrentUser = currentUser

	files := []string{
		"../../ui/base.gohtml",
		"../../ui/partials/navbar.gohtml",
		"../../ui/page/adminUsers.gohtml",
	}

	funcMap := template.FuncMap{
		"formatCreatedAt": app.users.FormatCreatedAt,
	}

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))

	err = tpl.ExecuteTemplate(res, "base", data)
	if err != nil {
		app.errorLog.Fatalln(err)
	}
}

func (a *application) editAppointmentHandler(res http.ResponseWriter, req *http.Request) {

}
