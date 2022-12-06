package main

import (
	"net/http"
)

// Asserts the HTTP method to chosen method before sending it to the request handler
func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			switch method := req.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				req.Method = method
			}
		}
		next.ServeHTTP(res, req)
	})
}

// Log all HTTP requests
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())
		next.ServeHTTP(res, req)
	})
}
