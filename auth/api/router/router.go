package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todo-app/api/handlers"
	"github.com/todo-app/api/middleware"
	"github.com/todo-app/internal/application"
)

func Get(app *application.App) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handlers.HealthCheck(app))

	// authMiddleware := alice.New(middleware.AuthenticationMiddleware)
	r.HandleFunc("/register", handlers.Register(app)).Methods(http.MethodPost)
	r.HandleFunc("/signin", handlers.Login(app)).Methods(http.MethodPost)
	r.HandleFunc("/currentuser", middleware.AuthenticationMiddleware(handlers.GetCurrentUser(app))).Methods(http.MethodGet)
	http.Handle("/", r)
	// Standard Middlewares applied on every request
	r.Use(middleware.SecureHeaders)
	r.Use(middleware.RequestLog)
	r.Use(middleware.PanicRecovery)
	return r
}
