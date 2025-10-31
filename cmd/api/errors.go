package main

import (
	//"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	//app.logger.Error(err.Error())
	app.errorResponse(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("conflict error: %s path: %s", r.Method, r.URL.Path, err)
	app.errorResponse(w, http.StatusConflict, err.Error())
}