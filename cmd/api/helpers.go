package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func readIDParam(r *http.Request, key string) (int64, error) {
	idParam := chi.URLParam(r, key)
	return strconv.ParseInt(idParam, 10, 64)
}
