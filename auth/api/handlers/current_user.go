package handlers

import (
	"encoding/json"
	"net/http"

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

		cookieValue, err := identity.GetCookie(r)

		if err != nil {
			logger.Error.Println(err)
			return
		}

		claims, err := identity.ExtractClaimsFromToken(cookieValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error.Println(err)
			return
		}

		logger.Info.Printf("Claims: %+v", claims)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cookieValue)
	}
}
