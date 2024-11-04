package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/hayohtee/social/internal/repository"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r, "userID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.repository.Users.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var response struct {
		ID        int64     `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}

	response.ID = user.ID
	response.Username = user.Username
	response.Email = user.Email
	response.CreatedAt = user.CreatedAt

	if err = app.writeJSON(w, http.StatusOK, envelope{"user": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	
}