package main

import (
	"net/http"

	uuid "github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"dental-clinic/internal/models"
)

// helpers
func convertToHash(object string) []byte {
	hashedValue, _ := bcrypt.GenerateFromPassword([]byte(object), bcrypt.MinCost)

	return hashedValue
}

func getUserFromCookie(res http.ResponseWriter, req *http.Request) models.User {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		// create cookie
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
	}
	// set cookie
	http.SetCookie(res, myCookie)

	// if user exists, get the user
	var myUser models.User
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = models.mapUsers[username]
	}

	return myUser
}

func alreadyLoggedIn(req *http.Request) bool {
	// if cookie doesn't exist, definitely not logged in
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}

	username := mapSessions[myCookie.Value] // is this a valid session?
	_, ok := models.mapUsers[username]      // does this person exist?

	return ok
}
