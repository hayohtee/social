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

	v := validator.New()

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(input.Content != "", "content", "must be provided")
	v.Check(len(input.Content) <= 1000, "content", "must be more than 1000 bytes long")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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

	var resp struct {
		ID        int64                  `json:"id"`
		UserID    int64                  `json:"user_id"`
		Title     string                 `json:"title"`
		Content   string                 `json:"content"`
		Tags      []string               `json:"tags"`
		CreatedAt time.Time              `json:"created_at"`
		UpdatedAt time.Time              `json:"updated_at"`
		Comments  []data.CommentWithUser `json:"comments"`
	}

	resp.ID = post.ID
	resp.UserID = post.UserID
	resp.Title = post.Title
	resp.Content = post.Content
	resp.Tags = post.Tags
	resp.CreatedAt = post.CreatedAt
	resp.UpdatedAt = post.UpdatedAt
	resp.Comments = comments

	if len(resp.Tags) == 0 {
		resp.Tags = []string{}
	}

	if len(resp.Comments) == 0 {
		resp.Comments = []data.CommentWithUser{}
	}

	if err = app.writeJSON(w, http.StatusOK, envelope{"post": resp}, nil); err != nil {
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
	post, ok := getPostFromContext(r.Context())
	if !ok {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
