package main

import (
	"net/http"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running and healthy
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
