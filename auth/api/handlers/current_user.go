package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/todo-app/internal"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/services"
	"github.com/todo-app/pkg/logger"
)

// TODO: Add auth middleware -- this should be protected
func GetCurrentUser(app *application.App) http.HandlerFunc {
	return getCurrentUser(app.IdentityService)
}

func getCurrentUser(service services.IdentityServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		claims, ok := identity.GetClaimsFromContext(r.Context())
		if !ok {
			internal.ErrBadRequest(errors.New("failed to get claims from request context"), "bad request").Send(w)
			return
		}

		// user, err := service.GetById()
		logger.Info.Printf("Claims: %+v", claims)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"Success": "true"})
	}
}
