package main

import (
	"fmt"
	"net/http"
)

// Asserts the HTTP method to chosen method before sending it to the request handler
func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("I AM IN THE MIDDLEWARE")
		fmt.Println("HTTP METHOD:", req.Method)
		if req.Method == http.MethodPost {
			switch method := req.PostFormValue("_method"); method {

			case http.MethodPut:
				fmt.Println(method)
				fallthrough
			case http.MethodPatch:
				fmt.Println(method)
				fallthrough
			case http.MethodDelete:
				fmt.Println(method)
				req.Method = method
			}
		}
		next.ServeHTTP(res, req)
	})
}

// Log all HTTP requests
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
