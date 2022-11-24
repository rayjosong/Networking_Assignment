package main

import (
	"html/template"
	"net/http"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	port = ":5221"
)

type User struct {
	Username  []byte // might need to be stored as hash value
	Password  []byte
	FirstName string
	LastName  string
	Role      string
}

var tpl *template.Template
var mapUsers = map[string]User{}
var mapSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	// creating admin user (this should not be the way to do it on a live system)
	mapUsers["admin"] = User{convertToHash("admin"), convertToHash("password"), "admin", "admin", "admin"}
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(port, nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	currentUser := getUser(res, req)

	// serve index.html
	tpl.ExecuteTemplate(res, "index.gohtml", currentUser)
}

func login(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		// redirect and return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		// check if user exist with username (draw from "db")
		myUser, ok := mapUsers[username]
		if !ok {
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}

		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
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

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}
