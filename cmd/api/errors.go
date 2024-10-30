package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem when processing the request")
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error, message string) {
	log.Printf("bad request error: %s path: %s error: %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, message)
}
