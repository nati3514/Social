package main

import (
	"net/http"
	"github.com/nati3514/Social/internal/store"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	
	app.store.Posts.Create(r.Context(), &store.Post{})
}