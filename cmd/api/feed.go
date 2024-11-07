package main

import (
	"net/http"

	"github.com/hayohtee/social/internal/data"
	"github.com/hayohtee/social/internal/validator"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	var filters data.Filters

	qs := r.URL.Query()
	v := validator.New()

	filters.Page = readInt(qs, "page", 1, v)
	filters.PageSize = readInt(qs, "page_size", 20, v)
	filters.Sort = readString(qs, "sort", "ASC")
	filters.Tags = readCSV(qs, "tags", []string{})
	filters.Search = readString(qs, "search", "")

	if data.ValidateFilters(v, filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// TODO: Replace hardcoded user_id with actual one
	feeds, err := app.repository.Posts.GetUserFeeds(r.Context(), 100, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if feeds == nil {
		feeds = []data.Feed{}
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"items": feeds}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
