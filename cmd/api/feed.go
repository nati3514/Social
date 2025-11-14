package main

import (
	"log"
	"net/http"

	"github.com/nati3514/Social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination, filters, sort
	fq := store.PaginatedFeedQuery{
		Limit: 20,
		Offset: 0,
		Sort: "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	
	// For now, we're using user ID 1 for testing
	// In a real app, you'd get this from the authenticated user's session
	userID := int64(1)
	
	log.Printf("Fetching feed for user ID: %d\n", userID)
	
	feed, err := app.store.Posts.GetUserFeed(ctx, userID, fq)
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
