package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todo-app/api/handlers"
	"github.com/todo-app/api/helpers"
	"github.com/todo-app/api/middleware"
	"github.com/todo-app/internal/application"
)

func Get(app *application.App) *mux.Router {
	r := mux.NewRouter()

	// Ovverride default Notfound/Method not allowed to get structured JSON response using our helper response methods
	r.NotFoundHandler = http.HandlerFunc(helpers.NotFoundErrResponse)
	r.MethodNotAllowedHandler = http.HandlerFunc(helpers.MethodNotAllowedResponse)

	r.HandleFunc("/health", handlers.HealthCheck(app)).Methods(http.MethodGet)
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
