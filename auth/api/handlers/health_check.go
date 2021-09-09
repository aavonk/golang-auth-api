package handlers

import (
	"net/http"

	"github.com/todo-app/api/helpers"
	"github.com/todo-app/internal/application"
	"github.com/todo-app/pkg/config"
)

func HealthCheck(app *application.App) http.HandlerFunc {

	return healthCheck(app.Confg)
}

func healthCheck(cfg *config.Confg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"status": "available",
			"system_info": map[string]string{
				"environment": cfg.GetEnvironment(),
				"version":     cfg.GetVersion(),
			},
		}

		err := helpers.SendJSON(w, http.StatusOK, response, nil)
		if err != nil {
			helpers.ServerErrReponse(w, r, err)
			return
		}

	}

}
