package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hayohtee/social/internal/data"
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
