package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nati3514/Social/internal/store"
)

type createPostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags" validate:"omitempty"`
	Version *int32    `json:"version" validate:"omitempty"`
}

type postCtxKey struct{}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with title, content, and tags
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body createPostPayload true "Post data"
// @Success 201 {object} store.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		//TODO: get user id from auth
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get a specific post by ID with all comments
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Success 200 {object} store.Post
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostsID(ctx, id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid post ID"))
			return
		}

		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, errors.New("post not found"))
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtxKey{}, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) (*store.Post, error) {
	if post, ok := r.Context().Value(postCtxKey{}).(*store.Post); ok {
		return post, nil
	}
	return nil, errors.New("post not found in context")
}

// UpdatePost godoc
// @Summary Update a post
// @Description Partially update a post with optimistic concurrency control
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Param post body UpdatePostRequest true "Post update data"
// @Success 200 {object} store.Post
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string "Version conflict"
// @Failure 500 {object} map[string]string
// @Router /posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid post ID"))
		return
	}

	// Parse the update request
	var input UpdatePostRequest
	if err := readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get a fresh copy of the post from the database
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, errors.New("post not found"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// Check version if provided
	if input.Version != nil && *input.Version != post.Version {
		app.errorResponse(w, http.StatusConflict, "edit conflict: post has been modified by another user")
		return
	}

	// Apply updates
	if input.Title != nil {
		if *input.Title == "" {
			app.badRequestResponse(w, r, errors.New("title cannot be empty"))
			return
		}
		post.Title = *input.Title
	}

	if input.Content != nil {
		if *input.Content == "" {
			app.badRequestResponse(w, r, errors.New("content cannot be empty"))
			return
		}
		post.Content = *input.Content
	}

	if input.Tags != nil {
		post.Tags = *input.Tags
	}

	// Attempt to update the post
	if err := app.store.Posts.Update(ctx, post); err != nil {
		switch {
		case errors.Is(err, store.ErrEditConflict):
			app.errorResponse(w, http.StatusConflict, "edit conflict: post has been modified by another user")
			return
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	// Get the updated post to return
	updatedPost, err := app.store.Posts.GetByID(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, updatedPost); err != nil {
		app.internalServerError(w, r, err)
	}
}
