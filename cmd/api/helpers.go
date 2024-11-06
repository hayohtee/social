package main

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hayohtee/social/internal/model"
	"github.com/hayohtee/social/internal/validator"
)

func readIDParam(r *http.Request, key string) (int64, error) {
	idParam := chi.URLParam(r, key)
	return strconv.ParseInt(idParam, 10, 64)
}

func getPostFromContext(ctx context.Context) (model.Post, bool) {
	post, ok := ctx.Value(postKey).(model.Post)
	return post, ok
}

func getUserFromContext(ctx context.Context) (model.User, bool) {
	user, ok := ctx.Value(userKey).(model.User)
	return user, ok
}

func readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}
	return i
}