package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hayohtee/social/internal/data"
	"github.com/hayohtee/social/internal/validator"
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

func readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
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

func readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func generateToken() (data.Token, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return data.Token{}, err
	}

	var token data.Token
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}
