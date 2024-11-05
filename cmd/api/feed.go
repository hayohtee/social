package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Replace hardcoded user_id with actual one
	feeds, err := app.repository.Posts.GetUserFeeds(r.Context(), 100)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"items": feeds}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
