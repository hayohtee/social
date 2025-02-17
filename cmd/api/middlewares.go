package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/hayohtee/social/internal/repository"
)

type contextKey string

const (
	postKey = contextKey("post")
	userKey = contextKey("user")
)

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := readIDParam(r, "postID")
		if err != nil || id <= 0 {
			app.notFoundResponse(w, r)
			return
		}

		post, err := app.repository.Posts.GetByID(r.Context(), id)
		if err != nil {
			log.Println(err)
			switch {
			case errors.Is(err, repository.ErrNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := readIDParam(r, "userID")
		if err != nil || id <= 0 {
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

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
