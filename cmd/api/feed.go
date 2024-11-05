package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.repository.Posts.GetUserFeed(r.Context(), 32)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"items": feeds}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}