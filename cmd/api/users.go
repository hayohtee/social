package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/hayohtee/social/internal/repository"
	"github.com/hayohtee/social/internal/validator"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
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

	if err := app.writeJSON(w, http.StatusOK, envelope{"user": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
		return
	}

	// TODO: revert back to auth userID from ctx
	var input struct {
		FollowerID int64 `json:"follower_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.FollowerID > 0, "follower_id", "must be a positive number")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.repository.Followers.Follow(r.Context(), user.ID, input.FollowerID); err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateKey):
			app.errorResponse(w, r, http.StatusConflict, "resource already exists")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "user followed successfully"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
		return
	}

	// TODO: revert back to auth userID from ctx
	var input struct {
		FollowerID int64 `json:"follower_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.FollowerID > 0, "follower_id", "must be a positive number")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.repository.Followers.UnFollow(r.Context(), user.ID, input.FollowerID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "user un-followed successfully"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
