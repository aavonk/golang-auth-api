package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/todo-app/api/helpers"

	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/services"
)

func register(service services.IdentityServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			helpers.UnprocessableErrResponse(w, r, err)
			return
		}

		var user domain.User

		err = json.Unmarshal(body, &user)

		if err != nil {
			helpers.UnprocessableErrResponse(w, r, err)

			return
		}

		createdUser, err := service.HandleRegister(&user)

		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

		userResponse := createdUser.ToHTTPResponse()

		err = identity.SetCookie(w, user)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)

			return
		}

		helpers.SendJSON(w, http.StatusCreated, userResponse, nil)

	}

}

func Register(app *application.App) http.HandlerFunc {
	return register(app.IdentityService)
}
