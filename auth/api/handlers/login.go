package handlers

import (
	"errors"
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
		err := helpers.ReadJSON(w, r, &loginReq)
		if err != nil {
			helpers.BadRequestErrResponseWithMsg(w, r, err)
			return
		}

		user, err := service.HandleLogin(&loginReq)
		if err != nil {
			switch {
			case errors.Is(err, identity.ErrInvalidCredentials):
				helpers.InvalidCredentialsResponse(w, r, err)
			case errors.Is(err, identity.ErrUserNotActivated):
				helpers.BadRequestErrResponseWithMsg(w, r, errors.New("you are unable to login due to your account not being activated. Please check your email and activate your account"))
			default:
				helpers.ServerErrReponse(w, r, err)
			}
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
