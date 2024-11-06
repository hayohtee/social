package main

import (
	"net/http"

	"github.com/hayohtee/social/internal/validator"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	var f filters

	qs := r.URL.Query()
	v := validator.New()

	f.Page = readInt(qs, "page", 1, v)
	f.PageSize = readInt(qs, "page_size", 20, v)
	f.Sort = readString(qs, "sort", "ASC")

	if ValidateFilters(v, f); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

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
