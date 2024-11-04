package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hayohtee/social/internal/data"
)

func readIDParam(r *http.Request, key string) (int64, error) {
	idParam := chi.URLParam(r, key)
	return strconv.ParseInt(idParam, 10, 64)
}

func getPostFromContext(ctx context.Context) (data.Post, bool) {
	post, ok := ctx.Value(postKey).(data.Post)
	return post, ok
}

func getUserFromContext(ctx context.Context) (data.User, bool) {
	user, ok := ctx.Value(userKey).(data.User)
	return user, ok
}
