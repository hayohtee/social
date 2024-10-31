package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
