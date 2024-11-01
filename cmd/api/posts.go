package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/hayohtee/social/internal/data"
	"github.com/hayohtee/social/internal/repository"
	"github.com/hayohtee/social/internal/validator"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var userID int64 = 1
	post := data.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID,
		Tags:    input.Tags,
	}

	v := validator.New()

	if data.ValidatePost(v, post); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.repository.Posts.Create(r.Context(), &post); err != nil {
		app.serverErrorResponse(w, r, err)
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

	if len(response.Tags) == 0 {
		response.Tags = make([]string, 0)
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"post": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post, ok := getPostFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
		return
	}

	comments, err := app.repository.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var response struct {
		ID        int64                  `json:"id"`
		UserID    int64                  `json:"user_id"`
		Title     string                 `json:"title"`
		Content   string                 `json:"content"`
		Tags      []string               `json:"tags"`
		CreatedAt time.Time              `json:"created_at"`
		UpdatedAt time.Time              `json:"updated_at"`
		Comments  []data.CommentWithUser `json:"comments"`
	}

	response.ID = post.ID
	response.UserID = post.UserID
	response.Title = post.Title
	response.Content = post.Content
	response.Tags = post.Tags
	response.CreatedAt = post.CreatedAt
	response.UpdatedAt = post.UpdatedAt
	response.Comments = comments

	if len(response.Tags) == 0 {
		response.Tags = []string{}
	}

	if len(response.Comments) == 0 {
		response.Comments = []data.CommentWithUser{}
	}

	if err = app.writeJSON(w, http.StatusOK, envelope{"post": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r, "postID")
	if err != nil || id < 0 {
		app.notFoundResponse(w, r)
		return
	}

	if err = app.repository.Posts.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err = app.writeJSON(w, http.StatusOK, envelope{"message": "post deleted successfully"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, ok := getPostFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
		return
	}

	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	v := validator.New()
	if data.ValidatePost(v, post); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.repository.Posts.Update(r.Context(), &post); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if len(post.Tags) == 0 {
		post.Tags = []string{}
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
