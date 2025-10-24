package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s error: %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusInternalServerError, "the server encountered a problem and could not complete your request")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request: %s error: %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusNotFound, "not found")
}
