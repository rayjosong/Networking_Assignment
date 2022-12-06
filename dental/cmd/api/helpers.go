package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	uuid "github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"dental-clinic/internal/models"
)

// helpers
func convertToHash(object string) []byte {
	hashedValue, _ := bcrypt.GenerateFromPassword([]byte(object), bcrypt.MinCost)

	return hashedValue
}

func (app *application) getUserFromCookie(res http.ResponseWriter, req *http.Request) (models.User, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Trapped panic: %s (%T)\n", r, r)
		}
	}()

	myCookie, err := req.Cookie("myCookie")
	fmt.Println(myCookie.Value)

	if err != nil {
		app.errorLog.Println("Need to create new cookie lah.", err)
		// create cookie
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
	}

	var myUser models.User
	// if user exists, get the user
	if username, ok := mapSessions[myCookie.Value]; ok {
		// retrieve user details from database
		myUser, err = app.users.Get(username)
		if err != nil {
			return myUser, err
		}
		return myUser, nil
	}

	return models.User{}, err
}

func (app *application) alreadyLoggedIn(req *http.Request) bool {
	// if cookie doesn't exist, definitely not logged in
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}

	_, ok := mapSessions[myCookie.Value] // is this a valid session?

	return ok
}

// Writes error message and ends a generic 500 Internal Server Error response to the user
func (app *application) serverError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)

	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends specific status code and corresponding description to the user (e.g. 400 "Bad Request")
func (app *application) clientError(res http.ResponseWriter, status int) {
	http.Error(res, http.StatusText(status), status)
}

func (app *application) notFound(res http.ResponseWriter) {
	app.clientError(res, http.StatusNotFound)
}

func (app *application) notAuthenticated(res http.ResponseWriter) {
	app.clientError(res, http.StatusUnauthorized)
}
