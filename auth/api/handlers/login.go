package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/todo-app/internal"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/services"
)

func Login(app *application.App) http.HandlerFunc {
	return login(app.IdentityService)
}

func login(service services.IdentityServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			internal.ErrUnprocessableEntity(err, "cannot parse body").Send(w)
			return
		}

		var loginReq identity.LoginRequest

		err = json.Unmarshal(body, &loginReq)

		if err != nil {
			internal.ErrUnprocessableEntity(err, "cannot parse body").Send(w)
			return
		}

		user, err := service.HandleLogin(&loginReq)
		if err != nil {
			// Error will read "invalid credentials"
			internal.ErrUnprocessableEntity(err, err.Error()).Send(w)
			return
		}

		err = identity.SetAndSaveSession(r, w, user)

		if err != nil {
			internal.ErrInternalServer(err, "internal error message").Send(w)
			return
		}

		userResponse := user.ToHTTPResponse()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&userResponse)

	}
}
