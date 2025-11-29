package main

import "net/http"

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email, max=255"`
	Password string `json:"password" validate:"required,min=3,max=72`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body map[string]string true "Username and password"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	Username : payload.Username,
	Email: payload.Email,
	
	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// store the user
	err := app

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}