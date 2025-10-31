package main

import (
	"encoding/json"
	"net/http"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJson(w, status, &envelope{Error: message})
}

// jsonResponse writes a JSON response with the given status code and data
func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    // Wrap the data in a "data" field
    response := map[string]any{
        "data": data,
    }
    
    return json.NewEncoder(w).Encode(response)
}

// errorResponse writes a JSON error response
func (app *application) errorResponse(w http.ResponseWriter, status int, message string) error {
    type errorResponse struct {
        Error string `json:"error"`
    }
    return app.jsonResponse(w, status, errorResponse{Error: message})
}
