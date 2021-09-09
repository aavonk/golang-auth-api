package handlers

import (
	"errors"
	"net/http"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/services"
)

func GetCurrentUser(app *application.App) http.HandlerFunc {
	return getCurrentUser(app.IdentityService)
}

func getCurrentUser(service services.IdentityServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := identity.GetClaimsFromContext(r.Context())
		if !ok {
			helpers.ServerErrReponse(w, r, errors.New("failed getting user claims context from request"))
			return
		}

		user := service.GetUserById(claims.UserId.String())

		if user.IsEmpty() {
			helpers.NotFoundErrResponse(w, r)
			return
		}

		err := helpers.SendJSON(w, http.StatusOK, user.ToHTTPResponse(), nil)

		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}
	}
}
