package main

import (
	"log"
	"net/http"
)

// GetUserFeed godoc
// @Summary Get user feed
// @Description Get personalized feed with posts from followed users
// @Tags Feed
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// For now, we're using user ID 1 for testing
	// In a real app, you'd get this from the authenticated user's session
	userID := int64(1)

	log.Printf("Fetching feed for user ID: %d\n", userID)

	feed, err := app.store.Posts.GetUserFeed(ctx, userID)
	if err != nil {
		log.Printf("Error fetching feed for user %d: %v\n", userID, err)
		app.internalServerError(w, r, err)
		return
	}

	log.Printf("Successfully fetched %d posts for user ID: %d\n", len(feed), userID)

	if err := app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   feed,
	}); err != nil {
		log.Printf("Error sending response: %v\n", err)
		app.internalServerError(w, r, err)
	}
}
