package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/todo-app/internal"

	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/services"
)

func register(service services.IdentityServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			internal.ErrUnprocessableEntity(err, err.Error()).Send(w)
			return
		}

		var user domain.User

		err = json.Unmarshal(body, &user)

		if err != nil {
			internal.ErrUnprocessableEntity(err, err.Error()).Send(w)
			return
		}

		createdUser, err := service.HandleRegister(&user)

		if err != nil {
			internal.ErrBadRequest(err, "Bad Request").Send(w)
			return
		}

		userResponse := createdUser.ToHTTPResponse()

		err = identity.SetAndSaveSession(r, w, user)
		if err != nil {
			internal.ErrInternalServer(err, "internal error message").Send(w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&userResponse)

	}

}

func Register(app *application.App) http.HandlerFunc {
	return register(app.IdentityService)
}
