package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/todo-app/api/handlers"
	"github.com/todo-app/api/middleware"
	"github.com/todo-app/internal/application"
)

func Get(app *application.App) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "time": time.Now().UTC().String()})
	})

	// authMiddleware := alice.New(middleware.AuthenticationMiddleware)
	r.HandleFunc("/register", handlers.Register(app)).Methods("POST")
	r.HandleFunc("/signin", handlers.Login(app)).Methods("POST")
	r.HandleFunc("/currentuser", middleware.AuthenticationMiddleware(handlers.GetCurrentUser(app))).Methods("GET")
	http.Handle("/", r)
	// Standard Middlewares applied on every request
	r.Use(middleware.SecureHeaders)
	r.Use(middleware.RequestLog)
	r.Use(middleware.PanicRecovery)
	return r
}
