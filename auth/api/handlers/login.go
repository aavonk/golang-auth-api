package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/services"
)

func Login(app *application.App) http.HandlerFunc {
	return login(app.IdentityService)
}

func login(service services.IdentityServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginReq identity.LoginRequest

		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			helpers.BadRequestErrResponse(w, r, err)
			return
		}

		user, err := service.HandleLogin(&loginReq)
		if err != nil {
			helpers.InvalidCredentialsResponse(w, r, err)

			return
		}

		err = identity.SetCookie(w, user)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

		userResponse := user.ToHTTPResponse()
		helpers.SendJSON(w, http.StatusOK, userResponse, nil)
	}
}
