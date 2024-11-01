package main

import (
	"fmt"
	"log"
	"net/http"
)

// The serverErrorResponse logs the detailed error message and uses the
// errorResponse helper method to send a 500 Internal Server Error status code
// and JSON response (containing a generic error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	message := http.StatusText(http.StatusInternalServerError)
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse uses the errorResponse helper method to send 404 Not Found status code
// and JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse uses the errorResponse helper method to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse uses the errorResponse helper method to send a 400 Bad Request status code
// and JSON response containing the error message to the client.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse writes 422 Unprocessable Entity and the contents of the error as JSON response.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}