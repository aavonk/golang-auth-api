package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/todo-app/internal/application"
	"github.com/todo-app/pkg/config"
)

func HealthCheck(app *application.App) http.HandlerFunc {

	return healthCheck(app.Confg)
}

func healthCheck(cfg *config.Confg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"status":      "available",
			"environment": cfg.GetEnvironment(),
			"version":     cfg.GetVersion(),
		}

		json.NewEncoder(w).Encode(response)
	}

}
