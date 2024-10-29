package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hayohtee/social/internal/data"
	"github.com/hayohtee/social/internal/repository"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := readJSON(w, r, &input); err != nil {
		log.Println(err)
		writeJSONError(w, http.StatusBadRequest, "bad request")
		return
	}

	var userID int64 = 1
	post := data.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID,
		Tags:    input.Tags,
	}

	if err := app.repository.Posts.Create(r.Context(), &post); err != nil {
		log.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var response struct {
		ID        int64     `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Tags      []string  `json:"tags"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	response.ID = post.ID
	response.Content = post.Content
	response.Title = post.Title
	response.Tags = post.Tags
	response.CreatedAt = post.CreatedAt
	response.UpdatedAt = post.UpdatedAt

	if err := writeJSON(w, http.StatusCreated, envelope{"post": response}); err != nil {
		log.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 16, 64)
	if err != nil || id < 0 {
		writeJSONError(w, http.StatusBadRequest, "id must be a positive number")
		return
	}

	post, err := app.repository.Posts.GetByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, repository.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
		default:
			writeJSONError(w, http.StatusInternalServerError, "error processing the request")
		}
		return
	}

	var resp struct {
		ID        int64     `json:"id"`
		UserID    int64     `json:"user_id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Tags      []string  `json:"tags"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	resp.ID = post.ID
	resp.UserID = post.UserID
	resp.Title = post.Title
	resp.Content = post.Content
	resp.Tags = post.Tags
	resp.CreatedAt = post.CreatedAt
	resp.UpdatedAt = post.UpdatedAt

	if err = writeJSON(w, http.StatusOK, envelope{"post": resp}); err != nil {
		log.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "error processing the request")
	}
}
