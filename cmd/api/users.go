package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/hayohtee/social/internal/data"
	"github.com/hayohtee/social/internal/repository"
	"github.com/hayohtee/social/internal/validator"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
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

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, &user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	token, err := generateToken()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.repository.Users.CreateAndInvite(r.Context(), &user, token.Hash, app.config.mail.exp)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exist")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, repository.ErrDuplicateUsername):
			v.AddError("username", "a user with this username already exist")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var response struct {
		data.User
		Token string `json:"token"`
	}

	response.User = user
	response.Token = token.PlainText

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if token == "" {
		app.badRequestResponse(w, r, errors.New("no token provided"))
		return
	}

	if err := app.repository.Users.Activate(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"message": "account activated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
